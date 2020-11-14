package mustard

import (
	"image"
	"image/draw"
	gg "thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateImageWidget - Creates and returns a new Image Widget
func CreateCanvasWidget(renderer func(*CanvasWidget)) *CanvasWidget {
	var widgets []Widget

	return &CanvasWidget{
		baseWidget: baseWidget{
			needsRepaint: true,
			widgets:      widgets,

			ref: "context",

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

	if currentContextSize.X != width || currentContextSize.Y != height {
		ctx.context = gg.NewContext(width, height)
		ctx.drawingContext = gg.NewContext(width, 12000)
		ctx.drawingRepaint = true
	}

	if ctx.drawingRepaint {
		ctx.renderer(ctx)
		ctx.drawingRepaint = false
	}

	ctx.context.DrawImage(ctx.drawingContext.Image(), left, ctx.offset)
	context.DrawImage(ctx.context.Image(), left, top)

	if ctx.buffer == nil || ctx.buffer.Bounds().Max.X != width && ctx.buffer.Bounds().Max.Y != height {
		ctx.buffer = image.NewRGBA(image.Rectangle{
			image.Point{}, image.Point{width, height},
		})
	}

	draw.Draw(ctx.buffer, image.Rectangle{
		image.Point{},
		image.Point{width, height},
	}, context.Image(), image.Point{left, top}, draw.Over)
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
