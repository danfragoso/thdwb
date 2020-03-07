package mustard

//CreateFrame - Creates and returns a new Frame
func CreateFrame(orientation FrameOrientation) *Frame {
	var widgets []interface{}

	return &Frame{
		widget: widget{
			ref: "frame",

			needsRepaint: true,
			widgets:      widgets,

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
		frame.needsRepaint = true
	}
}

//SetWidth - Sets the frame width
func (frame *Frame) SetWidth(width int) {
	frame.box.width = width
	frame.fixedWidth = true
}

//SetHeight - Sets the frame height
func (frame *Frame) SetHeight(height int) {
	frame.box.height = height
	frame.fixedHeight = true
}

//SetHeight - Sets the frame height
func (frame *Frame) GetHeight() int {
	return frame.box.height
}

func drawRootFrame(window *Window) {
	drawFrame(window, window.rootFrame, 0, 0, window.width, window.height)
}

func drawFrame(window *Window, frame *Frame, top, left, width, height int) {
	context := window.context
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
				frame.box.top = childrenLayout[i].box.top
				frame.box.left = childrenLayout[i].box.left
				frame.box.width = childrenLayout[i].box.width
				frame.box.height = childrenLayout[i].box.height

				drawFrame(window, frame, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)
			case *LabelWidget:
				label := frame.widgets[i].(*LabelWidget)
				label.box.top = childrenLayout[i].box.top
				label.box.left = childrenLayout[i].box.left
				label.box.width = childrenLayout[i].box.width
				label.box.height = childrenLayout[i].box.height

				drawLabelWidget(context, label, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)
			case *TextWidget:
				text := frame.widgets[i].(*TextWidget)
				text.box.top = childrenLayout[i].box.top
				text.box.left = childrenLayout[i].box.left
				text.box.width = childrenLayout[i].box.width
				text.box.height = childrenLayout[i].box.height

				drawTextWidget(context, text, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)

			case *ImageWidget:
				image := frame.widgets[i].(*ImageWidget)
				image.box.top = childrenLayout[i].box.top
				image.box.left = childrenLayout[i].box.left
				image.box.width = childrenLayout[i].box.width
				image.box.height = childrenLayout[i].box.height

				drawImageWidget(context, image, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)
			case *ContextWidget:
				ctx := frame.widgets[i].(*ContextWidget)
				ctx.box.top = childrenLayout[i].box.top
				ctx.box.left = childrenLayout[i].box.left
				ctx.box.width = childrenLayout[i].box.width
				ctx.box.height = childrenLayout[i].box.height

				drawContextWidget(context, ctx, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)

			case *ButtonWidget:
				button := frame.widgets[i].(*ButtonWidget)
				button.box.top = childrenLayout[i].box.top
				button.box.left = childrenLayout[i].box.left
				button.box.width = childrenLayout[i].box.width
				button.box.height = childrenLayout[i].box.height

				drawButtonWidget(context, button, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)

			case *InputWidget:
				input := frame.widgets[i].(*InputWidget)
				input.box.top = childrenLayout[i].box.top
				input.box.left = childrenLayout[i].box.left
				input.box.width = childrenLayout[i].box.width
				input.box.height = childrenLayout[i].box.height

				drawInputWidget(window, input, childrenLayout[i].box.top, childrenLayout[i].box.left, childrenLayout[i].box.width, childrenLayout[i].box.height)
			}
		}
	}
}
