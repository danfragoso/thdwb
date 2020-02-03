package mustard

import (
	"github.com/fogleman/gg"
)

//CreateFrame - Creates and returns a new Frame
func CreateFrame(orientation FrameOrientation) *Frame {
	var widgets []interface{}

	return &Frame{
		widget: widget{
			top:  0,
			left: 0,

			width:  100,
			height: 100,

			ref: "frame",

			dirty:   true,
			widgets: widgets,

			backgroundColor: "#fff"},

		orientation: orientation,
	}
}

//AttachWidget - Attaches widgets to a frame el
func (frame *Frame) AttachWidget(widget interface{}) {
	frame.widgets = append(frame.widgets, widget)
}

//SetBackgroundColor - Sets the frame background color
func (frame *Frame) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		frame.backgroundColor = backgroundColor
		frame.dirty = true
	}
}

//SetWidth - Sets the frame width
func (frame *Frame) SetWidth(width int) {
	frame.width = width
	frame.fixedWidth = true
}

//SetHeight - Sets the frame height
func (frame *Frame) SetHeight(height int) {
	frame.height = height
	frame.fixedHeight = true
}

func drawRootFrame(window *Window) {
	drawFrame(window.context, window.rootFrame, 0, 0, window.width, window.height)
}

func drawFrame(context *gg.Context, frame *Frame, top, left, width, height int) {
	context.SetHexColor(frame.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	childrenLen := len(frame.widgets)
	if childrenLen > 0 {
		childrenWidgets := getCoreWidgets(frame.widgets)
		childrenLayout := calculateChildrenWidgetsLayout(childrenWidgets, top, left, width, height, frame.orientation)

		for i := 0; i < childrenLen; i++ {
			switch frame.widgets[i].(type) {
			case *Frame:
				frame := frame.widgets[i].(*Frame)
				frame.top = childrenLayout[i].top
				frame.left = childrenLayout[i].left
				frame.width = childrenLayout[i].width
				frame.height = childrenLayout[i].height

				drawFrame(context, frame, childrenLayout[i].top, childrenLayout[i].left, childrenLayout[i].width, childrenLayout[i].height)
			case *LabelWidget:
				label := frame.widgets[i].(*LabelWidget)
				label.top = childrenLayout[i].top
				label.left = childrenLayout[i].left
				label.width = childrenLayout[i].width
				label.height = childrenLayout[i].height

				drawLabelWidget(context, label, childrenLayout[i].top, childrenLayout[i].left, childrenLayout[i].width, childrenLayout[i].height)
			case *ImageWidget:
				image := frame.widgets[i].(*ImageWidget)
				image.top = childrenLayout[i].top
				image.left = childrenLayout[i].left
				image.width = childrenLayout[i].width
				image.height = childrenLayout[i].height

				drawImageWidget(context, image, childrenLayout[i].top, childrenLayout[i].left, childrenLayout[i].width, childrenLayout[i].height)
			case *ContextWidget:
				ctx := frame.widgets[i].(*ContextWidget)
				ctx.top = childrenLayout[i].top
				ctx.left = childrenLayout[i].left
				ctx.width = childrenLayout[i].width
				ctx.height = childrenLayout[i].height

				drawContextWidget(context, ctx, childrenLayout[i].top, childrenLayout[i].left, childrenLayout[i].width, childrenLayout[i].height)
			}
		}
	}

	//debugLayout(surface, top, left, width, height)
}
