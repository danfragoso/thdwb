package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"

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

func main() {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	mustard.SetGLFWHints()
	perf := profiler.CreateProfiler()

	url := os.Args[1]

	perf.Start("fetching")
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	perf.Stop("fetching")

	perf.Start("parsing")
	parsedDocument := ketchup.ParseDocument(htmlString)
	perf.Stop("parsing")

	parsedDocument.Profiler = perf

	perf.Start("ui-creation")
	app := mustard.CreateNewApp("THDWB")
	window := mustard.CreateNewWindow("THDWB", 600, 600)
	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)

	appBar := mustard.CreateFrame(mustard.VerticalFrame)

	titleBar := mustard.CreateLabelWidget("THDWB - " + url)
	titleBar.SetFontColor("#fff")

	logo := mustard.CreateImageWidget("logo.png")
	logo.SetWidth(20)

	appBar.SetHeight(28)
	appBar.AttachWidget(logo)
	appBar.AttachWidget(titleBar)
	appBar.SetBackgroundColor("#5f6368")

	rootFrame.AttachWidget(appBar)

	statusBar := mustard.CreateFrame(mustard.HorizontalFrame)
	statusBar.SetBackgroundColor("#babcbe")
	statusBar.SetHeight(20)

	statusLabel := mustard.CreateLabelWidget("Processed Events:")
	statusLabel.SetFontSize(16)
	frameEvents := 0

	rootFrame.AttachWidget(statusBar)
	statusBar.AttachWidget(statusLabel)

	viewPort := mustard.CreateContextWidget(func(ctx *gg.Context) {
		/*
			Parsing the document again is a very hacky solution to solve
			layout problems, the solution to those problems is to never modify
			the DOM Tree itself, it should be deep cloned as the render tree
			which will be modified inside this callback.
		*/
		perf.Start("render")
		bun.RenderDocument(ctx, ketchup.ParseDocument(parsedDocument.RawDocument))
		perf.Stop("render")

		profiles := perf.GetAllProfiles()
		fmt.Println("-----------------")
		for _, profile := range profiles {
			fmt.Println(profile.GetName(), "took", profile.GetElapsedTime())
		}
	})

	rootFrame.AttachWidget(viewPort)
	window.SetRootFrame(rootFrame)
	app.AddWindow(window)
	window.Show()
	perf.Stop("ui-creation")

	debugFrame := createDebugFrame(window, parsedDocument)
	rootFrame.AttachWidget(debugFrame)

	app.Run(func() {
		frameEvents++
		width, height := window.GetSize()
		statusLabel.SetContent("Processed Events: " + strconv.Itoa(frameEvents) + "; Resolution: " + strconv.Itoa(width) + "X" + strconv.Itoa(height))
	})
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

	inputTitle := mustard.CreateInputWidget()
	window.RegisterInput(inputTitle)

	debugBar.AttachWidget(inputTitle)

	debugBar.AttachWidget(toggleDebugButton)
	debugBar.AttachWidget(debugTitle)
	debugFrame.AttachWidget(debugBar)
	debugFrame.AttachWidget(debugContent)
	return debugFrame
}
