package bun

import hotdog "thdwb/hotdog"

func createRenderTree(root *hotdog.NodeDOM) *hotdog.NodeDOM {
	if root.Style.Display == "none" {
		return nil
	}

	node := &hotdog.NodeDOM{
		Style:      root.Style,
		Element:    root.Element,
		Content:    root.Content,
		Attributes: root.Attributes,
	}

	node.RenderBox = &hotdog.RenderBox{}
	for _, child := range root.Children {
		r := createRenderTree(child)
		if r != nil {
			r.Parent = node
			node.Children = append(node.Children, r)
		}
	}

	return node
}
