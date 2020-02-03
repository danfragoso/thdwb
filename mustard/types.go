package mustard

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type App struct {
	name    string
	windows []*Window
}

type Window struct {
	title  string
	width  int
	height int

	isDirty bool
	visible bool

	glw     *glfw.Window
	context *gg.Context
	backend *glBackend

	rootFrame *Frame
}

type glBackend struct {
	program uint32

	vao uint32
	vbo uint32

	texture uint32
	quad    []float32
}

type widget struct {
	top    int
	left   int
	width  int
	height int

	dirty       bool
	fixedWidth  bool
	fixedHeight bool

	widgets []interface{}

	backgroundColor string

	ref    string
	cursor *glfw.Cursor

	focusable  bool
	selectable bool

	focused  bool
	selected bool
}

type FrameOrientation int

const (
	//VerticalFrame - Vertical frame orientation
	VerticalFrame FrameOrientation = iota

	//HorizontalFrame - Horizontal frame orientation
	HorizontalFrame
)

//Frame - Layout frame type
type Frame struct {
	widget

	orientation FrameOrientation
}

type LabelWidget struct {
	widget
	content string

	fontSize  float64
	fontColor string
}

type ImageWidget struct {
	widget

	path string
	img  image.Image
}

type ContextWidget struct {
	widget

	context  *gg.Context
	renderer func(*gg.Context)
}
