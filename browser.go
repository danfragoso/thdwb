package main

import (
	bun "thdwb/bun"
	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	sauce "thdwb/sauce"
	structs "thdwb/structs"
)

func loadDocument(browser *structs.WebBrowser, link string, callback func()) {
	URL := sauce.ParseURL(link)

	if URL.Scheme == "" && URL.Host == "" {
		URL = sauce.ParseURL(browser.Document.URL.Scheme + "://" + browser.Document.URL.Host + URL.Path)
	}

	resource := sauce.GetResource(URL)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)
	parsedDocument.URL = URL

	browser.Document = parsedDocument
	callback()
}

func loadDocumentFromAsset(document []byte) *structs.HTMLDocument {
	parsedDocument := ketchup.ParseDocument(string(document))
	parsedDocument.URL = sauce.ParseURL("thdwb://homepage/")

	return parsedDocument
}

func loadDocumentFromUrl(browser *structs.WebBrowser, statusLabel *mustard.LabelWidget, urlInput *mustard.InputWidget, viewPort *mustard.CanvasWidget) {
	statusLabel.SetContent("Loading: " + urlInput.GetValue())

	go loadDocument(browser, urlInput.GetValue(), func() {
		browser.History.Push(browser.Document.URL.String())

		ctx := viewPort.GetContext()
		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		perf.Start("render")
		bun.RenderDocument(ctx, browser.Document)
		perf.Stop("render")

		statusLabel.SetContent(createStatusLabel(perf))
		viewPort.SetOffset(0)
		viewPort.SetDrawingRepaint(true)
		viewPort.RequestRepaint()
		statusLabel.RequestRepaint()

		urlInput.SetValue(browser.Document.URL.String())
	})
}

func createStatusLabel(perf *profiler.Profiler) string {
	return "Loaded; " +
		"Render took: " + perf.GetProfile("render").GetElapsedTime().String() + "; "
}
