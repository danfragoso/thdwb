package mustard

import (
	"image"

	gg "thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/goki/freetype/truetype"
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
	hasOverlay  bool
	// This flag is active when drawing is happening on a thread that is not the
	// main one; When the asyncFlag is active the window frame processor function
	// should pool the status of the drawing routine and when it ends reflow? the
	// entire window surface?
	asyncFlag bool

	glw            *glfw.Window
	context        *gg.Context
	overlayContext *gg.Context
	backend        *glBackend
	frameBuffer    *image.RGBA

	defaultCursor *glfw.Cursor
	pointerCursor *glfw.Cursor

	registeredButtons []*ButtonWidget
	registeredInputs  []*InputWidget
	activeInput       *InputWidget
	rootFrame         *Frame

	cursorX float64
	cursorY float64

	pointerPositionEventListeners []func(float64, float64)
	scrollEventListeners          []func(int)
	clickEventListeners           []func()

	menuEntries []*menuEntry
}

type menuEntry struct {
	entryText string
	action    func()
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

func (box *box) SetCoords(top, left, width, height int) {
	box.top = top
	box.left = left
	box.width = width
	box.height = height
}

func (box *box) GetCoords() (int, int, int, int) {
	return box.top, box.left, box.width, box.height
}

type widget struct {
	box         box
	computedBox box

	font *truetype.Font

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

	window *Window
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

type CanvasWidget struct {
	widget

	context        *gg.Context
	drawingContext *gg.Context

	renderer func(*CanvasWidget)

	scrollable bool
	offset     int

	drawingRepaint bool
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
	context        *gg.Context
	fontColor      string
	cursorPosition int
	returnCallback func()
}

type ScrollBarWidget struct {
	widget

	orientation ScrollBarOrientation
	selected    bool
	thumbSize   float64
	thumbColor  string

	scrollerSize   float64
	scrollerOffset float64
}

type ScrollBarOrientation int

const (
	VerticalScrollBar ScrollBarOrientation = iota
	HorizontalScrollBar
)
