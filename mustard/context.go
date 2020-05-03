package mustard

import (
	gg "thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateImageWidget - Creates and returns a new Image Widget
func CreateCanvasWidget(renderer func(*gg.Context)) *CanvasWidget {
	var widgets []interface{}

	return &CanvasWidget{
		widget: widget{
			needsRepaint: true,
			widgets:      widgets,

			ref: "context",

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",
		},

		context:  nil,
		renderer: renderer,
	}
}

//AttachWidget - Attaches a new widget to the window
func (canvas *CanvasWidget) AttachWidget(widget interface{}) {
	canvas.widgets = append(canvas.widgets, widget)
}

//SetWidth - Sets the label width
func (canvas *CanvasWidget) SetWidth(width int) {
	canvas.box.width = width
	canvas.fixedWidth = true
	canvas.RequestReflow()
}

//SetHeight - Sets the label height
func (canvas *CanvasWidget) SetHeight(height int) {
	canvas.box.height = height
	canvas.fixedHeight = true
	canvas.RequestReflow()
}

func (canvas *CanvasWidget) EnableScrolling() {
	canvas.scrollable = true
}

func (canvas *CanvasWidget) DisableScrolling() {
	canvas.scrollable = false
	canvas.offset = 0
}

func (canvas *CanvasWidget) GetContext() *gg.Context {
	return canvas.context
}

func (ctx *CanvasWidget) draw() {
	context := ctx.window.context
	top, left, width, height := ctx.computedBox.GetCoords()
	if ctx.context == nil || ctx.context.Width() != width || ctx.context.Height() != height {
		if ctx.scrollable {
			createCtxScrollBar(ctx)
			width -= 12
		}

		ctx.context = gg.NewContext(width, height)

		ctx.context.SetRGB(1, 1, 1)
		ctx.context.Clear()

		ctx.renderer(ctx.context)
	}

	context.DrawImage(ctx.context.Image(), left, top)
	ctx.needsRepaint = false
}

func createCtxScrollBar(ctx *CanvasWidget) {
	top, _, width, height := ctx.computedBox.GetCoords()
	context := ctx.window.context

	//Scroll Track
	context.SetHexColor("#c1c1c1")
	context.DrawRectangle(float64(width-12), float64(top), 12, float64(height))
	context.Fill()

	//Scroll Arrow
	context.SetHexColor("#ff0000")
	context.DrawRectangle(float64(width-12), 30, 10, 10)
	context.Fill()

	//Scroll Thumb
	context.SetHexColor("#565656")
	context.DrawRectangle(float64(width-12), float64(top), 12, 200)
	context.Fill()
}
