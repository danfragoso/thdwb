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
}

//SetHeight - Sets the label height
func (context *ContextWidget) SetHeight(height int) {
	context.box.height = height
	context.fixedHeight = true
}

func (context *ContextWidget) GetContext() *gg.Context {
	return context.context
}

func drawContextWidget(context *gg.Context, widget *ContextWidget, top, left, width, height int) {
	if widget.context == nil || widget.context.Width() != width || widget.context.Height() != height {
		widget.context = gg.NewContext(width, height)

		widget.context.SetRGB(1, 1, 1)
		widget.context.Clear()

		widget.renderer(widget.context)
	}

	context.DrawImage(widget.context.Image(), left, top)
}
