package bun

import (
	"fmt"
	"strings"

	gg "thdwb/gg"
	structs "thdwb/structs"
)

func RenderDocument(ctx *gg.Context, document *structs.HTMLDocument) {
	html := document.RootElement.FindChildByName("html")

	renderTree := createRenderTree(html)
	document.RenderTree = renderTree

	renderTree.RenderBox.Width = float64(ctx.Width())
	renderTree.RenderBox.Height = float64(ctx.Height())
	body := renderTree.FindChildByName("body")

	ctx.SetRGB(1, 1, 1)
	ctx.Clear()

	layoutNode(ctx, body)
	paintNode(ctx, body)

	pRender(body, "-")
	fmt.Print("\n")
}

func createRenderTree(root *structs.NodeDOM) *structs.NodeDOM {
	if root.Style.Display == "none" {
		return nil
	}

	node := &structs.NodeDOM{
		Style:      root.Style,
		Element:    root.Element,
		Content:    root.Content,
		Attributes: root.Attributes,
	}

	node.RenderBox = &structs.RenderBox{}
	for _, child := range root.Children {
		r := createRenderTree(child)
		if r != nil {
			r.Parent = node
			node.Children = append(node.Children, r)
		}
	}

	return node
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

func pRender(TreeDOM *structs.NodeDOM, d string) {
	fmt.Println(d, getElementName(TreeDOM), TreeDOM.RenderBox.String())
	nodeChildren := getNodeChildren(TreeDOM)

	for i := 0; i < len(nodeChildren); i++ {
		pRender(nodeChildren[i], d+"-")
	}
}

func layoutNode(ctx *gg.Context, node *structs.NodeDOM) {
	switch node.Style.Display {
	case "block":
		calcBlockDimension(ctx, node)
		calcBlockPosition(ctx, node)

	case "inline":
		calcInlineDimension(ctx, node)
		calcInlinePosition(ctx, node)

	}

	for _, child := range node.Children {
		layoutNode(ctx, child)

		switch node.Style.Display {
		case "block":
			node.RenderBox.Height += child.RenderBox.Height
		case "inline":
			node.RenderBox.Width += child.RenderBox.Width
			if child.RenderBox.Height > node.RenderBox.Height {
				node.RenderBox.Height = child.RenderBox.Height
			}
		}
	}
}

func calcInlineDimension(ctx *gg.Context, node *structs.NodeDOM) {
	content := strings.TrimSpace(node.Content)

	if len(content) > 0 {
		ctx.SetFont(sansSerif[node.Parent.Style.FontWeight], node.Parent.Style.FontSize)
		node.RenderBox.Width, node.RenderBox.Height = ctx.MeasureString(node.Content)
	}
}

func calcInlinePosition(ctx *gg.Context, node *structs.NodeDOM) {
	node.RenderBox.Top = node.Parent.RenderBox.Top

	prevSibling := node.PreviousSibling()
	if prevSibling != nil {
		node.RenderBox.Left = prevSibling.RenderBox.Left + prevSibling.RenderBox.Width
	} else {
		node.RenderBox.Left = node.Parent.RenderBox.Left
	}
}

func paintInlineNode(ctx *gg.Context, node *structs.NodeDOM) {
	content := strings.TrimSpace(node.Content)

	if len(content) > 0 {
		ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
		ctx.SetRGBA(node.Parent.Style.BackgroundColor.R, node.Parent.Style.BackgroundColor.G, node.Parent.Style.BackgroundColor.B, node.Parent.Style.BackgroundColor.A)
		ctx.Fill()

		ctx.SetFont(sansSerif[node.Parent.Style.FontWeight], node.Parent.Style.FontSize)
		ctx.SetRGBA(node.Parent.Style.Color.R, node.Parent.Style.Color.G, node.Parent.Style.Color.B, node.Parent.Style.Color.A)
		ctx.DrawString(node.Content, node.RenderBox.Left, node.RenderBox.Height+node.RenderBox.Top)
	}
}

func calcBlockDimension(ctx *gg.Context, node *structs.NodeDOM) {
	node.RenderBox.Width = node.Parent.RenderBox.Width
}

func calcBlockPosition(ctx *gg.Context, node *structs.NodeDOM) {
	prevSibling := node.PreviousRealSibling()
	if prevSibling != nil {
		node.RenderBox.Top = prevSibling.RenderBox.Top + prevSibling.RenderBox.Height
	}
}

func paintBlockNode(ctx *gg.Context, node *structs.NodeDOM) {
	ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()
}

func paintNode(ctx *gg.Context, node *structs.NodeDOM) {
	switch node.Style.Display {
	case "block":
		paintBlockNode(ctx, node)
	case "inline":
		paintInlineNode(ctx, node)
	}

	for _, child := range node.Children {
		paintNode(ctx, child)
	}
}

func spaintNode(ctx *gg.Context, node *structs.NodeDOM) {
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
