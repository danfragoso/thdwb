package ketchup

import (
	hotdog "github.com/danfragoso/thdwb/hotdog"
	mayo "github.com/danfragoso/thdwb/mayo"

	"golang.org/x/net/html"
)

func buildKetchupNode(node *html.Node, document *hotdog.Document) *hotdog.NodeDOM {
	var element, content string

	ketchupNode := &hotdog.NodeDOM{}
	attributes := retrieveAttributes(node)

	children := retrieveChildren(node)
	for _, child := range children {
		ketchupChild := buildKetchupNode(child, document)
		ketchupChild.Parent = ketchupNode

		ketchupNode.Children = append(
			ketchupNode.Children,
			ketchupChild,
		)
	}

	switch node.Type {
	case html.TextNode:
		element = "html:text"
		content = node.Data
		break

	case html.ElementNode:
		element = node.Data

	case html.DoctypeNode:
		element = "html:doctype"

	case html.RawNode:
		element = "html:raw"
	}

	ketchupNode.Element = element
	ketchupNode.Content = content

	ketchupNode.Attributes = attributes

	ketchupNode.Document = document

	ketchupNode.NeedsReflow = true
	ketchupNode.NeedsRepaint = true

	ketchupNode.Style = mayo.GetElementStylesheet(element, attributes)
	ketchupNode.RenderBox = &hotdog.RenderBox{}

	return ketchupNode
}

func retrieveChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	if node.FirstChild == nil {
		return children
	}

	child := node.FirstChild
	children = append(children, child)

	for child.NextSibling != nil {
		child = child.NextSibling
		children = append(children, child)
	}

	return children
}

func retrieveAttributes(node *html.Node) []*hotdog.Attribute {
	var attributes []*hotdog.Attribute
	for _, attribute := range node.Attr {
		attributes = append(attributes, &hotdog.Attribute{
			Name:  attribute.Key,
			Value: attribute.Val,
		})
	}

	return attributes
}
