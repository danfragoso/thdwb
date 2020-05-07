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

	browser := &structs.WebBrowser{
		Document: loadDocumentFromAsset(assets.HomePage()),
	}

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)

	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, goButton, urlInput := createMainBar(window, browser)
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

	window.AttachPointerPositionEventListener(func(pointerX, pointerY float64) {
		statusBarOffset := 62.
		processPointerPositionEvent(browser, pointerX, pointerY-statusBarOffset)
	})

	window.AttachScrollEventListener(func(direction int) {
		//
	})

	rootFrame.AttachWidget(viewPort)
	rootFrame.AttachWidget(debugFrame)

	window.SetRootFrame(rootFrame)
	window.Show()

	app.AddWindow(window)
	app.Run(func() {})
}

func processPointerPositionEvent(browser *structs.WebBrowser, x, y float64) {
	browser.Document.SelectedElement = browser.Document.RootElement.CalcPointIntersection(x, y)
	browser.Viewport.RequestRepaint()
}
