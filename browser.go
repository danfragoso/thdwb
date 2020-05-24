package main

import (
	bun "thdwb/bun"
	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	sauce "thdwb/sauce"
	structs "thdwb/structs"
)

func loadDocument(browser *structs.WebBrowser, link string) {
	URL := sauce.ParseURL(link)

	if URL.Scheme == "" && URL.Host == "" {
		URL = sauce.ParseURL(browser.Document.URL.Scheme + "://" + browser.Document.URL.Host + URL.Path)
	}

	resource := sauce.GetResource(URL)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)
	parsedDocument.URL = resource.URL

	browser.Document = parsedDocument

	if browser.History.PageCount() == 0 || browser.History.Last().String() != resource.URL.String() {
		browser.History.Push(resource.URL)
	}
}

func loadDocumentFromUrl(browser *structs.WebBrowser, statusLabel *mustard.LabelWidget, urlInput *mustard.InputWidget, viewPort *mustard.CanvasWidget) {
	statusLabel.SetContent("Loading: " + urlInput.GetValue())

	loadDocument(browser, urlInput.GetValue())

	perf.Start("render")
	bun.RenderDocument(viewPort.GetContext(), browser.Document)
	perf.Stop("render")

	statusLabel.SetContent(createStatusLabel(perf))
	viewPort.SetOffset(0)
	viewPort.SetDrawingRepaint(true)
	viewPort.RequestRepaint()
	statusLabel.RequestRepaint()

	urlInput.SetValue(browser.Document.URL.String())
}

func createStatusLabel(perf *profiler.Profiler) string {
	return "Loaded; " +
		"Render took: " + perf.GetProfile("render").GetElapsedTime().String() + "; "
}
