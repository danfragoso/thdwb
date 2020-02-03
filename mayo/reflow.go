package mayo

import (
	structs "../structs"
)

func ReflowNode(node *structs.NodeDOM, prev *structs.NodeDOM, siblingsOffset float64) float64 {
	for i := 0; i < len(node.Children); i++ {
		siblingsOffset += ReflowNode(node.Children[i], node.Children[i], siblingsOffset)
	}

	node.Style.Top = node.Style.Height + siblingsOffset

	if node == prev {
		return node.Style.Height
	}

	return node.Style.Height
}
