package mustard

import (
	"image"

	gg "../gg"

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

	needsReflow bool
	visible     bool

	glw     *glfw.Window
	context *gg.Context
	backend *glBackend

	defaultCursor     *glfw.Cursor
	registeredButtons []*ButtonWidget
	registeredInputs  []*InputWidget
	activeInput       *InputWidget
	rootFrame         *Frame

	cursorX float64
	cursorY float64
}

type glBackend struct {
	program uint32

	vao uint32
	vbo uint32

	texture uint32
	quad    []float32
}

type box struct {
	top    int
	left   int
	width  int
	height int
}

type widget struct {
	box         box
	computedBox box

	needsRepaint bool
	fixedWidth   bool
	fixedHeight  bool

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

type TextWidget struct {
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

type ButtonWidget struct {
	widget
	content string

	icon      image.Image
	fontSize  float64
	fontColor string
	selected  bool
	padding   float64
	onClick   func()
}

type InputWidget struct {
	widget

	value          string
	selected       bool
	active         bool
	padding        float64
	fontSize       float64
	fontColor      string
	returnCallback func()
}
