package mustard

import (
	gg "github.com/danfragoso/thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateImageWidget - Creates and returns a new Image Widget
func CreateCanvasWidget(renderer func(*CanvasWidget)) *CanvasWidget {
	var widgets []Widget

	return &CanvasWidget{
		baseWidget: baseWidget{
			needsRepaint: true,
			widgets:      widgets,

			widgetType: canvasWidget,

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",
		},

		context:        gg.NewContext(0, 0),
		drawingContext: gg.NewContext(0, 0),
		renderer:       renderer,
		drawingRepaint: true,
	}
}

//SetWidth - Sets the label width
func (canvas *CanvasWidget) SetWidth(width float64) {
	canvas.box.width = width
	canvas.fixedWidth = true
	canvas.RequestReflow()
}

//SetHeight - Sets the label height
func (canvas *CanvasWidget) SetHeight(height float64) {
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

func (canvas *CanvasWidget) SetOffset(offset int) {
	canvas.offset = offset
}

func (canvas *CanvasWidget) GetOffset() int {
	return canvas.offset
}

func (canvas *CanvasWidget) GetContext() *gg.Context {
	return canvas.drawingContext
}

func (canvas *CanvasWidget) SetContext(ctx *gg.Context) {
	canvas.drawingContext = ctx
}

func (cavas *CanvasWidget) SetDrawingRepaint(repaint bool) {
	cavas.drawingRepaint = repaint
}

func (ctx *CanvasWidget) draw() {
	context := ctx.window.context
	top, left, width, height := ctx.computedBox.GetCoords()
	currentContextSize := ctx.context.Image().Bounds().Size()

	if currentContextSize.X != int(width) || currentContextSize.Y != int(height) {
		ctx.context = gg.NewContext(int(width), int(height))
		ctx.drawingContext = gg.NewContext(int(width), 12000)
		ctx.drawingRepaint = true
	}

	if ctx.drawingRepaint {
		ctx.renderer(ctx)
		ctx.drawingRepaint = false
	}

	ctx.context.DrawImage(ctx.drawingContext.Image(), int(left), ctx.offset)
	context.DrawImage(ctx.context.Image(), int(left), int(top))
	copyWidgetToBuffer(ctx, context.Image())
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
