package main

import (
	"fmt"
	"os"
	"runtime"

	bun "./bun"
	gg "./gg"
	ketchup "./ketchup"
	mustard "./mustard"
	profiler "./profiler"
	structs "./structs"

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
	url := os.Args[1]

	browser := &structs.WebBrowser{
		Document: loadDocument(url),
	}

	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)
	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar, statusLabel, menuButton, goButton := createMainBar(window, browser)
	debugFrame := createDebugFrame(window, browser)
	fmt.Println(statusLabel, goButton, menuButton)
	rootFrame.AttachWidget(appBar)

	// window.RegisterButton(menuButton, func() {
	// 	if debugFrame.GetHeight() != 300 {
	// 		debugFrame.SetHeight(300)
	// 	} else {
	// 		debugFrame.SetHeight(0)
	// 	}
	// })

	viewPort := mustard.CreateContextWidget(func(ctx *gg.Context) {
		perf.Start("parse")
		parsedDoc := ketchup.ParseDocument(browser.Document.RawDocument)
		perf.Stop("parse")

		perf.Start("render")
		bun.RenderDocument(ctx, parsedDoc)
		perf.Stop("render")

		statusLabel.SetContent("Loaded; " +
			"Render: " + perf.GetProfile("render").GetElapsedTime().String() + "; " +
			"Parsing: " + perf.GetProfile("parse").GetElapsedTime().String() + "; ")
	})

	rootFrame.AttachWidget(viewPort)
	rootFrame.AttachWidget(debugFrame)
	window.SetRootFrame(rootFrame)
	app.AddWindow(window)
	window.Show()
	perf.Stop("ui-creation")

	app.Run(func() {})
}
