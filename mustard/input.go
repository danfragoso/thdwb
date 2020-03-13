package mustard

import (
	gg "../gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateInputWidget - Creates and returns a new Input Widget
func CreateInputWidget() *InputWidget {
	var widgets []interface{}

	return &InputWidget{
		widget: widget{

			needsRepaint: true,
			widgets:      widgets,

			ref: "input",

			cursor: glfw.CreateStandardCursor(glfw.IBeamCursor),

			backgroundColor: "#fff",
		},

		fontSize:  20,
		fontColor: "#000",
	}
}

//AttachWidget - Attaches a new widget to the window
func (input *InputWidget) AttachWidget(widget interface{}) {
	input.widgets = append(input.widgets, widget)
}

//SetWidth - Sets the input width
func (input *InputWidget) SetWidth(width int) {
	input.box.width = width
	input.fixedWidth = true
	input.RequestReflow()
}

//SetHeight - Sets the input height
func (input *InputWidget) SetHeight(height int) {
	input.box.height = height
	input.fixedHeight = true
	input.RequestReflow()
}

//SetFontSize - Sets the input font size
func (input *InputWidget) SetFontSize(fontSize float64) {
	input.fontSize = fontSize
	input.needsRepaint = true
}

func (input *InputWidget) SetReturnCallback(returnCallback func()) {
	input.returnCallback = returnCallback
}

//SetFontColor - Sets the input font color
func (input *InputWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		input.fontColor = fontColor
		input.needsRepaint = true
	}
}

//SetFontColor - Sets the input font color
func (input *InputWidget) SetValue(value string) {
	input.value = value
	input.needsRepaint = true
}

//SetFontColor - Sets the input font color
func (input *InputWidget) GetValue() string {
	return input.value
}

//SetBackgroundColor - Sets the input background color
func (input *InputWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		input.backgroundColor = backgroundColor
		input.needsRepaint = true
	}
}

func (input *InputWidget) draw() {
	// For some reason the gg Clip function is reaaaaaally slow to clip text?
	// A faster way to do this is to create a context the size of the input
	// and draw the text on this context then draw the input context.Image()
	// on the main context.

	top, left, width, height := input.computedBox.GetCoords()
	if input.context == nil || input.context.Width() != width || input.context.Height() != height {
		input.context = gg.NewContext(width, height)

		input.context.SetRGB(1, 1, 1)
		input.context.Clear()
	}

	window := input.window
	context := input.context

	if input.selected {
		context.SetHexColor("#e4e4e4")
	} else {
		context.SetHexColor("#efefef")
	}

	if input.active {
		context.SetHexColor("#fff")
	}

	context.DrawRectangle(0, 0, float64(width), float64(height))
	context.Fill()

	context.SetHexColor("#2f2f2f")
	context.LoadFontFace("roboto.ttf", input.fontSize)
	w, _ := context.MeasureString(input.value)

	if w > float64(width)-input.fontSize && input.active {
		context.DrawStringAnchored(input.value, float64(width)-input.fontSize+4, float64(height)/2+2+input.fontSize/4, 1, 0)
		context.SetRGB(1, 1, 1)
		context.DrawRectangle(-1, float64(height)/2-input.fontSize*1.2/2, 6, input.fontSize*1.2)
		context.Fill()
	} else {
		context.DrawString(input.value, input.fontSize/4+4, float64(height)/2+2+input.fontSize/4)
	}

	context.Fill()

	if input.active {
		context.SetHexColor("#000")
		context.DrawRectangle(
			input.fontSize/4+4+w,
			float64(height)/2-input.fontSize/2+.5,
			1.3,
			float64(input.fontSize),
		)
		context.Fill()
	}

	window.context.DrawImage(context.Image(), left, top)
	window.context.SetHexColor("#000")
	window.context.SetLineWidth(.4)

	window.context.DrawRectangle(
		float64(left)+1+input.padding,
		float64(top)+1+input.padding,
		float64(width)-2-(input.padding*2),
		float64(height)-2-(input.padding*2),
	)

	window.context.SetLineJoinRound()
	window.context.Stroke()
	input.needsRepaint = false
}
