package main

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"

	bun "./bun"
	gg "./gg"
	ketchup "./ketchup"
	mustard "./mustard"
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

	url := os.Args[1]
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)

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
		bun.RenderDocument(ctx, ketchup.ParseDocument(parsedDocument.RawDocument))
	})

	rootFrame.AttachWidget(viewPort)
	window.SetRootFrame(rootFrame)
	app.AddWindow(window)
	window.Show()

	if len(os.Args) >= 3 && os.Args[2] == "debug" {
		debugFrame := createDebugFrame(parsedDocument)
		rootFrame.AttachWidget(debugFrame)
	}

	app.Run(func() {
		frameEvents++
		width, height := window.GetSize()
		statusLabel.SetContent("Processed Events: " + strconv.Itoa(frameEvents) + "; Resolution: " + strconv.Itoa(width) + "X" + strconv.Itoa(height))
	})
}

func createDebugFrame(document *structs.HTMLDocument) *mustard.Frame {
	debugFrame := mustard.CreateFrame(mustard.HorizontalFrame)
	debugBar := mustard.CreateFrame(mustard.HorizontalFrame)
	debugContent := mustard.CreateFrame(mustard.VerticalFrame)

	debugTitleBar := mustard.CreateLabelWidget("Debug View")
	debugTitleBar.SetFontSize(15)
	debugTitleBar.SetFontColor("#111")

	vd := mustard.CreateFrame(mustard.VerticalFrame)
	vd.SetHeight(1)
	vd.SetBackgroundColor("#bdbdbd")

	hd := mustard.CreateFrame(mustard.VerticalFrame)
	hd.SetWidth(1)
	hd.SetBackgroundColor("#bdbdbd")

	debugBar.SetHeight(22)
	debugBar.AttachWidget(vd)
	debugBar.AttachWidget(debugTitleBar)

	debugBar.SetBackgroundColor("#ccc")

	debugFrame.SetBackgroundColor("#fff")
	debugFrame.AttachWidget(debugBar)
	debugFrame.AttachWidget(debugContent)

	documentSource := mustard.CreateTextWidget(document.RawDocument)
	documentSource.SetFontSize(11)
	documentSource.SetBackgroundColor("#fcfcfc")

	treeStr, _ := json.MarshalIndent(document.RootElement, "", "  ")
	jsonTree := mustard.CreateTextWidget(string(treeStr))
	jsonTree.SetFontSize(11)
	jsonTree.SetBackgroundColor("#fcfcfc")

	debugContent.AttachWidget(documentSource)
	debugContent.AttachWidget(hd)
	debugContent.AttachWidget(jsonTree)

	return debugFrame
}
