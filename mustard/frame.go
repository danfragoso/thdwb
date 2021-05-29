package mustard

//CreateFrame - Creates and returns a new Frame
func CreateFrame(orientation FrameOrientation) *Frame {
	var widgets []Widget

	return &Frame{
		baseWidget: baseWidget{
			widgetType: frameWidget,

			needsRepaint: true,
			widgets:      widgets,

			backgroundColor: "#fff"},

		orientation: orientation,
	}
}

//SetBackgroundColor - Sets the frame background color
func (frame *Frame) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		frame.backgroundColor = backgroundColor
		frame.needsRepaint = true
	}
}

//SetWidth - Sets the frame width
func (frame *Frame) SetWidth(width float64) {
	frame.box.width = width
	frame.fixedWidth = true
	frame.RequestReflow()
}

//SetHeight - Sets the frame height
func (frame *Frame) SetHeight(height float64) {
	frame.box.height = height
	frame.fixedHeight = true
	frame.RequestReflow()
}

//SetHeight - Sets the frame height
func (frame *Frame) GetHeight() float64 {
	return frame.box.height
}

func drawRootFrame(window *Window) {
	window.rootFrame.computedBox.SetCoords(0, 0, float64(window.width), float64(window.height))

	window.rootFrame.draw()
}

func (frame *Frame) draw() {
	top, left, width, height := frame.computedBox.GetCoords()
	window := frame.window
	context := window.context
	context.SetHexColor(frame.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	copyWidgetToBuffer(frame, context.Image())

	childrenLen := len(frame.widgets)
	if childrenLen > 0 {
		childrenWidgets := getCoreWidgets(frame.widgets)
		childrenLayout := calculateChildrenWidgetsLayout(childrenWidgets, top, left, width, height, frame.orientation)

		for idx, child := range frame.Widgets() {
			child.ComputedBox().SetCoords(childrenLayout[idx].box.GetCoords())
			child.draw()
		}
	}
}
