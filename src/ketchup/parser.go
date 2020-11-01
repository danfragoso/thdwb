package ketchup

import (
	"strings"
	structs "thdwb/structs"

	"golang.org/x/net/html"
)

func ParseHTMLDocument(document string) *structs.Document {
	parsedDoc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		panic(err)
	}

	HTMLDocument := &structs.Document{}
	HTMLDocument.RawDocument = document

	HTMLDocument.DOM = buildKetchupNode(parsedDoc, HTMLDocument)
	return HTMLDocument
}
