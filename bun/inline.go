package bun

import (
	gg "thdwb/gg"
	structs "thdwb/structs"
)

func paintInlineElement(ctx *gg.Context, node *structs.NodeDOM) {
	ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.SetFont(sansSerif[node.Style.FontWeight], node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.RenderBox.Left, node.RenderBox.Top, 0, 0, node.RenderBox.Width, 1.5, gg.AlignLeft)
	ctx.Fill()
}

func calculateInlineLayout(ctx *gg.Context, node *structs.NodeDOM, childIdx int) {
	ctx.SetFont(sansSerif[node.Style.FontWeight], node.Style.FontSize)

	if childIdx > 0 && node.Parent.Children[childIdx-1] != nil {
		prev := node.Parent.Children[childIdx-1]
		if prev.Style.Display == "inline" {
			node.RenderBox.Top = prev.RenderBox.Top
			node.RenderBox.Left = prev.RenderBox.Left + prev.RenderBox.Width
		} else {
			node.RenderBox.Top = prev.RenderBox.Top + prev.RenderBox.Height
			node.RenderBox.Left = node.Parent.RenderBox.Left
		}
	} else {
		node.RenderBox.Top = node.Parent.RenderBox.Top
		node.RenderBox.Left = node.Parent.RenderBox.Left
	}

	node.RenderBox.Width, node.RenderBox.Height = ctx.MeasureMultilineString(node.Content, 1.5)
	node.RenderBox.Height++
}
