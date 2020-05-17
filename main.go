package main

import (
	"runtime"

	assets "thdwb/assets"
	bun "thdwb/bun"
	gg "thdwb/gg"
	ketchup "thdwb/ketchup"
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

	browser := &structs.WebBrowser{Document: loadDocumentFromAsset(assets.HomePage()), History: &structs.History{}}
	browser.History.Push("thdwb://homepage/")

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)

	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, goButton, backButton, urlInput := createMainBar(window, browser)
	urlInput.SetReturnCallback(func() {
		goButton.Click()
	})

	debugFrame := createDebugFrame(window, browser)
	rootFrame.AttachWidget(appBar)
	browser.Document = ketchup.ParseDocument(browser.Document.RawDocument)

	viewPort := mustard.CreateCanvasWidget(func(ctx *gg.Context) {
		bun.RenderDocument(ctx, browser.Document)
	})

	browser.Viewport = viewPort
	//viewPort.EnableScrolling()
	window.RegisterButton(menuButton, func() {
		if debugFrame.GetHeight() != 300 {
			debugFrame.SetHeight(300)
		} else {
			debugFrame.SetHeight(0)
		}
	})

	window.RegisterButton(goButton, func() {
		loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
	})

	window.RegisterButton(backButton, func() {
		browser.History.Pop()
		urlInput.SetValue(browser.History.Last())
		loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
	})

	window.AttachPointerPositionEventListener(func(pointerX, pointerY float64) {
		statusBarOffset := 62.
		processPointerPositionEvent(browser, pointerX, pointerY-statusBarOffset)
	})

	window.AttachScrollEventListener(func(direction int) {
		scrollStep := 10
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
				loadDocumentFromUrl(browser, statusLabel, urlInput, viewPort)
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
	//browser.Viewport.SetDrawingRepaint(true)
	//browser.Viewport.RequestRepaint()
}
