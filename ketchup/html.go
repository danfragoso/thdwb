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

func ParseDocument(document string) *structs.HTMLDocument {
	HTMLDocument := &structs.HTMLDocument{}

	HTMLDocument.RawDocument = document
	lastNode := HTMLDocument.RootElement
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
					lastNode.Content = strings.TrimSpace(contentString)
				}

				lastNode = lastNode.Parent
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
				}

				if currentTagName == "html" {
					HTMLDocument.RootElement = currentNode
					lastNode = HTMLDocument.RootElement
				} else {
					lastNode.Children = append(lastNode.Children, currentNode)

					if !isVoidElement(currentTagName) {
						lastNode = currentNode
					}
				}
			}

			document = document[currentTagIndex[1]:len(document)]
		}

		if !xmlTag.MatchString(document) {
			parseDocument = false
		}
	}

	return HTMLDocument
}
