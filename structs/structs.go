package structs

import (
	"thdwb/mustard"
	profiler "thdwb/profiler"
)

type WebBrowser struct {
	Document       *HTMLDocument
	ActiveDocument *Document
	Documents      []*Document
	Viewport       *mustard.CanvasWidget
}

type HTMLDocument struct {
	Title       string
	RootElement *NodeDOM
	URL         string
	RawDocument string
	OffsetY     int
	Styles      []*StyleElement
	Profiler    *profiler.Profiler

	SelectedElement *NodeDOM
	DebugFlag       bool
}

type Document struct {
	Title       string
	Path        string
	ContentType string
	RawDocument string
	DOM         *NodeDOM
}

type RenderBox struct {
	Node *NodeDOM

	Top  float64
	Left float64

	Width  float64
	Height float64

	MarginTop    float64
	MarginLeft   float64
	MarginRight  float64
	MarginBottom float64

	PaddingTop    float64
	PaddingLeft   float64
	PaddingRight  float64
	PaddingBottom float64
}

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

//Resource "HTTP resource struct definition"
type Resource struct {
	Body        string
	ContentType string
	Code        int
}

//Attribute "Generic key:value attribute definition"
type Attribute struct {
	Name  string
	Value string
}

//Stylesheet "Stylesheet definition for DOM Nodes"
type Stylesheet struct {
	Color           *ColorRGBA
	BackgroundColor *ColorRGBA

	FontSize float64
	Display  string
	Position string

	Width  float64
	Height float64
	Top    float64
	Left   float64
}

//StyleElement "hmtl <style> element"
type StyleElement struct {
	Selector string
	Style    *Stylesheet
}

//ColorRGBA "RGBA color model"
type ColorRGBA struct {
	R float64
	G float64
	B float64
	A float64
}
