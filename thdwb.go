package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	bun "./bun"
	gg "./gg"
	ketchup "./ketchup"
	mustard "./mustard"
	profiler "./profiler"
	sauce "./sauce"
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

	perf.Start("fetching")
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	perf.Stop("fetching")

	perf.Start("parsing")
	parsedDocument := ketchup.ParseDocument(htmlString)
	parsedDocument.URL = url
	perf.Stop("parsing")

	parsedDocument.Profiler = perf

	perf.Start("ui-creation")
	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)
	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar := createMainBar(window, parsedDocument)
	rootFrame.AttachWidget(appBar)

	viewPort := createViewport(window, parsedDocument)
	rootFrame.AttachWidget(viewPort)

	debugFrame := createDebugFrame(window, parsedDocument)
	rootFrame.AttachWidget(debugFrame)

	window.SetRootFrame(rootFrame)
	app.AddWindow(window)
	window.Show()
	perf.Stop("ui-creation")

	app.Run(func() {})
}

func createDebugFrame(window *mustard.Window, document *structs.HTMLDocument) *mustard.Frame {
	debugFrame := mustard.CreateFrame(mustard.HorizontalFrame)
	debugBar := mustard.CreateFrame(mustard.VerticalFrame)
	debugContent := mustard.CreateFrame(mustard.VerticalFrame)

	debugBar.SetHeight(22)
	debugBar.SetBackgroundColor("#ddd")
	toggleDebugButton := mustard.CreateButtonWidget("*")
	toggleDebugButton.SetFontSize(16)
	toggleDebugButton.SetWidth(22)
	toggleDebugButton.SetPadding(2)
	debugFrame.SetHeight(400)

	source := mustard.CreateTextWidget(document.RawDocument)
	source.SetFontSize(12)

	dv := mustard.CreateFrame(mustard.HorizontalFrame)
	dv.SetBackgroundColor("#999")
	dv.SetWidth(1)

	jsonByte, _ := json.MarshalIndent(document.RootElement, "", "  ")
	json := mustard.CreateTextWidget(string(jsonByte))
	json.SetWidth(200)
	json.SetFontSize(12)

	debugContent.AttachWidget(json)
	debugContent.AttachWidget(dv)
	debugContent.AttachWidget(source)

	window.RegisterButton(toggleDebugButton, func() {
		if debugFrame.GetHeight() == 22 {
			debugFrame.SetHeight(400)
		} else {
			debugFrame.SetHeight(22)
		}
	})

	debugTitle := mustard.CreateLabelWidget("Show Source")
	debugTitle.SetFontSize(16)

	debugBar.AttachWidget(toggleDebugButton)
	debugBar.AttachWidget(debugTitle)
	debugFrame.AttachWidget(debugBar)
	debugFrame.AttachWidget(debugContent)
	debugFrame.SetHeight(22)
	return debugFrame
}

func createMainBar(window *mustard.Window, document *structs.HTMLDocument) *mustard.Frame {
	appBar := mustard.CreateFrame(mustard.HorizontalFrame)
	appBar.SetHeight(40)

	inputFrame := mustard.CreateFrame(mustard.VerticalFrame)
	urlInput := mustard.CreateInputWidget()
	urlInput.SetValue(document.URL)
	icon := mustard.CreateFrame(mustard.VerticalFrame)
	img := mustard.CreateImageWidget("logo.png")
	img.SetWidth(50)
	icon.AttachWidget(img)
	icon.SetBackgroundColor("#ddd")
	icon.SetWidth(100)

	inputFrame.AttachWidget(icon)
	inputFrame.AttachWidget(urlInput)

	buttonFrame := mustard.CreateFrame(mustard.VerticalFrame)
	button := mustard.CreateButtonWidget("Ir")
	button.SetPadding(2)
	button.SetWidth(40)

	window.RegisterButton(button, func() {
		fmt.Println(urlInput.GetValue())
	})

	buttonFrame.AttachWidget(button)
	buttonFrame.SetWidth(100)
	buttonFrame.SetBackgroundColor("#ddd")
	inputFrame.AttachWidget(buttonFrame)
	window.RegisterInput(urlInput)

	dv := mustard.CreateFrame(mustard.HorizontalFrame)
	dv.SetBackgroundColor("#ddd")
	dv.SetHeight(4)

	appBar.AttachWidget(dv)
	appBar.AttachWidget(inputFrame)
	appBar.AttachWidget(dv)
	return appBar
}

func createViewport(window *mustard.Window, document *structs.HTMLDocument) *mustard.ContextWidget {
	return mustard.CreateContextWidget(func(ctx *gg.Context) {
		/*
			Parsing the document again is a very hacky solution to solve
			layout problems, the solution to those problems is to never modify
			the DOM Tree itself, it should be deep cloned as the render tree
			which will be modified inside this callback.
		*/
		perf.Start("render")
		bun.RenderDocument(ctx, ketchup.ParseDocument(document.RawDocument))
		perf.Stop("render")

		profiles := perf.GetAllProfiles()
		fmt.Println("-----------------")
		for _, profile := range profiles {
			fmt.Println(profile.GetName(), "took", profile.GetElapsedTime())
		}
	})
}
