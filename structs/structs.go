package structs

import (
	profiler "thdwb/profiler"
)

type WebBrowser struct {
	Document       *HTMLDocument
	ActiveDocument *Document
	Documents      []*Document
	SelectedNode   *NodeDOM
}

type HTMLDocument struct {
	Title       string
	RootElement *NodeDOM
	URL         string
	RawDocument string
	OffsetY     int
	Styles      []*StyleElement
	Profiler    *profiler.Profiler
	PointerXPos float64
	PointerYPos float64
	DebugFlag   bool
}

type Document struct {
	Title       string
	Path        string
	ContentType string
	RawDocument string
	DOM         *NodeDOM
}

type RenderBox struct {
	Node        *NodeDOM
	NeedsReflow bool

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
	Element    string       `json:"element"`
	Content    string       `json:"content"`
	Children   []*NodeDOM   `json:"children"`
	Attributes []*Attribute `json:"attributes"`
	Style      *Stylesheet  `json:"style"`
	Parent     *NodeDOM     `json:"-"`
	RenderBox  *RenderBox   `json:"-"`
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
