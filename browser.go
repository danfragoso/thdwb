package main

import (
	"strings"

	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	sauce "thdwb/sauce"
	structs "thdwb/structs"
)

func loadDocument(browser *structs.WebBrowser, link string) {
	URL := sauce.ParseURL(link)

	if URL.Scheme == "" && URL.Host == "" {
		if !strings.HasPrefix(URL.Path, "/") {
			URL.Path = "/" + URL.Path
		}

		URL = sauce.ParseURL(browser.Document.URL.Scheme + "://" + browser.Document.URL.Host + URL.Path)
	}

	resource := sauce.GetResource(URL, browser)
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
	statusLabel.RequestRepaint()

	loadDocument(browser, urlInput.GetValue())
	viewPort.SetOffset(0)
	viewPort.SetDrawingRepaint(true)
	viewPort.RequestRepaint()
	urlInput.SetValue(browser.Document.URL.String())
}

func createStatusLabel(perf *profiler.Profiler) string {
	return "Loaded; " +
		"Render took: " + perf.GetProfile("render").GetElapsedTime().String() + "; "
}
