package bun

import (
	"bytes"
	"fmt"
	"github.com/danfragoso/thdwb/assets"
	gg "github.com/danfragoso/thdwb/gg"
	hotdog "github.com/danfragoso/thdwb/hotdog"
	"github.com/danfragoso/thdwb/sauce"
	"image"
)

func paintInlineElement(ctx *gg.Context, node *hotdog.NodeDOM) {
	ctx.DrawRectangle(node.RenderBox.Left, node.RenderBox.Top, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(node.Style.BackgroundColor.R, node.Style.BackgroundColor.G, node.Style.BackgroundColor.B, node.Style.BackgroundColor.A)
	ctx.Fill()

	if node.Element == "img" {
		// We should not download the image again, as we already have this saved.
		// We could save the image on the node or have it cached on the
		im, err := fetchNodeImage(node)

		if err != nil {
			fmt.Println(err)
			// Use the stand-in error image.
			im, _, _ = image.Decode(bytes.NewReader(assets.ErrorImage()))
		}
		ctx.DrawImage(im, int(node.RenderBox.Left), int(node.RenderBox.Top))
	}

	ctx.SetRGBA(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B, node.Style.Color.A)
	ctx.SetFont(sansSerif[node.Style.FontWeight], node.Style.FontSize)
	ctx.DrawStringWrapped(node.Content, node.RenderBox.Left, node.RenderBox.Top, 0, 0, node.RenderBox.Width, 1, gg.AlignLeft)
	ctx.Fill()
}

func fetchNodeImage(node *hotdog.NodeDOM) (image.Image, error) {
	imgPath := node.Attr("src")

	imgURL, err := node.Document.URL.Parse(imgPath)
	if err != nil {
		return nil, err
	}

	// Fetch the image (either from cache or network)
	data, err := sauce.GetImage(imgURL)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}
	return im, nil
}

func calculateInlineLayout(ctx *gg.Context, node *hotdog.NodeDOM, childIdx int) {
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

	if node.Element == "img" {
		im, err := fetchNodeImage(node)
		if err != nil {
			fmt.Println(err)
			// Use the stand-in error image.
			im, _, _ = image.Decode(bytes.NewReader(assets.ErrorImage()))
		}
		imgSize := im.Bounds().Size()

		node.RenderBox.Width = float64(imgSize.X)
		node.RenderBox.Height = float64(imgSize.Y)
	} else {
		if node.RenderBox.Width == 0 {
			node.RenderBox.Width = node.Parent.RenderBox.Width
		}

		node.RenderBox.Height = ctx.MeasureStringWrapped(node.Content, node.RenderBox.Width, 1)
		mW, _ := ctx.MeasureString(node.Content)
		if mW < node.RenderBox.Width {
			node.RenderBox.Width = mW
		}
	}

	node.RenderBox.Height++
}
