package mustard

import (
	"fmt"

	"github.com/danfragoso/thdwb/structs"
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

func walkDOM(TreeDOM *structs.NodeDOM, d int) {
	fmt.Println(d, getElementName(TreeDOM))
	nodeChildren := getNodeChildren(TreeDOM)

	for i := 0; i < len(nodeChildren); i++ {
		walkDOM(nodeChildren[i], d+1)
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

func renderNode(node *structs.NodeDOM, browserWindow *structs.AppWindow, vOffset float64) {
	sizeStep := node.Style.FontSize

	if node.Style.Display == "block" {
		if node.Style.Color != nil {
			browserWindow.Viewport.SetFillStyle(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B)
		} else {
			browserWindow.Viewport.SetFillStyle("#000")
		}

		browserWindow.Viewport.SetFont("roboto.ttf", sizeStep)
		browserWindow.Viewport.FillText(node.Content, 0, vOffset+node.Style.Top)
	}

	children := getNodeChildren(node)

	for i := 0; i < len(children); i++ {
		if isNodeInsideViewportBounds(browserWindow, children[i], vOffset+children[i].Style.Top) {
			renderNode(children[i], browserWindow, vOffset+children[i].Style.Top)
		}
	}
}

func isNodeInsideViewportBounds(browserWindow *structs.AppWindow, node *structs.NodeDOM, vOffset float64) bool {

	if vOffset > float64(browserWindow.ViewportHeight) {
		return false
	}

	return true
}
