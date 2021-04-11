package ketchup

import (
	hotdog "github.com/danfragoso/thdwb/hotdog"
	"strings"

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
