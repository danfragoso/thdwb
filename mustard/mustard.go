package mustard

import (
	"fmt"

	"github.com/danfragoso/thdwb/structs"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func getNodeContent(NodeDOM *structs.NodeDOM) string {
	return NodeDOM.Content
}

func getElementName(NodeDOM *structs.NodeDOM) string {
	return NodeDOM.Element
}

func getNodeChildren(NodeDOM *structs.NodeDOM) []*structs.NodeDOM {
	return NodeDOM.Children
}

func walkDOM(DOM_Tree *structs.NodeDOM, d int) {
	fmt.Println(d, getElementName(DOM_Tree))
	nodeChildren := getNodeChildren(DOM_Tree)

	for i := 0; i < len(nodeChildren); i++ {
		walkDOM(nodeChildren[i], d+1)
	}
}

func renderNode(NodeDOM *structs.NodeDOM, cr *cairo.Context, x float64, y float64) {
	nodeChildren := getNodeChildren(NodeDOM)
	sizeStep := NodeDOM.Style.FontSize

	if NodeDOM.Style.Display == "block" {
		cr.SetSourceRGB(NodeDOM.Style.Color.R, NodeDOM.Style.Color.G, NodeDOM.Style.Color.B)
		cr.SelectFontFace("Arial", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
		cr.SetFontSize(sizeStep)
		cr.Translate(x, y+sizeStep+2)
		cr.ShowText(getNodeContent(NodeDOM))
		cr.Translate(0, 2)
		cr.Fill()
	}

	for i := 0; i < len(nodeChildren); i++ {
		renderNode(nodeChildren[i], cr, x, y*float64(i))
	}
}

func GetPageTitle(DOM_Tree *structs.NodeDOM) string {
	nodeChildren := getNodeChildren(DOM_Tree)
	pageTitle := "Sem Titulo"

	if getElementName(DOM_Tree) == "title" {
		return getNodeContent(DOM_Tree)
	} else {
		for i := 0; i < len(nodeChildren); i++ {
			nPageTitle := GetPageTitle(nodeChildren[i])

			if nPageTitle != "Sem Titulo" {
				pageTitle = nPageTitle
			}
		}
	}

	return pageTitle
}

func DrawDOM(DOM_Tree *structs.NodeDOM) func(drawingArea *gtk.DrawingArea, cr *cairo.Context) {
	return func(drawingArea *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(1, 1, 1)
		cr.Paint()
		renderNode(DOM_Tree, cr, 0, 0)
		cr.IdentityMatrix()
	}
}
