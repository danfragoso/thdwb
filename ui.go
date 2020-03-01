package main

import (
	"encoding/json"

	mustard "./mustard"
	structs "./structs"
)

func createDebugFrame(window *mustard.Window, browser *structs.WebBrowser) *mustard.Frame {
	debugFrame := mustard.CreateFrame(mustard.HorizontalFrame)
	debugBar := mustard.CreateFrame(mustard.VerticalFrame)
	debugContent := mustard.CreateFrame(mustard.VerticalFrame)

	debugBar.SetHeight(22)
	debugBar.SetBackgroundColor("#eee")
	toggleDebugButton := mustard.CreateButtonWidget("*")
	toggleDebugButton.SetFontSize(16)
	toggleDebugButton.SetWidth(22)
	toggleDebugButton.SetPadding(2)
	debugFrame.SetHeight(400)

	source := mustard.CreateTextWidget(browser.Document.RawDocument)
	source.SetFontSize(12)

	dv := mustard.CreateFrame(mustard.HorizontalFrame)
	dv.SetBackgroundColor("#999")
	dv.SetWidth(1)

	jsonByte, _ := json.MarshalIndent(browser.Document.RootElement, "", "  ")
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

func createMainBar(window *mustard.Window, browser *structs.WebBrowser) *mustard.Frame {
	appBar := mustard.CreateFrame(mustard.HorizontalFrame)
	appBar.SetHeight(40)

	inputFrame := mustard.CreateFrame(mustard.VerticalFrame)
	urlInput := mustard.CreateInputWidget()
	urlInput.SetValue(browser.Document.URL)
	icon := mustard.CreateFrame(mustard.VerticalFrame)
	img := mustard.CreateImageWidget("logo.png")
	img.SetWidth(50)
	icon.AttachWidget(img)
	icon.SetBackgroundColor("#ddd")
	icon.SetWidth(50)

	inputFrame.AttachWidget(icon)
	inputFrame.AttachWidget(urlInput)
	urlInput.SetFontSize(15)

	buttonFrame := mustard.CreateFrame(mustard.VerticalFrame)
	button := mustard.CreateButtonWidget("Ir")
	button.SetPadding(2)
	button.SetWidth(40)

	window.RegisterButton(button, func() {
		browser.Document = loadDocument(urlInput.GetValue())

		window.RequestRepaint()
	})

	buttonFrame.AttachWidget(button)
	buttonFrame.SetWidth(50)
	buttonFrame.SetBackgroundColor("#ddd")
	inputFrame.AttachWidget(buttonFrame)
	window.RegisterInput(urlInput)

	dv := mustard.CreateFrame(mustard.HorizontalFrame)
	dv.SetBackgroundColor("#ddd")
	dv.SetHeight(6)

	appBar.AttachWidget(dv)
	appBar.AttachWidget(inputFrame)
	appBar.AttachWidget(dv)
	return appBar
}
