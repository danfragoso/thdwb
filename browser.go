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
