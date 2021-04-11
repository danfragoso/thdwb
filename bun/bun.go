package bun

import (
	"fmt"

	gg "github.com/danfragoso/thdwb/gg"
	hotdog "github.com/danfragoso/thdwb/hotdog"
)

func RenderDocument(ctx *gg.Context, document *hotdog.Document, experimentalLayout bool) error {
	if !experimentalLayout {
		body, err := document.DOM.FindChildByName("body")
		if err != nil {
			// TODO: Handle documents without body elements by synthesizing one.
			return err
		}

		document.DOM.RenderBox.Width = float64(ctx.Width())
		document.DOM.RenderBox.Height = float64(ctx.Height())

		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		layoutDOM(ctx, body, 0)
	} else {
		html, err := document.DOM.FindChildByName("html")
		if err != nil {
			return err
		}

		renderTree := createRenderTree(html)
		renderTree.RenderBox.Width = float64(ctx.Width())
		renderTree.RenderBox.Height = float64(ctx.Height())

		layoutNode(ctx, renderTree)
		paintNode(ctx, renderTree)
		paintText(ctx, renderTree)

		renderTree.Print(0)
	}

	return nil
}

func getNodeContent(NodeDOM *hotdog.NodeDOM) string {
	return NodeDOM.Content
}

func getElementName(NodeDOM *hotdog.NodeDOM) string {
	return NodeDOM.Element
}

func getNodeChildren(NodeDOM *hotdog.NodeDOM) []*hotdog.NodeDOM {
	return NodeDOM.Children
}

func walkDOM(TreeDOM *hotdog.NodeDOM, d string) {
	fmt.Println(d, getElementName(TreeDOM))
	nodeChildren := getNodeChildren(TreeDOM)

	for i := 0; i < len(nodeChildren); i++ {
		walkDOM(nodeChildren[i], d+"-")
	}
}

func layoutDOM(ctx *gg.Context, node *hotdog.NodeDOM, childIdx int) {
	nodeChildren := getNodeChildren(node)

	node.RenderBox = &hotdog.RenderBox{}
	calculateNode(ctx, node, childIdx)

	for i := 0; i < len(nodeChildren); i++ {
		layoutDOM(ctx, nodeChildren[i], i)
		node.RenderBox.Height += nodeChildren[i].RenderBox.Height
	}

	paintNode(ctx, node)
}

func paintNode(ctx *gg.Context, node *hotdog.NodeDOM) {
	switch node.Style.Display {
	case "block":
		paintBlockElement(ctx, node)
	case "inline":
		paintInlineElement(ctx, node)
	case "list-item":
		paintListItemElement(ctx, node)
	}
}

func calculateNode(ctx *gg.Context, node *hotdog.NodeDOM, postion int) {
	switch node.Style.Display {
	case "block":
		calculateBlockLayout(ctx, node, postion)
	case "inline":
		calculateInlineLayout(ctx, node, postion)
	case "list-item":
		calculateListItemLayout(ctx, node, postion)
	}
}

func GetPageTitle(TreeDOM *hotdog.NodeDOM) string {
	nodeChildren := getNodeChildren(TreeDOM)
	pageTitle := "Sem Titulo"

	if getElementName(TreeDOM) == "title" {
		return getNodeContent(TreeDOM)
	}

	for i := 0; i < len(nodeChildren); i++ {
		nPageTitle := GetPageTitle(nodeChildren[i])

		if nPageTitle != "Sem Titulo" {
			pageTitle = nPageTitle
		}
	}

	return pageTitle
}
