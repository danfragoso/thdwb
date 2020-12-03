package ketchup

import (
	"strings"
	hotdog "thdwb/hotdog"

	"golang.org/x/net/html"
)

func ParseHTMLDocument(document string) *hotdog.Document {
	parsedDoc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		panic(err)
	}

	HTMLDocument := &hotdog.Document{}
	HTMLDocument.RawDocument = document

	HTMLDocument.DOM = buildKetchupNode(parsedDoc, HTMLDocument)
	return HTMLDocument
}
