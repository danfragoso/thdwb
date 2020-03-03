package main

import (
	ketchup "./ketchup"
	sauce "./sauce"
	structs "./structs"
)

func loadDocument(url string) *structs.HTMLDocument {
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)
	parsedDocument.URL = url

	return parsedDocument
}

func loadDocumentFromAsset(document []byte) *structs.HTMLDocument {
	parsedDocument := ketchup.ParseDocument(string(document))
	parsedDocument.URL = "thdwb://homepage/"

	return parsedDocument
}
