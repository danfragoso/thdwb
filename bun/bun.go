package bun

import (
	"fmt"

	gg "thdwb/gg"
	structs "thdwb/structs"
)

func RenderDocument(ctx *gg.Context, document *structs.HTMLDocument) {
	body := document.RootElement.Children[1]

	document.RootElement.RenderBox.Width = float64(ctx.Width())
	document.RootElement.RenderBox.Height = float64(ctx.Height())

	ctx.SetRGB(1, 1, 1)
	ctx.Clear()

	layoutDOM(ctx, body, 0)
	if document.SelectedElement != nil {
		//paintDebugRect(ctx, document.SelectedElement)
	}
}

func paintDebugRect(ctx *gg.Context, node *structs.NodeDOM) {
	debugString := node.Element + " {" + fmt.Sprint(node.RenderBox.Top, node.RenderBox.Left, node.RenderBox.Width, node.RenderBox.Height) + "}"
	ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(.2, .8, .4, .3)
	ctx.Fill()

	w, h := ctx.MeasureString(debugString)
	ctx.DrawRectangle(node.RenderBox.Left+node.RenderBox.Width-w-4, node.RenderBox.Top+node.RenderBox.Height, w+4, h+4)
	ctx.SetRGB(1, 1, 0)
	ctx.Fill()

	ctx.SetRGB(0, 0, 0)
	ctx.DrawString(debugString, node.RenderBox.Left+node.RenderBox.Width-w-2, node.RenderBox.Top+node.RenderBox.Height+h+2)
	ctx.Fill()
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

	node.RenderBox = &structs.RenderBox{}
	calculateNode(ctx, node, childIdx)

	for i := 0; i < len(nodeChildren); i++ {
		layoutDOM(ctx, nodeChildren[i], i)
		node.RenderBox.Height += nodeChildren[i].RenderBox.Height
	}

	paintNode(ctx, node)
}

func paintNode(ctx *gg.Context, node *structs.NodeDOM) {
	switch node.Style.Display {
	case "block":
		paintBlockElement(ctx, node)
	case "inline":
		paintInlineElement(ctx, node)
	case "list-item":
		paintListItemElement(ctx, node)
	}
}

func calculateNode(ctx *gg.Context, node *structs.NodeDOM, postion int) {
	switch node.Style.Display {
	case "block":
		calculateBlockLayout(ctx, node, postion)
	case "inline":
		calculateInlineLayout(ctx, node, postion)
	case "list-item":
		calculateListItemLayout(ctx, node, postion)
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
