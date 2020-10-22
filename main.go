package main

import (
	"flag"
	"fmt"
	"image"
	"runtime"
	bun "thdwb/bun"
	"thdwb/gg"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	structs "thdwb/structs"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var perf *profiler.Profiler

func main() {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	mustard.SetGLFWHints()

	perf = profiler.CreateProfiler()

	defaultPath := "./settings.json"
	settingsPath := flag.String("settings", defaultPath, "This flag sets the location for the browser settings file.")
	flag.Parse()

	settings := LoadSettings(*settingsPath)

	browser := &structs.WebBrowser{Document: &structs.HTMLDocument{}, History: &structs.History{}}

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", settings.WindowWidth, settings.WindowHeight, settings.HiDPI)
	window.EnableContextMenus()
	browser.Window = window

	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, goButton, backButton, urlInput := createMainBar(window, browser)
	urlInput.SetReturnCallback(func() {
		goButton.Click()
	})

	rootFrame.AttachWidget(appBar)

	loadDocument(browser, settings.Homepage)
	urlInput.SetValue(browser.Document.URL.String())

	scrollBar := mustard.CreateScrollBarWidget(mustard.VerticalScrollBar)
	scrollBar.SetTrackColor("#ccc")
	scrollBar.SetThumbColor("#aaa")
	scrollBar.SetWidth(12)

	viewPort := mustard.CreateCanvasWidget(func(canvas *mustard.CanvasWidget) {
		go func() {
			perf.Start("render")
			ctxBounds := canvas.GetContext().Image().Bounds()
			drawingContext := gg.NewContext(ctxBounds.Max.X, ctxBounds.Max.Y)

			err := bun.RenderDocument(drawingContext, browser.Document)
			if err != nil {
				structs.Log("render", "Can't render page: "+err.Error())
			}

			canvas.SetContext(drawingContext)
			canvas.RequestRepaint()
			perf.Stop("render")

			statusLabel.SetContent(createStatusLabel(perf))
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

	window.RegisterButton(menuButton, func() {
		window.AddContextMenuEntry("Home", func() {
			urlInput.SetValue("thdwb://homepage/")
			goButton.Click()
		})

		window.AddContextMenuEntry("History", func() {
			urlInput.SetValue("thdwb://history/")
			goButton.Click()
		})

		window.AddContextMenuEntry("About", func() {
			urlInput.SetValue("thdwb://about/")
			goButton.Click()
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

	window.RegisterButton(goButton, func() {
		go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
	})

	window.RegisterButton(backButton, func() {
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
						backButton.Click()
					})
					window.AddContextMenuEntry("Reload", func() {
						goButton.Click()
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

func processPointerPositionEvent(browser *structs.WebBrowser, x, y float64) {
	y -= float64(browser.Viewport.GetOffset())
	browser.Document.SelectedElement = browser.Document.RootElement.CalcPointIntersection(x, y)

	if browser.Document.SelectedElement != nil && browser.Document.SelectedElement.Element == "a" {
		browser.Window.SetCursor("pointer")
		browser.StatusLabel.SetContent(browser.Document.SelectedElement.Attr("href"))
	} else {
		browser.Window.SetCursor("default")
		browser.StatusLabel.SetContent(createStatusLabel(perf))
	}

	if browser.Document.DebugFlag &&
		browser.Document.SelectedElement != nil &&
		browser.Document.SelectedElement.Element != "html" {
		showDebugOverlay(browser)
	}

	browser.StatusLabel.RequestRepaint()
}

func showDebugOverlay(browser *structs.WebBrowser) {
	browser.Window.RemoveStaticOverlay("debugOverlay")

	debugEl := browser.Document.SelectedElement
	top, left, _, height := debugEl.RenderBox.GetRect()
	ctx := gg.NewContext(int(browser.Document.RootElement.RenderBox.Width), int(height+20))
	paintDebugRect(ctx, debugEl)

	overlay := mustard.CreateStaticOverlay("debugOverlay", ctx, image.Point{
		int(left), int(top) + browser.Viewport.GetTop() + browser.Viewport.GetOffset(),
	})

	browser.Window.AddStaticOverlay(overlay)
}

func paintDebugRect(ctx *gg.Context, node *structs.NodeDOM) {
	debugString := node.Element + " {" + fmt.Sprint(node.RenderBox.Top, node.RenderBox.Left, node.RenderBox.Width, node.RenderBox.Height) + "}"
	ctx.DrawRectangle(0, 0, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(.2, .8, .4, .3)
	ctx.Fill()

	w, h := ctx.MeasureString(debugString)

	if node.RenderBox.Width < w {
		ctx.DrawRectangle(0, node.RenderBox.Height, w+4, h+4)
		ctx.SetRGB(1, 1, 0)
		ctx.Fill()

		ctx.SetRGB(0, 0, 0)
		ctx.DrawString(debugString, 2, node.RenderBox.Height+h)
		ctx.Fill()
	} else {
		ctx.DrawRectangle(node.RenderBox.Width-w-2, node.RenderBox.Height, w+4, h+4)
		ctx.SetRGB(1, 1, 0)
		ctx.Fill()

		ctx.SetRGB(0, 0, 0)
		ctx.DrawString(debugString, node.RenderBox.Width-w, node.RenderBox.Height+h)
		ctx.Fill()
	}
}
