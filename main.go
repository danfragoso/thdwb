package main

import (
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

	browser := &structs.WebBrowser{Document: &structs.HTMLDocument{}, History: &structs.History{}}

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)
	browser.Window = window

	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, goButton, backButton, urlInput := createMainBar(window, browser)
	urlInput.SetReturnCallback(func() {
		goButton.Click()
	})

	debugFrame := createDebugFrame(window, browser)
	rootFrame.AttachWidget(appBar)

	loadDocument(browser, "thdwb://homepage")
	urlInput.SetValue(browser.Document.URL.String())

	viewPort := mustard.CreateCanvasWidget(func(canvas *mustard.CanvasWidget) {
		go func() {
			perf.Start("render")
			ctxBounds := canvas.GetContext().Image().Bounds()
			drawingContext := gg.NewContext(ctxBounds.Max.X, ctxBounds.Max.Y)

			bun.RenderDocument(drawingContext, browser.Document)
			canvas.SetContext(drawingContext)
			canvas.RequestRepaint()
			perf.Stop("render")

			statusLabel.SetContent(createStatusLabel(perf))
			statusLabel.RequestRepaint()
		}()
	})

	browser.Viewport = viewPort
	window.RegisterButton(menuButton, func() {
		if debugFrame.GetHeight() != 300 {
			debugFrame.SetHeight(300)
		} else {
			debugFrame.SetHeight(0)
		}
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
		if direction > 0 {
			if viewPort.GetOffset() < 0 {
				viewPort.SetOffset(viewPort.GetOffset() + scrollStep)
			}
		} else {
			viewPort.SetOffset(viewPort.GetOffset() - scrollStep)
		}

		browser.Viewport.SetDrawingRepaint(false)
		viewPort.RequestRepaint()
	})

	window.AttachClickEventListener(func() {
		if browser.Document.SelectedElement != nil {
			if browser.Document.SelectedElement.Element == "a" {
				href := browser.Document.SelectedElement.Attr("href")
				urlInput.SetValue(href)
				go loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
			}
		}
	})

	rootFrame.AttachWidget(viewPort)
	rootFrame.AttachWidget(debugFrame)

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
	} else {
		browser.Window.SetCursor("default")
	}

	//browser.Viewport.SetDrawingRepaint(true)
	//browser.Viewport.RequestRepaint()
}
