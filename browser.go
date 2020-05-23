package main

import (
	"net/url"
	bun "thdwb/bun"
	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	sauce "thdwb/sauce"
	structs "thdwb/structs"
)

func loadDocument(browser *structs.WebBrowser, link string, callback func()) {
	URL := parseURL(link)
	resource := sauce.GetResource(URL)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)
	parsedDocument.URL = URL

	browser.Document = parsedDocument
	callback()
}

func loadDocumentFromAsset(document []byte) *structs.HTMLDocument {
	parsedDocument := ketchup.ParseDocument(string(document))
	parsedDocument.URL = parseURL("thdwb://homepage/")

	return parsedDocument
}

func loadDocumentFromUrl(browser *structs.WebBrowser, statusLabel *mustard.LabelWidget, urlInput *mustard.InputWidget, viewPort *mustard.CanvasWidget) {
	statusLabel.SetContent("Loading: " + urlInput.GetValue())
	browser.History.Push(urlInput.GetValue())

	go loadDocument(browser, urlInput.GetValue(), func() {
		ctx := viewPort.GetContext()
		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		perf.Start("parse")
		parsedDoc := ketchup.ParseDocument(browser.Document.RawDocument)
		perf.Stop("parse")

		perf.Start("render")
		bun.RenderDocument(ctx, parsedDoc)
		perf.Stop("render")

		statusLabel.SetContent(createStatusLabel(perf))
		viewPort.SetOffset(0)
		viewPort.SetDrawingRepaint(true)
		viewPort.RequestRepaint()
		statusLabel.RequestRepaint()
	})
}

func createStatusLabel(perf *profiler.Profiler) string {
	return "Loaded; " +
		"Render took: " + perf.GetProfile("render").GetElapsedTime().String() + "; " +
		"Parse took: " + perf.GetProfile("parse").GetElapsedTime().String() + "; "
}

func parseURL(link string) *url.URL {
	URL, err := url.Parse(link)
	if err != nil {
		panic("Err parsing URL: " + link)
	}

	return URL
}
