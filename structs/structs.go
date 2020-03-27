package structs

import (
	profiler "thdwb/profiler"
)

type WebBrowser struct {
	Document *HTMLDocument
}

type HTMLDocument struct {
	Title       string
	RootElement *NodeDOM
	URL         string
	RawDocument string
	Styles      []*StyleElement
	Profiler    *profiler.Profiler
}

//NodeDOM "DOM Node Struct definition"
type NodeDOM struct {
	Element    string       `json:"element"`
	Content    string       `json:"content"`
	Children   []*NodeDOM   `json:"children"`
	Attributes []*Attribute `json:"attributes"`
	Style      *Stylesheet  `json:"style"`
	Parent     *NodeDOM     `json:"-"`
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
