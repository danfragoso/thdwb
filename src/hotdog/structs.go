package hotdog

import (
	"fmt"
	"net/url"
	"thdwb/mustard"
	profiler "thdwb/profiler"
)

type WebBrowser struct {
	ActiveDocument *Document
	Documents      []*Document

	Viewport    *mustard.CanvasWidget
	StatusLabel *mustard.LabelWidget
	History     *History
	Window      *mustard.Window
	Profiler    *profiler.Profiler
	BuildInfo   *BuildInfo
	Settings    *Settings
}

type Document struct {
	Title       string
	ContentType string
	URL         *url.URL

	RawDocument string
	DOM         *NodeDOM

	DebugFlag       bool
	SelectedElement *NodeDOM

	OffsetY int
}
type BuildInfo struct {
	GitRevision string
	GitBranch   string

	HostInfo  string
	BuildTime string
}

type History struct {
	previousPages []*url.URL
	nextPages     []*url.URL
}

func (history *History) NextPages() []*url.URL {
	return history.nextPages
}

func (history *History) AllPages() []*url.URL {
	return history.previousPages
}

func (history *History) PageCount() int {
	return len(history.previousPages)
}

func (history *History) Push(URL *url.URL) {
	history.nextPages = nil
	history.previousPages = append(history.previousPages, URL)
}

func (history *History) Last() *url.URL {
	return history.previousPages[len(history.previousPages)-1]
}

func (history *History) PopNext() {
	if len(history.nextPages) > 0 {
		history.previousPages = append(history.previousPages, history.nextPages[len(history.nextPages)-1])
		history.nextPages = nil
	}
}

func (history *History) Pop() {
	if len(history.previousPages) > 0 {
		history.nextPages = append(history.nextPages, history.previousPages[len(history.previousPages)-1])
		history.previousPages = history.previousPages[:len(history.previousPages)-1]
	}
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

func (box *RenderBox) GetRect() (float64, float64, float64, float64) {
	return box.Top, box.Left, box.Width, box.Height
}

type NoSuchElementError string

func (e NoSuchElementError) Error() string {
	return fmt.Sprintf("no such element: %q", string(e))
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
