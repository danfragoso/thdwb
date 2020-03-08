package mustard

import (
	gg "../gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateImageWidget - Creates and returns a new Image Widget
func CreateContextWidget(renderer func(*gg.Context)) *ContextWidget {
	var widgets []interface{}

	return &ContextWidget{
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
func (context *ContextWidget) AttachWidget(widget interface{}) {
	context.widgets = append(context.widgets, widget)
}

//SetWidth - Sets the label width
func (context *ContextWidget) SetWidth(width int) {
	context.box.width = width
	context.fixedWidth = true
	context.RequestReflow()
}

//SetHeight - Sets the label height
func (context *ContextWidget) SetHeight(height int) {
	context.box.height = height
	context.fixedHeight = true
	context.RequestReflow()
}

func (context *ContextWidget) GetContext() *gg.Context {
	return context.context
}

func (ctx *ContextWidget) draw() {
	context := ctx.window.context
	top, left, width, height := ctx.computedBox.GetCoords()
	if ctx.context == nil || ctx.context.Width() != width || ctx.context.Height() != height {
		ctx.context = gg.NewContext(width, height)

		ctx.context.SetRGB(1, 1, 1)
		ctx.context.Clear()

		ctx.renderer(ctx.context)
	}

	context.DrawImage(ctx.context.Image(), left, top)
	ctx.needsRepaint = false
}
