package main

import (
	assets "thdwb/assets"
	mustard "thdwb/mustard"
	hotdog "thdwb/hotdog"
)

func createMainBar(window *mustard.Window, browser *hotdog.WebBrowser) (*mustard.Frame, *mustard.LabelWidget, *mustard.ButtonWidget, *mustard.ButtonWidget, *mustard.ButtonWidget, *mustard.ButtonWidget, *mustard.InputWidget) {
	appBar := mustard.CreateFrame(mustard.HorizontalFrame)
	appBar.SetHeight(62)

	inputFrame := mustard.CreateFrame(mustard.VerticalFrame)
	urlInput := mustard.CreateInputWidget()
	icon := mustard.CreateFrame(mustard.VerticalFrame)
	img := mustard.CreateImageWidget(assets.Logo())

	previousButton := mustard.CreateButtonWidget("", assets.ArrowLeft())
	previousButton.SetWidth(30)

	nextButton := mustard.CreateButtonWidget("", assets.ArrowRight())
	nextButton.SetWidth(30)

	reloadButton := mustard.CreateButtonWidget("", assets.Reload())
	reloadButton.SetWidth(30)

	toolsButton := mustard.CreateButtonWidget("", assets.Menu())
	toolsButton.SetWidth(34)

	rv := mustard.CreateFrame(mustard.HorizontalFrame)
	rv.SetBackgroundColor("#ddd")
	rv.SetWidth(5)

	img.SetWidth(50)
	icon.AttachWidget(img)
	icon.SetBackgroundColor("#ddd")
	icon.SetWidth(50)

	inputFrame.AttachWidget(icon)
	inputFrame.AttachWidget(previousButton)
	inputFrame.AttachWidget(rv)
	inputFrame.AttachWidget(nextButton)
	inputFrame.AttachWidget(rv)
	inputFrame.AttachWidget(reloadButton)
	inputFrame.AttachWidget(rv)
	inputFrame.AttachWidget(rv)
	inputFrame.AttachWidget(urlInput)
	inputFrame.AttachWidget(rv)
	inputFrame.AttachWidget(toolsButton)
	inputFrame.AttachWidget(rv)

	urlInput.SetFontSize(15)

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

	return appBar, statusLabel, toolsButton, nextButton, previousButton, reloadButton, urlInput
}
