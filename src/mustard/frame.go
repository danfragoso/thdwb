package mustard

import (
	"image"
	"image/draw"
)

//CreateFrame - Creates and returns a new Frame
func CreateFrame(orientation FrameOrientation) *Frame {
	var widgets []Widget

	return &Frame{
		baseWidget: baseWidget{
			ref: "frame",

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

	if frame.buffer == nil || frame.buffer.Bounds().Max.X != width && frame.buffer.Bounds().Max.Y != height {
		frame.buffer = image.NewRGBA(image.Rectangle{
			image.Point{}, image.Point{width, height},
		})
	}
	draw.Draw(frame.buffer, image.Rectangle{
		image.Point{},
		image.Point{width, height},
	}, context.Image(), image.Point{left, top}, draw.Over)

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
