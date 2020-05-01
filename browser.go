package main

import (
	bun "thdwb/bun"
	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	sauce "thdwb/sauce"
	structs "thdwb/structs"
)

func loadDocument(browser *structs.WebBrowser, url string, callback func()) {
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)
	parsedDocument.URL = url

	browser.Document = parsedDocument
	callback()
}

func loadDocumentFromAsset(document []byte) *structs.HTMLDocument {
	parsedDocument := ketchup.ParseDocument(string(document))
	parsedDocument.URL = "thdwb://homepage/"

	return parsedDocument
}

func loadDocumentFromUrl(browser *structs.WebBrowser, statusLabel *mustard.LabelWidget, urlInput *mustard.InputWidget, viewPort *mustard.ContextWidget) {
	if urlInput.GetValue() != browser.Document.URL {
		statusLabel.SetContent("Loading: " + urlInput.GetValue())

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
			viewPort.RequestRepaint()
			statusLabel.RequestRepaint()
		})
	}
}

func createStatusLabel(perf *profiler.Profiler) string {
	return "Loaded; " +
		"Render: " + perf.GetProfile("render").GetElapsedTime().String() + "; " +
		"Parsing: " + perf.GetProfile("parse").GetElapsedTime().String() + "; "
}
