package mustard

import (
	"fmt"
	"log"

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
	fmt.Println(d, getNodeContent(DOM_Tree))
	nodeChildren := getNodeChildren(DOM_Tree)

	for i := 0; i < len(nodeChildren); i++ {
		walkDOM(nodeChildren[i], d+1)
	}
}

func renderH2(content string, cr *cairo.Context, x float64, y float64) {
	sizeStep := float64(28)
	cr.SetSourceRGB(0, 0, 0)
	cr.SelectFontFace("Arial", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(sizeStep)
	cr.Translate(x, y+sizeStep+2)
	cr.ShowText(content)
	cr.Translate(0, 10)
	cr.Fill()
}

func renderH1(content string, cr *cairo.Context, x float64, y float64) {
	sizeStep := float64(36)
	cr.SetSourceRGB(1, 0, 0)
	cr.SelectFontFace("Arial", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(sizeStep)
	cr.Translate(x, y+sizeStep+2)
	cr.ShowText(content)
	cr.Translate(0, 10)
	cr.Fill()
}

func renderP(content string, cr *cairo.Context, x float64, y float64) {
	sizeStep := float64(14)
	cr.SetSourceRGB(0, 0, 0)
	cr.SelectFontFace("Arial", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	cr.SetFontSize(sizeStep)
	cr.Translate(x, y+sizeStep+5)
	cr.ShowText(content)
	cr.Translate(0, 5)
	cr.Fill()
}

func renderB(content string, cr *cairo.Context, x float64, y float64) {
	sizeStep := float64(14)
	cr.SetSourceRGB(0, 0, 0)
	cr.SelectFontFace("Arial", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(sizeStep)
	cr.Translate(x, y+sizeStep+1)
	cr.ShowText(content)
	cr.Translate(0, 1)
	cr.Fill()
}

func renderI(content string, cr *cairo.Context, x float64, y float64) {
	sizeStep := float64(14)
	cr.SetSourceRGB(0, 0, 0)
	cr.SelectFontFace("Arial", cairo.FONT_SLANT_ITALIC, cairo.FONT_WEIGHT_NORMAL)
	cr.SetFontSize(sizeStep)
	cr.Translate(x, y+sizeStep+1)
	cr.ShowText(content)
	cr.Translate(0, 1)
	cr.Fill()
}

func renderDefault(content string, cr *cairo.Context, x float64, y float64) {
	sizeStep := float64(14)
	cr.SetSourceRGB(0, 0, 0)
	cr.SelectFontFace("Arial", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_NORMAL)
	cr.SetFontSize(sizeStep)
	cr.Translate(x, y+sizeStep+2)
	cr.ShowText(content)
	cr.Translate(0, 2)
	cr.Fill()
}

func renderNode(NodeDOM *structs.NodeDOM, cr *cairo.Context, x float64, y float64) {
	nodeChildren := getNodeChildren(NodeDOM)

	switch el := getElementName(NodeDOM); el {
	case "h1":
		renderH1(getNodeContent(NodeDOM), cr, x, y)
	case "h2":
		renderH2(getNodeContent(NodeDOM), cr, x, y)
	case "p":
		renderP(getNodeContent(NodeDOM), cr, x, y)
	case "i":
		renderI(getNodeContent(NodeDOM), cr, x, y)
	case "b":
		renderB(getNodeContent(NodeDOM), cr, x, y)
	default:
		renderDefault(getNodeContent(NodeDOM), cr, x, y)
	}

	for i := 0; i < len(nodeChildren); i++ {
		renderNode(nodeChildren[i], cr, x, y*float64(i))
	}
}

func getPageTitle(DOM_Tree *structs.NodeDOM) string {
	if getElementName(DOM_Tree) == "title" {
		pageTitle := getNodeContent(DOM_Tree)

		if pageTitle == "" {
			return "Sem Título"
		} else {
			return pageTitle
		}
	} else {
		return getPageTitle(DOM_Tree.Children[0])
	}
}

func drawDOM(DOM_Tree *structs.NodeDOM) func(drawingArea *gtk.DrawingArea, cr *cairo.Context) {
	return func(drawingArea *gtk.DrawingArea, cr *cairo.Context) {
		renderNode(DOM_Tree, cr, 0, 0)
	}
}

func RenderDOM(DOM_Tree *structs.NodeDOM) {
	gtk.Init(nil)

	browserWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	drawingArea, _ := gtk.DrawingAreaNew()

	browserWindow.Add(drawingArea)

	header, err := gtk.HeaderBarNew()
	if err != nil {
		log.Fatal("Could not create header bar:", err)
	}

	header.SetShowCloseButton(true)
	header.SetTitle(getPageTitle(DOM_Tree) + " - THDWB")

	browserWindow.SetTitlebar(header)
	browserWindow.Connect("destroy", gtk.MainQuit)
	browserWindow.ShowAll()

	html := DOM_Tree.Children[0]
	drawingArea.Connect("draw", drawDOM(html.Children[1]))
	gtk.Main()
}
