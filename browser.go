package main

import (
	ketchup "./ketchup"
	sauce "./sauce"
	structs "./structs"
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
