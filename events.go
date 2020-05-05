package main

import (
	structs "thdwb/structs"
)

func processPointerPosition(browser *structs.WebBrowser, x, y float64) {
	browser.SelectedNode = nil
	browser.SelectedNode = calcIntersection(browser.Document.RootElement, x, y)
}

func calcIntersection(node *structs.NodeDOM, x, y float64) *structs.NodeDOM {
	var intersectedNode *structs.NodeDOM
	if x > float64(node.RenderBox.Left) &&
		x < float64(node.RenderBox.Left+node.RenderBox.Width) &&
		y > float64(node.RenderBox.Top) &&
		y < float64(node.RenderBox.Top+node.RenderBox.Height) {
		intersectedNode = node
	}

	for i := 0; i < len(node.Children); i++ {
		tempNode := calcIntersection(node.Children[i], x, y)
		if tempNode != nil {
			intersectedNode = tempNode
		}
	}

	return intersectedNode
}
