package ketchup

import (
	"regexp"
	"strings"
)

var xmlTag = regexp.MustCompile(`(\<.+?\>)|(\<//?\w+\>\\?)`)
var clTag = regexp.MustCompile(`\<\/\w+\>`)
var tagContent = regexp.MustCompile(`(.+?)\<\/`)
var tagName = regexp.MustCompile(`(\<\w+)`)
var attr = regexp.MustCompile(`[ ]\w+=".+?"`)

type DOM_Node struct {
	Element  string      `json:"element"`
	Content  string      `json:"content"`
	Children []*DOM_Node `json:"children"`
	parent   *DOM_Node
}

func ParseHTML(document string) *DOM_Node {
	DOM_Tree := &DOM_Node{
		Element:  "root",
		Content:  "THDWB",
		Children: []*DOM_Node{},
		parent:   nil,
	}

	lastNode := DOM_Tree
	parseDocument := xmlTag.MatchString(document)
	document = strings.ReplaceAll(document, "\n", "")

	for parseDocument == true {
		var currentNode *DOM_Node

		currentTag := xmlTag.FindString(document)
		currentTagIndex := xmlTag.FindStringIndex(document)

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

			lastNode = lastNode.parent
		} else {
			currentTagName := tagName.FindString(currentTag)
			currentNode = &DOM_Node{
				Element:  strings.Trim(currentTagName, "<"),
				Content:  "",
				Children: []*DOM_Node{},
				parent:   lastNode,
			}

			lastNode.Children = append(lastNode.Children, currentNode)
			lastNode = currentNode
		}

		document = document[currentTagIndex[1]:len(document)]

		if !xmlTag.MatchString(document) {
			parseDocument = false
		}
	}

	return DOM_Tree
}
