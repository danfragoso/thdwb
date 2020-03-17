package bun

import (
	"fmt"

	assets "../assets"
	gg "../gg"
	structs "../structs"
)

func RenderDocument(ctx *gg.Context, document *structs.HTMLDocument) {
	//tree.Children[0] is head
	body := document.RootElement.Children[1]

	document.RootElement.Style.Width = float64(ctx.Width())
	document.RootElement.Style.Height = float64(ctx.Height())

	layoutDOM(ctx, body, 0)
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

func layoutDOM(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	nodeChildren := getNodeChildren(node)

	if node.Style.Display == "block" {
		calculateBlockLayout(ctx, node, childIdx)

		for i := 0; i < len(nodeChildren); i++ {
			layoutDOM(ctx, nodeChildren[i], i)
		}

		paintBlockElement(ctx, node)
	}
}

func paintBlockElement(ctx *gg.Context, node *structs.NodeDOM) {
	ctx.DrawRectangle(node.Style.Left, node.Style.Top, node.Style.Width, node.Style.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.Style.Left, node.Style.Top+1, 0, 0, node.Style.Width, 1.5, gg.AlignLeft)
	ctx.Fill()
}

func calculateBlockLayout(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	if node.Style.Width == 0 {
		node.Style.Width = node.Parent.Style.Width
	}

	if node.Style.Height == 0 && len(node.Content) > 0 {
		ctx.LoadAssetFont(assets.SansSerif(), node.Style.FontSize)
		node.Style.Height = ctx.MeasureStringWrapped(node.Content, node.Style.Width, 1.5) + 2 + ctx.FontHeight()*.5
	}

	if childIdx > 0 {
		prev := node.Parent.Children[childIdx-1]

		if prev.Style.Display == "block" {
			node.Style.Top = prev.Style.Top + prev.Style.Height
		} else {
			node.Style.Top = prev.Style.Top
		}
	}
}

func GetPageTitle(TreeDOM *structs.NodeDOM) string {
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
