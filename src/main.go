package main

import (
	"flag"
	"runtime"

	bun "thdwb/bun"
	gg "thdwb/gg"
	hotdog "thdwb/hotdog"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	mustard.SetGLFWHints()

	defaultPath := "./settings.json"
	settingsPath := flag.String("settings", defaultPath, "This flag sets the location for the browser settings file.")
	flag.Parse()

	settings := hotdog.LoadSettings(*settingsPath)

	browser := &hotdog.WebBrowser{
		ActiveDocument: &hotdog.Document{},

		History:  &hotdog.History{},
		Profiler: profiler.CreateProfiler(),
		Settings: settings,

		BuildInfo: &hotdog.BuildInfo{
			GitRevision: gitRevision,
			GitBranch:   gitBranch,
			HostInfo:    hostInfo,
			BuildTime:   buildTime,
		},
	}

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", settings.WindowWidth, settings.WindowHeight, settings.HiDPI)
	window.EnableContextMenus()
	browser.Window = window

	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, nextButton, previousButton, reloadButton, urlInput := createMainBar(window, browser)
	rootFrame.AttachWidget(appBar)

	loadDocument(browser, settings.Homepage)
	urlInput.SetValue(browser.ActiveDocument.URL.String())

	scrollBar := mustard.CreateScrollBarWidget(mustard.VerticalScrollBar)
	scrollBar.SetTrackColor("#ccc")
	scrollBar.SetThumbColor("#aaa")
	scrollBar.SetWidth(12)

	viewPort := mustard.CreateCanvasWidget(func(canvas *mustard.CanvasWidget) {
		go func() {
			browser.Profiler.Start("render")
			ctxBounds := canvas.GetContext().Image().Bounds()
			drawingContext := gg.NewContext(ctxBounds.Max.X, ctxBounds.Max.Y)

			err := bun.RenderDocument(drawingContext, browser.ActiveDocument, settings.ExperimentalLayout)
			if err != nil {
				hotdog.Log("render", "Can't render page: "+err.Error())
			}

			canvas.SetContext(drawingContext)
			canvas.RequestRepaint()
			browser.Profiler.Stop("render")

			statusLabel.SetContent(createStatusLabel(browser.Profiler))
			statusLabel.RequestRepaint()
			canvas.RequestRepaint()

			scrollBar.SetScrollerOffset(0)

			body, err := browser.ActiveDocument.DOM.FindChildByName("body")
			if err != nil {
				hotdog.Log("render", "can't find body element: "+err.Error())
				return
			}

			scrollBar.SetScrollerSize(body.RenderBox.Height)
			scrollBar.RequestReflow()
		}()
	})

	browser.Viewport = viewPort
	browser.StatusLabel = statusLabel

	urlInput.SetReturnCallback(func() {
		loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
	})

	window.RegisterButton(menuButton, func() {
		window.AddContextMenuEntry("Home", func() {
			urlInput.SetValue("thdwb://homepage/")
			loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		})

		window.AddContextMenuEntry("History", func() {
			urlInput.SetValue("thdwb://history/")
			loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		})

		window.AddContextMenuEntry("About", func() {
			urlInput.SetValue("thdwb://about/")
			loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		})

		if browser.ActiveDocument.DebugFlag {
			window.AddContextMenuEntry("Disable debug mode", func() {
				browser.Window.RemoveStaticOverlay("debugOverlay")
				browser.ActiveDocument.DebugFlag = false

				if browser.ActiveDocument.DebugWindow != nil {
					app.DestroyWindow(browser.ActiveDocument.DebugWindow)
					browser.ActiveDocument.DebugWindow = nil
					browser.ActiveDocument.DebugTree = nil
				}
			})
		} else {
			window.AddContextMenuEntry("Enable debug mode", func() {
				browser.ActiveDocument.DebugFlag = true
			})
		}

		if browser.ActiveDocument.DebugFlag {
			if browser.ActiveDocument.DebugWindow != nil {
				window.AddContextMenuEntry("Hide Tree", func() {
					app.DestroyWindow(browser.ActiveDocument.DebugWindow)
					browser.ActiveDocument.DebugWindow = nil
					browser.ActiveDocument.DebugTree = nil
				})
			} else {
				window.AddContextMenuEntry("Show Tree", func() {
					tree := mustard.CreateTreeWidget()

					browser.ActiveDocument.DebugWindow = mustard.CreateNewWindow("HTML tree view", 600, 800, true)
					browser.ActiveDocument.DebugTree = tree

					rFrame := mustard.CreateFrame(mustard.HorizontalFrame)
					tree.SetFontSize(14)
					rFrame.AttachWidget(tree)

					browser.ActiveDocument.DebugWindow.RegisterTree(tree)
					browser.ActiveDocument.DebugWindow.SetRootFrame(rFrame)
					browser.ActiveDocument.DebugWindow.Show()

					app.AddWindow(browser.ActiveDocument.DebugWindow)

					treeNodeDOM := treeNodeFromDOM(browser.ActiveDocument.DOM)
					tree.SetSelectCallback(func(selectedNode *mustard.TreeWidgetNode) {
						if browser.ActiveDocument.DebugFlag {
							child, _ := browser.ActiveDocument.DOM.FindByXPath(selectedNode.Value)
							browser.ActiveDocument.SelectedElement = child
							showDebugOverlay(browser)
						}
					})

					tree.RemoveNodes()
					tree.AddNode(treeNodeDOM)
					tree.RequestRepaint()
				})
			}
		}

		window.DrawContextMenu()
	})

	window.RegisterButton(reloadButton, func() {
		loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
	})

	window.RegisterButton(nextButton, func() {
		if len(browser.History.NextPages()) > 0 {
			browser.History.PopNext()
			urlInput.SetValue(browser.History.Last().String())
			loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		}
	})

	window.RegisterButton(previousButton, func() {
		if browser.History.PageCount() > 1 {
			browser.History.Pop()
			urlInput.SetValue(browser.History.Last().String())
			loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		}
	})

	window.AttachPointerPositionEventListener(func(pointerX, pointerY float64) {
		if viewPort.IsPointInside(pointerX, pointerY) {
			offset := float64(appBar.GetHeight())
			processPointerPositionEvent(browser, pointerX, pointerY-offset)
		} else {
			browser.ActiveDocument.SelectedElement = nil
		}
	})

	window.AttachScrollEventListener(func(direction int) {
		scrollStep := 20

		body, err := browser.ActiveDocument.DOM.FindChildByName("body")
		if err != nil {
			hotdog.Log("render", "Can't find body element: "+err.Error())
			return
		}

		if direction > 0 {
			if viewPort.GetOffset() < 0 {
				viewPort.SetOffset(viewPort.GetOffset() + scrollStep)
			}
		} else {
			documentOffset := viewPort.GetOffset() + int(body.RenderBox.Height)

			if documentOffset >= viewPort.GetHeight() {
				viewPort.SetOffset(viewPort.GetOffset() - scrollStep)
			}
		}

		scrollBar.SetScrollerOffset(float64(viewPort.GetOffset()))
		scrollBar.SetScrollerSize(body.RenderBox.Height)
		scrollBar.RequestReflow()

		browser.Viewport.SetDrawingRepaint(false)
		viewPort.RequestRepaint()

		browser.Window.RemoveStaticOverlay("debugOverlay")
	})

	window.AttachClickEventListener(func(key mustard.MustardKey) {
		if viewPort.IsPointInside(window.GetCursorPosition()) {
			if key == mustard.MouseLeft {
				if browser.ActiveDocument.SelectedElement != nil {
					if browser.ActiveDocument.SelectedElement.Element == "a" {
						href := browser.ActiveDocument.SelectedElement.Attr("href")
						urlInput.SetValue(href)
						loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					}
				}
			} else {
				if browser.ActiveDocument.SelectedElement != nil {
					window.AddContextMenuEntry("Back", func() {
						previousButton.Click()
					})
					window.AddContextMenuEntry("Reload", func() {
						loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					})
					window.AddContextMenuEntry("History", func() {
						urlInput.SetValue("thdwb://history")
						loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					})
					window.AddContextMenuEntry("Home", func() {
						urlInput.SetValue("thdwb://homepage")
						loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					})

					window.DrawContextMenu()
				}
			}
		}
	})

	viewArea := mustard.CreateFrame(mustard.VerticalFrame)
	viewArea.AttachWidget(viewPort)
	viewArea.AttachWidget(scrollBar)

	rootFrame.AttachWidget(viewArea)

	window.SetRootFrame(rootFrame)
	window.Show()

	app.AddWindow(window)
	app.Run(func() {})
}
