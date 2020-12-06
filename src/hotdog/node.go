package hotdog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

//NodeDOM "DOM Node Struct definition"
type NodeDOM struct {
	Element string `json:"element"`
	Content string `json:"content"`

	Children   []*NodeDOM   `json:"children"`
	Attributes []*Attribute `json:"attributes"`
	Style      *Stylesheet  `json:"style"`
	Parent     *NodeDOM     `json:"-"`
	RenderBox  *RenderBox   `json:"-"`

	NeedsReflow  bool `json:"-"`
	NeedsRepaint bool `json:"-"`

	Document *Document `json:"-"`
}

func (node *NodeDOM) Print(d int) {
	spacing := strings.Repeat("-", d)
	fmt.Printf("|%s> %s\n", spacing, node.Element)

	for _, child := range node.Children {
		child.Print(d + 1)
	}
}

func (node *NodeDOM) JSON() string {
	res, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "{}"
	}

	return string(res)
}

func (node *NodeDOM) FindChildByName(childName string) (*NodeDOM, error) {
	if node.Element == childName {
		return node, nil
	}

	for _, child := range node.Children {
		foundChild, err := child.FindChildByName(childName)
		if err != nil {
			var noChild NoSuchElementError
			if errors.As(err, &noChild) {
				// No child with that element name, continue in other branches of the element tree
				continue
			}

			// Some other error
			return nil, err
		}

		return foundChild, nil
	}

	return nil, NoSuchElementError(childName)
}

func (node *NodeDOM) Attr(attrName string) string {
	for _, attribute := range node.Attributes {
		if attribute.Name == attrName {
			return attribute.Value
		}
	}

	return ""
}

func (node *NodeDOM) CalcPointIntersection(x, y float64) *NodeDOM {
	var intersectedNode *NodeDOM
	if x > float64(node.RenderBox.Left) &&
		x < float64(node.RenderBox.Left+node.RenderBox.Width) &&
		y > float64(node.RenderBox.Top) &&
		y < float64(node.RenderBox.Top+node.RenderBox.Height) {
		intersectedNode = node
	}

	for i := 0; i < len(node.Children); i++ {
		tempNode := node.Children[i].CalcPointIntersection(x, y)
		if tempNode != nil {
			intersectedNode = tempNode
		}
	}

	return intersectedNode
}

func (node NodeDOM) RequestRepaint() {
	node.NeedsRepaint = true

	for _, childNode := range node.Children {
		childNode.RequestRepaint()
	}
}

func (node NodeDOM) RequestReflow() {
	node.NeedsReflow = true

	for _, childNode := range node.Children {
		childNode.RequestReflow()
	}
}
