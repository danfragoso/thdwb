package main

import (
	"flag"
	"runtime"
	bun "thdwb/bun"
	"thdwb/gg"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	structs "thdwb/structs"

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

	settings := LoadSettings(*settingsPath)

	browser := &structs.WebBrowser{
		Document: &structs.HTMLDocument{},
		History:  &structs.History{},
		Profiler: profiler.CreateProfiler(),
		BuildInfo: &structs.BuildInfo{
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

	appBar, statusLabel, menuButton, nextButton, previousButton, urlInput := createMainBar(window, browser)
	rootFrame.AttachWidget(appBar)

	loadDocument(browser, settings.Homepage)
	urlInput.SetValue(browser.Document.URL.String())

	scrollBar := mustard.CreateScrollBarWidget(mustard.VerticalScrollBar)
	scrollBar.SetTrackColor("#ccc")
	scrollBar.SetThumbColor("#aaa")
	scrollBar.SetWidth(12)

	viewPort := mustard.CreateCanvasWidget(func(canvas *mustard.CanvasWidget) {
		go func() {
			browser.Profiler.Start("render")
			ctxBounds := canvas.GetContext().Image().Bounds()
			drawingContext := gg.NewContext(ctxBounds.Max.X, ctxBounds.Max.Y)

			err := bun.RenderDocument(drawingContext, browser.Document)
			if err != nil {
				structs.Log("render", "Can't render page: "+err.Error())
			}

			canvas.SetContext(drawingContext)
			canvas.RequestRepaint()
			browser.Profiler.Stop("render")

			statusLabel.SetContent(createStatusLabel(browser.Profiler))
			statusLabel.RequestRepaint()
			canvas.RequestRepaint()

			scrollBar.SetScrollerOffset(0)

			body, err := browser.Document.RootElement.FindChildByName("body")
			if err != nil {
				structs.Log("render", "can't find body element: "+err.Error())
				return
			}

			scrollBar.SetScrollerSize(body.RenderBox.Height)
			scrollBar.RequestReflow()
		}()
	})

	browser.Viewport = viewPort
	browser.StatusLabel = statusLabel

	urlInput.SetReturnCallback(func() {
		go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
	})

	window.RegisterButton(menuButton, func() {
		window.AddContextMenuEntry("Home", func() {
			urlInput.SetValue("thdwb://homepage/")
			go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		})

		window.AddContextMenuEntry("History", func() {
			urlInput.SetValue("thdwb://history/")
			go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		})

		window.AddContextMenuEntry("About", func() {
			urlInput.SetValue("thdwb://about/")
			go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		})

		if browser.Document.DebugFlag {
			window.AddContextMenuEntry("Disable element inspector", func() {
				browser.Window.RemoveStaticOverlay("debugOverlay")
				browser.Document.DebugFlag = false
			})
		} else {
			window.AddContextMenuEntry("Enable element inspector", func() {
				browser.Document.DebugFlag = true
			})
		}

		window.DrawContextMenu()
	})

	window.RegisterButton(nextButton, func() {
		if len(browser.History.NextPages()) > 0 {
			browser.History.PopNext()
			urlInput.SetValue(browser.History.Last().String())
			go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		}
	})

	window.RegisterButton(previousButton, func() {
		if browser.History.PageCount() > 1 {
			browser.History.Pop()
			urlInput.SetValue(browser.History.Last().String())
			go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
		}
	})

	window.AttachPointerPositionEventListener(func(pointerX, pointerY float64) {
		if viewPort.IsPointInside(pointerX, pointerY) {
			offset := float64(appBar.GetHeight())
			processPointerPositionEvent(browser, pointerX, pointerY-offset)
		} else {
			browser.Document.SelectedElement = nil
		}
	})

	window.AttachScrollEventListener(func(direction int) {
		scrollStep := 20

		body, err := browser.Document.RootElement.FindChildByName("body")
		if err != nil {
			structs.Log("render", "Can't find body element: "+err.Error())
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
				if browser.Document.SelectedElement != nil {
					if browser.Document.SelectedElement.Element == "a" {
						href := browser.Document.SelectedElement.Attr("href")
						urlInput.SetValue(href)
						go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					}
				}
			} else {
				if browser.Document.SelectedElement != nil {
					window.AddContextMenuEntry("Back", func() {
						previousButton.Click()
					})
					window.AddContextMenuEntry("Reload", func() {
						go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					})
					window.AddContextMenuEntry("History", func() {
						urlInput.SetValue("thdwb://history")
						go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
					})
					window.AddContextMenuEntry("Home", func() {
						urlInput.SetValue("thdwb://homepage")
						go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
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
