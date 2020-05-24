package main

import (
	"encoding/json"

	assets "thdwb/assets"
	mustard "thdwb/mustard"
	structs "thdwb/structs"
)

func createDebugFrame(window *mustard.Window, browser *structs.WebBrowser) *mustard.Frame {
	debugFrame := mustard.CreateFrame(mustard.HorizontalFrame)
	debugBar := mustard.CreateFrame(mustard.VerticalFrame)
	debugContent := mustard.CreateFrame(mustard.VerticalFrame)

	debugBar.SetHeight(22)
	debugBar.SetBackgroundColor("#eee")
	debugFrame.SetHeight(400)

	source := mustard.CreateTextWidget("")
	source.SetFontSize(12)

	dv := mustard.CreateFrame(mustard.HorizontalFrame)
	dv.SetBackgroundColor("#999")
	dv.SetWidth(1)

	jsonByte, _ := json.MarshalIndent("", "", "  ")
	json := mustard.CreateTextWidget(string(jsonByte))
	json.SetWidth(200)
	json.SetFontSize(12)

	debugContent.AttachWidget(json)
	debugContent.AttachWidget(dv)
	debugContent.AttachWidget(source)

	debugTitle := mustard.CreateLabelWidget("Source")
	debugTitle.SetBackgroundColor("#eee")
	debugTitle.SetFontSize(16)

	debugBar.AttachWidget(debugTitle)
	debugFrame.AttachWidget(debugBar)
	debugFrame.AttachWidget(debugContent)
	debugFrame.SetHeight(0)
	return debugFrame
}

func createMainBar(window *mustard.Window, browser *structs.WebBrowser) (*mustard.Frame, *mustard.LabelWidget, *mustard.ButtonWidget, *mustard.ButtonWidget, *mustard.ButtonWidget, *mustard.InputWidget) {
	appBar := mustard.CreateFrame(mustard.HorizontalFrame)
	appBar.SetHeight(62)

	inputFrame := mustard.CreateFrame(mustard.VerticalFrame)
	urlInput := mustard.CreateInputWidget()
	icon := mustard.CreateFrame(mustard.VerticalFrame)
	img := mustard.CreateImageWidget(assets.Logo())

	backButton := mustard.CreateButtonWidget(assets.ArrowLeft())
	backButton.SetWidth(30)

	rv := mustard.CreateFrame(mustard.HorizontalFrame)
	rv.SetBackgroundColor("#ddd")
	rv.SetWidth(5)

	img.SetWidth(50)
	icon.AttachWidget(img)
	icon.SetBackgroundColor("#ddd")
	icon.SetWidth(50)

	inputFrame.AttachWidget(icon)
	inputFrame.AttachWidget(backButton)
	inputFrame.AttachWidget(rv)
	inputFrame.AttachWidget(urlInput)
	urlInput.SetFontSize(15)

	buttonFrame := mustard.CreateFrame(mustard.VerticalFrame)

	goButton := mustard.CreateButtonWidget(assets.ArrowRight())
	goButton.SetWidth(30)
	goButton.SetPadding(1)

	toolsButton := mustard.CreateButtonWidget(assets.Menu())
	toolsButton.SetWidth(34)
	toolsButton.SetPadding(1)

	buttonFrame.AttachWidget(goButton)
	buttonFrame.AttachWidget(toolsButton)
	buttonFrame.SetWidth(68)
	buttonFrame.SetBackgroundColor("#ddd")
	inputFrame.AttachWidget(buttonFrame)
	window.RegisterInput(urlInput)

	dv := mustard.CreateFrame(mustard.HorizontalFrame)
	dv.SetBackgroundColor("#ddd")
	dv.SetHeight(6)

	pv := mustard.CreateFrame(mustard.HorizontalFrame)
	pv.SetBackgroundColor("#bfbfbf")
	pv.SetHeight(1)

	statusBar := mustard.CreateFrame(mustard.HorizontalFrame)
	statusLabel := mustard.CreateLabelWidget("The HotDog Web Browser")
	statusLabel.SetBackgroundColor("#ddd")
	statusLabel.SetFontColor("#333")
	statusLabel.SetFontSize(15)
	statusBar.AttachWidget(statusLabel)
	statusBar.SetHeight(20)

	appBar.AttachWidget(dv)
	appBar.AttachWidget(inputFrame)
	appBar.AttachWidget(dv)
	appBar.AttachWidget(pv)
	appBar.AttachWidget(statusBar)
	appBar.AttachWidget(pv)

	return appBar, statusLabel, toolsButton, goButton, backButton, urlInput
}
