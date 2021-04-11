package bun

import (
	gg "github.com/danfragoso/thdwb/gg"
	hotdog "github.com/danfragoso/thdwb/hotdog"
)

func paintListItemElement(ctx *gg.Context, node *hotdog.NodeDOM) {
	ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	ctx.DrawCircle(node.RenderBox.Left-15, node.RenderBox.Top+node.Style.FontSize/2, 3)
	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.SetFont(sansSerif[node.Style.FontWeight], node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.RenderBox.Left, node.RenderBox.Top+1, 0, 0, node.RenderBox.Width, 1.5, gg.AlignLeft)
	ctx.Fill()
}

func calculateListItemLayout(ctx *gg.Context, node *hotdog.NodeDOM, childIdx int) {
	if node.Style.Width == 0 {
		node.RenderBox.Width = node.Parent.RenderBox.Width - 30
	}

	if node.Style.Height == 0 && len(node.Content) > 0 {
		ctx.SetFont(sansSerif[node.Style.FontWeight], node.Style.FontSize)
		node.RenderBox.Height = ctx.MeasureStringWrapped(node.Content, node.RenderBox.Width, 1.5) + 2 + ctx.FontHeight()*.5
	}

	if childIdx > 0 {
		prev := node.Parent.Children[childIdx-1]
		node.RenderBox.Top = prev.RenderBox.Top + prev.RenderBox.Height
	} else {
		node.RenderBox.Top = node.Parent.RenderBox.Top
	}

	node.RenderBox.Left = 30
}
