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

	appBar := createMainBar(window, browser)
	rootFrame.AttachWidget(appBar)

	viewPort := createViewport(window, browser)
	rootFrame.AttachWidget(viewPort)

	debugFrame := createDebugFrame(window, browser)
	rootFrame.AttachWidget(debugFrame)

	window.SetRootFrame(rootFrame)
	app.AddWindow(window)
	window.Show()
	perf.Stop("ui-creation")

	app.Run(func() {})
}

func createViewport(window *mustard.Window, browser *structs.WebBrowser) *mustard.ContextWidget {
	return mustard.CreateContextWidget(func(ctx *gg.Context) {
		/*
			Parsing the document again is a very hacky solution to solve
			layout problems, the solution to those problems is to never modify
			the DOM Tree itself, it should be deep cloned as the render tree
			which will be modified inside this callback.
		*/
		perf.Start("render")
		bun.RenderDocument(ctx, ketchup.ParseDocument(browser.Document.RawDocument))
		perf.Stop("render")

		profiles := perf.GetAllProfiles()
		fmt.Println("-----------------")
		for _, profile := range profiles {
			fmt.Println(profile.GetName(), "took", profile.GetElapsedTime())
		}
	})
}
