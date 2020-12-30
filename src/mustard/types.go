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

	hiDPI bool

	needsReflow bool
	visible     bool
	// This flag is active when drawing is happening on a thread that is not the
	// main one; When the asyncFlag is active the window frame processor function
	// should pool the status of the drawing routine and when it ends reflow? the
	// entire window surface?
	asyncFlag bool

	glw         *glfw.Window
	context     *gg.Context
	backend     *glBackend
	frameBuffer *image.RGBA

	defaultCursor  *glfw.Cursor
	pointerCursor  *glfw.Cursor
	selectedWidget Widget

	registeredTrees   []*TreeWidget
	registeredButtons []*ButtonWidget
	registeredInputs  []*InputWidget
	activeInput       *InputWidget
	rootFrame         *Frame

	cursorX float64
	cursorY float64

	pointerPositionEventListeners []func(float64, float64)
	scrollEventListeners          []func(int)
	clickEventListeners           []func(MustardKey)

	overlays         []*Overlay
	hasActiveOverlay bool

	staticOverlays   []*Overlay
	hasStaticOverlay bool

	contextMenu *contextMenu
}

type Overlay struct {
	ref string

	active bool

	top  float64
	left float64

	width  float64
	height float64

	position image.Point

	backgroundBuffer *image.RGBA
	buffer           *image.RGBA
}

type contextMenu struct {
	overlay       *Overlay
	entries       []*menuEntry
	selectedEntry *menuEntry
}

type menuEntry struct {
	entryText string
	action    func()

	top  float64
	left float64

	width  float64
	height float64
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

type Widget interface {
	Buffer() *image.RGBA
	SetNeedsRepaint(bool)
	NeedsRepaint() bool
	Widgets() []Widget
	ComputedBox() *box
	SetWindow(*Window)
	BaseWidget() *baseWidget

	draw()
}

type baseWidget struct {
	box            box
	computedBox    box
	widgetPosition widgetPosition

	font *truetype.Font

	needsRepaint bool
	fixedWidth   bool
	fixedHeight  bool

	widgets []Widget

	backgroundColor string

	widgetType widgetType
	cursor     *glfw.Cursor

	focusable  bool
	selectable bool

	focused  bool
	selected bool

	buffer *image.RGBA
	window *Window
}

type cursorType int

const (
	//DefaultCursor - Default arrow cursor
	DefaultCursor cursorType = iota
	//PointerCursor - Pointer cursor
	PointerCursor
)

type widgetType int

const (
	buttonWidget widgetType = iota
	canvasWidget
	frameWidget
	imageWidget
	inputWidget
	labelWidget
	treeWidget
	scrollbarWidget
	textWidget
)

type widgetPosition int

const (
	PositionRelative widgetPosition = iota
	PositionAbsolute
)

type FrameOrientation int

const (
	//VerticalFrame - Vertical frame orientation
	VerticalFrame FrameOrientation = iota

	//HorizontalFrame - Horizontal frame orientation
	HorizontalFrame
)

type MustardKey int

const (
	MouseLeft MustardKey = iota
	MouseRight
)

//Frame - Layout frame type
type Frame struct {
	baseWidget

	orientation FrameOrientation
}

type LabelWidget struct {
	baseWidget
	content string

	fontSize  float64
	fontColor string
}

type TreeWidget struct {
	baseWidget

	fontSize  float64
	fontColor string
	nodes     []*TreeWidgetNode

	openIcon  image.Image
	closeIcon image.Image
}

func (widget *TreeWidget) RemoveNodes() {
	widget.nodes = nil
}

func (widget *TreeWidget) AddNode(childNode *TreeWidgetNode) {
	widget.nodes = append(widget.nodes, childNode)
}

func CreateTreeWidgetNode(content string) *TreeWidgetNode {
	return &TreeWidgetNode{
		Content: content,
		box:     box{},
	}
}

type TreeWidgetNode struct {
	Content  string
	Parent   *TreeWidgetNode
	Children []*TreeWidgetNode

	isOpen bool
	index  int
	box    box
}

func (node *TreeWidgetNode) Toggle() {
	if node.isOpen {
		node.isOpen = false
	} else {
		node.isOpen = true
	}
}

func (node *TreeWidgetNode) Close() {
	node.isOpen = false

}
func (node *TreeWidgetNode) Open() {
	node.isOpen = true
}

func (node *TreeWidgetNode) AddNode(childNode *TreeWidgetNode) {
	childNode.Parent = node
	childNode.index = len(node.Children)
	node.Children = append(node.Children, childNode)
}

func (node *TreeWidgetNode) NextSibling() *TreeWidgetNode {
	selfIdx := node.index
	if selfIdx+1 < len(node.Parent.Children) {
		return node.Parent.Children[selfIdx+1]
	}

	return nil
}

func (node *TreeWidgetNode) PreviousSibling() *TreeWidgetNode {
	selfIdx := node.index
	if selfIdx-1 >= 0 {
		return node.Parent.Children[selfIdx-1]
	}

	return nil
}

type TextWidget struct {
	baseWidget
	content string

	fontSize  float64
	fontColor string
}

type ImageWidget struct {
	baseWidget

	path string
	img  image.Image
}

type CanvasWidget struct {
	baseWidget

	context        *gg.Context
	drawingContext *gg.Context

	renderer func(*CanvasWidget)

	scrollable bool
	offset     int

	drawingRepaint bool
}

type ButtonWidget struct {
	baseWidget
	content string

	icon      image.Image
	fontSize  float64
	fontColor string
	selected  bool
	padding   float64
	onClick   func()
}

type InputWidget struct {
	baseWidget

	value           string
	selected        bool
	active          bool
	padding         float64
	fontSize        float64
	context         *gg.Context
	fontColor       string
	cursorFloat     bool
	cursorPosition  int
	cursorDirection bool
	returnCallback  func()
}

type ScrollBarWidget struct {
	baseWidget

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
