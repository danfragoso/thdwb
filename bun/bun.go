package bun

import (
	"fmt"

	structs "../structs"
	"github.com/fogleman/gg"
)

func RenderTree(ctx *gg.Context, tree *structs.NodeDOM) {
	renderNode(ctx, tree, 10)
}

func getNodeContent(NodeDOM *structs.NodeDOM) string {
	return NodeDOM.Content
}

func getElementName(NodeDOM *structs.NodeDOM) string {
	return NodeDOM.Element
}

func getNodeChildren(NodeDOM *structs.NodeDOM) []*structs.NodeDOM {
	return NodeDOM.Children
}

func walkDOM(TreeDOM *structs.NodeDOM, d string) {
	fmt.Println(d, getElementName(TreeDOM))
	nodeChildren := getNodeChildren(TreeDOM)

	for i := 0; i < len(nodeChildren); i++ {
		walkDOM(nodeChildren[i], d+"-")
	}
}

func renderNode(ctx *gg.Context, node *structs.NodeDOM, d int) {
	ctx.SetHexColor("#000")
	ctx.DrawString(node.Content, float64(d*10), float64(d*10))
	ctx.Fill()

	nodeChildren := getNodeChildren(node)

	for i := 0; i < len(nodeChildren); i++ {
		renderNode(ctx, nodeChildren[i], d+1)
	}
}

func getPageTitle(TreeDOM *structs.NodeDOM) string {
	nodeChildren := getNodeChildren(TreeDOM)
	pageTitle := "Sem Titulo"

	if getElementName(TreeDOM) == "title" {
		return getNodeContent(TreeDOM)
	}

	for i := 0; i < len(nodeChildren); i++ {
		nPageTitle := getPageTitle(nodeChildren[i])

		if nPageTitle != "Sem Titulo" {
			pageTitle = nPageTitle
		}
	}

	return pageTitle
}
