package ketchup

import (
	"strings"
	structs "thdwb/structs"

	"golang.org/x/net/html"
)

func ParseHTMLDocument(document string) *structs.HTMLDocument {
	parsedDoc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		panic(err)
	}

	HTMLDocument := &structs.HTMLDocument{}
	HTMLDocument.RawDocument = document

	HTMLDocument.RootElement = buildKetchupNode(parsedDoc, HTMLDocument)
	return HTMLDocument
}
