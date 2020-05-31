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
	frame.RequestReflow()
}

//SetHeight - Sets the frame height
func (frame *Frame) SetHeight(height int) {
	frame.box.height = height
	frame.fixedHeight = true
	frame.RequestReflow()
}

//SetHeight - Sets the frame height
func (frame *Frame) GetHeight() int {
	return frame.box.height
}

func drawRootFrame(window *Window) {
	window.rootFrame.computedBox.SetCoords(0, 0, window.width, window.height)

	window.rootFrame.draw()
}

func (frame *Frame) draw() {
	top, left, width, height := frame.computedBox.GetCoords()
	window := frame.window
	context := window.context
	context.SetHexColor(frame.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	frame.needsRepaint = false
	childrenLen := len(frame.widgets)
	if childrenLen > 0 {
		childrenWidgets := getCoreWidgets(frame.widgets)
		childrenLayout := calculateChildrenWidgetsLayout(childrenWidgets, top, left, width, height, frame.orientation)

		for i := 0; i < childrenLen; i++ {
			switch frame.widgets[i].(type) {
			case *Frame:
				frame := frame.widgets[i].(*Frame)
				frame.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				frame.draw()

			case *LabelWidget:
				label := frame.widgets[i].(*LabelWidget)
				label.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				label.draw()

			case *TextWidget:
				text := frame.widgets[i].(*TextWidget)
				text.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				text.draw()

			case *ImageWidget:
				image := frame.widgets[i].(*ImageWidget)
				image.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				image.draw()

			case *CanvasWidget:
				ctx := frame.widgets[i].(*CanvasWidget)
				ctx.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				ctx.draw()

			case *ButtonWidget:
				button := frame.widgets[i].(*ButtonWidget)
				button.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				button.draw()

			case *InputWidget:
				input := frame.widgets[i].(*InputWidget)
				input.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				input.draw()

			case *ScrollBarWidget:
				scrollBar := frame.widgets[i].(*ScrollBarWidget)
				scrollBar.computedBox.SetCoords(childrenLayout[i].box.GetCoords())
				scrollBar.draw()
			}

		}
	}
}
