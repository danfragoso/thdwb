package structs

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
)

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
	Body string
	Code int
}

//Attribute "Generic key:value attribute definition"
type Attribute struct {
	Name  string
	Value string
}

//Stylesheet "Stylesheet definition for DOM Nodes"
type Stylesheet struct {
	Color    *ColorRGBA
	FontSize float64
	Display  string
}

//ColorRGBA "RGBA color model"
type ColorRGBA struct {
	R float64
	G float64
	B float64
	A float64
}

//UIElement "User interface elements"
type UIElement struct {
	ID    string
	WType string

	X float64
	Y float64
	W float64
	H float64

	Canvas *canvas.Canvas
	Cursor *glfw.Cursor

	Focused  bool
	Selected bool

	Text string
}

//Redraw "UIElement redraw function"
func (el UIElement) Redraw() {
	switch elementType := el.WType; elementType {
	case "input":
		DrawInput(el)
	case "box":
		DrawBox(el)
	}
}

//DrawBox "Draws the box element on it`s Canvas"
func DrawBox(el UIElement) {
	el.Canvas.SetFillStyle("#E1E2E1")
	el.Canvas.FillRect(el.X, el.Y, el.W, el.H)
}

//DrawInput "Draws the input element on it`s Canvas"
func DrawInput(el UIElement) {
	if el.Focused {
		el.Canvas.SetFillStyle("#C0C0C0")
	} else {
		el.Canvas.SetFillStyle("#CCC")
	}

	if el.Selected {
		el.Canvas.SetFillStyle("#FFF")
	}

	el.Canvas.FillRect(el.X, el.Y, el.W, el.H)

	if el.Text != "" {
		el.Canvas.SetFillStyle("#403F40")
		el.Canvas.SetFont("roboto.ttf", 16)
		el.Canvas.FillText(el.Text, el.X+5, el.Y+14+el.Y/2)
	}
}

//AppWindow "MustardUi Application Window Struct"
type AppWindow struct {
	Initialized bool

	Width  int
	Height int

	ViewportWidth  int
	ViewportHeight int

	AddressbarWidth  int
	AddressbarHeight int

	CursorX float64
	CursorY float64

	DefaultCursor *glfw.Cursor

	Title  string
	Redraw bool
	Resize bool

	ViewportOffset int

	Addressbar *canvas.Canvas
	Viewport   *canvas.Canvas

	AddressbarBackend *goglbackend.GoGLBackend
	ViewportBackend   *goglbackend.GoGLBackend
	GlfwWindow        *glfw.Window

	Location string

	UIElements []*UIElement
	DOM        *NodeDOM
}
