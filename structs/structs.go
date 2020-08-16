package structs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"thdwb/mustard"
	profiler "thdwb/profiler"
)

type WebBrowser struct {
	Document       *HTMLDocument
	ActiveDocument *Document
	Documents      []*Document
	Viewport       *mustard.CanvasWidget
	StatusLabel    *mustard.LabelWidget
	History        *History
	Window         *mustard.Window
}

type HTMLDocument struct {
	Title       string
	RootElement *NodeDOM
	RenderTree  *NodeDOM
	URL         *url.URL
	RawDocument string
	OffsetY     int
	Styles      []*StyleElement
	Profiler    *profiler.Profiler

	SelectedElement *NodeDOM
	DebugFlag       bool
}

type History struct {
	pages    []*url.URL
	allPages []*url.URL
}

func (history *History) AllPages() []*url.URL {
	return history.allPages
}

func (history *History) PageCount() int {
	return len(history.pages)
}

func (history *History) Push(URL *url.URL) {
	history.pages = append(history.pages, URL)
	history.allPages = append(history.allPages, URL)
}

func (history *History) Last() *url.URL {
	return history.pages[len(history.pages)-1]
}

func (history *History) Pop() {
	if len(history.pages) > 0 {
		history.pages = history.pages[:len(history.pages)-1]
	}
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

func (box *RenderBox) String() string {
	return fmt.Sprintf("[top: %f, left: %f, width: %f, height: %f]", box.Top, box.Left, box.Width, box.Height)
}

func (box *RenderBox) GetRect() (float64, float64, float64, float64) {
	return box.Top, box.Left, box.Width, box.Height
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

	Document *HTMLDocument `json:"-"`
}

func (node *NodeDOM) PreviousRealSibling() *NodeDOM {
	var nodeIdx int
	for idx, child := range node.Parent.Children {
		if child == node {
			nodeIdx = idx
			break
		}
	}

	for i := nodeIdx; i > 0; i-- {
		if node.Parent.Children[i-1] != nil && node.Parent.Children[i-1].Element != "html:text" {
			return node.Parent.Children[i-1]
		}
	}

	return nil
}

func (node *NodeDOM) NextRealSibling() *NodeDOM {
	for i, child := range node.Parent.Children {
		if child == node && i < len(node.Parent.Children) &&
			node.Parent.Children[i+1].Element != "html:text" {
			return node.Parent.Children[i+1]
		}
	}

	return nil
}

func (node *NodeDOM) PreviousSibling() *NodeDOM {
	var nodeIdx int
	for idx, child := range node.Parent.Children {
		if child == node {
			nodeIdx = idx
			break
		}
	}

	for i := nodeIdx; i > 0; i-- {
		if node.Parent.Children[i-1] != nil {
			return node.Parent.Children[i-1]
		}
	}

	return nil
}

func (node *NodeDOM) NextSibling() *NodeDOM {
	for i, child := range node.Parent.Children {
		if child == node && i < len(node.Parent.Children) {
			return node.Parent.Children[i+1]
		}
	}

	return nil
}

func (node *NodeDOM) FindChildByName(childName string) *NodeDOM {
	if node.Element == childName {
		return node
	}

	for _, child := range node.Children {
		foundChild := child.FindChildByName(childName)

		if foundChild != nil {
			return foundChild
		}
	}

	return nil
}

func (node *NodeDOM) Attr(attrName string) string {
	for _, attribute := range node.Attributes {
		if attribute.Name == attrName {
			return attribute.Value
		}
	}

	return ""
}

func (node *NodeDOM) JSON() string {
	res, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "{}"
	}

	return string(res)
}

func (node *NodeDOM) CalcPointIntersection(x, y float64) *NodeDOM {
	var intersectedNode *NodeDOM
	if node == nil {
		return intersectedNode
	}

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
	URL         *url.URL
	Key         string
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

	FontSize   float64
	FontWeight int

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

type ResourceCache struct {
	cachedResources []*Resource
}

func (cache *ResourceCache) AddResource(resource *Resource) {
	cache.cachedResources = append(cache.cachedResources, resource)
}

func (cache *ResourceCache) GetResource(resourceKey string) *Resource {
	for _, resource := range cache.cachedResources {
		if resource.Key == resourceKey {
			return resource
		}
	}

	return nil
}

type CachedImage struct {
	Key   string
	Image []byte
}
type ImgCache struct {
	cachedImages []*CachedImage
}

func (cache *ImgCache) AddImage(key string, value []byte) {
	cache.cachedImages = append(cache.cachedImages,
		&CachedImage{
			Key:   key,
			Image: value,
		},
	)
}

func (cache *ImgCache) GetImage(imageKey string) *CachedImage {
	for _, image := range cache.cachedImages {
		if image.Key == imageKey {
			return image
		}
	}

	return nil
}

func Log(component, msg string) {
	str := "(" + "\033[95m" + component + "\033[0m" + ")"
	fmt.Println(str, msg)
}
