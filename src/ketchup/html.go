package ketchup

import (
	"regexp"
	"strings"

	mayo "thdwb/mayo"
	structs "thdwb/structs"
)

var xmlTag = regexp.MustCompile(`(\<.+?\>)|(\<//?\w+\>\\?)`)
var clTag = regexp.MustCompile(`\<\/\w+\>`)
var selfClosingTag = regexp.MustCompile(`\<.+\/\>`)
var tagContent = regexp.MustCompile(`(.+?)\<\/`)
var tagName = regexp.MustCompile(`(\<\w+)`)
var attr = regexp.MustCompile(`\w+=".+?"`)

func extractAttributes(tag string) []*structs.Attribute {
	rawAttrArray := attr.FindAllString(tag, -1)
	elementAttrs := []*structs.Attribute{}

	for i := 0; i < len(rawAttrArray); i++ {
		attrStringSlice := strings.Split(rawAttrArray[i], "=")
		attr := &structs.Attribute{
			Name:  attrStringSlice[0],
			Value: strings.Trim(attrStringSlice[1], "\""),
		}

		elementAttrs = append(elementAttrs, attr)
	}

	return elementAttrs
}

func isVoidElement(tagName string) bool {
	var isVoid bool
	switch tagName {
	case "area",
		"base",
		"br",
		"col",
		"command",
		"embed",
		"hr",
		"img",
		"input",
		"keygen",
		"link",
		"meta",
		"param",
		"source",
		"track",
		"wbr":
		isVoid = true
	default:
		isVoid = false
	}

	return isVoid
}

func ParsePlainText(document string) *structs.Document {
	documentTitle := "Text Document"
	textDocument := &structs.Document{
		Title: documentTitle,

		RawDocument: document,
		DOM: &structs.NodeDOM{
			Element: "html", NeedsReflow: true, NeedsRepaint: true,
			Style:     mayo.GetElementStylesheet("html", []*structs.Attribute{}),
			RenderBox: &structs.RenderBox{},
		},
	}

	textDocument.DOM.Document = textDocument
	textDocument.DOM.Children = []*structs.NodeDOM{
		&structs.NodeDOM{Element: "head", Document: textDocument,
			Style:     mayo.GetElementStylesheet("head", []*structs.Attribute{}),
			RenderBox: &structs.RenderBox{}, Parent: textDocument.DOM,
		},
		&structs.NodeDOM{
			Element: "body", NeedsReflow: true, NeedsRepaint: true,
			Style:     mayo.GetElementStylesheet("body", []*structs.Attribute{}),
			RenderBox: &structs.RenderBox{}, Document: textDocument,
			Parent: textDocument.DOM,
		},
	}

	documentLines := strings.Split(document, "\n")
	body, _ := textDocument.DOM.FindChildByName("body")
	for _, line := range documentLines {
		body.Children = append(body.Children, &structs.NodeDOM{
			Element: "p", Content: line, RenderBox: &structs.RenderBox{},
			Style:  mayo.GetElementStylesheet("p", []*structs.Attribute{}),
			Parent: body,
		})
	}

	return textDocument
}

func ParseHTML(document string) *structs.Document {
	HTMLDocument := &structs.Document{}

	HTMLDocument.RawDocument = document
	lastNode := HTMLDocument.DOM
	parseDocument := xmlTag.MatchString(document)
	document = strings.ReplaceAll(document, "\n", "")

	for parseDocument == true {
		var currentNode *structs.NodeDOM

		currentTag := xmlTag.FindString(document)
		currentTagIndex := xmlTag.FindStringIndex(document)

		if string(currentTag[1]) == "!" {
			document = strings.Replace(document, currentTag, "", 1)
		} else {
			if clTag.MatchString(currentTag) {
				contentStringMatch := tagContent.FindStringSubmatch(document)
				contentString := ""

				if len(contentStringMatch) > 1 {
					contentString = contentStringMatch[1]
				}

				if clTag.MatchString(contentString) {
					lastNode.Content = ""
				} else {
					if lastNode != nil {
						lastNode.Content = strings.TrimSpace(contentString)
					}
				}

				if lastNode.Parent != nil {
					lastNode = lastNode.Parent
				}
			} else {
				currentTagName := strings.Trim(tagName.FindString(currentTag), "<")

				extractedAttributes := extractAttributes(currentTag)
				elementStylesheet := mayo.GetElementStylesheet(currentTagName, extractedAttributes)

				currentNode = &structs.NodeDOM{
					Element:    currentTagName,
					Content:    "",
					Children:   []*structs.NodeDOM{},
					Attributes: extractedAttributes,
					Style:      elementStylesheet,
					Parent:     lastNode,

					NeedsReflow:  true,
					NeedsRepaint: true,
					RenderBox:    &structs.RenderBox{},

					Document: HTMLDocument,
				}

				if currentTagName == "html" {
					HTMLDocument.DOM = currentNode
					lastNode = HTMLDocument.DOM
				} else {
					lastNode.Children = append(lastNode.Children, currentNode)

					if !isVoidElement(currentTagName) {
						lastNode = currentNode
					}
				}
			}

			document = document[currentTagIndex[1]:]
		}

		if !xmlTag.MatchString(document) {
			parseDocument = false
		}
	}

	return HTMLDocument
}
