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
	input.padding = 4
	top, left, width, height := input.computedBox.GetCoords()
	totalPadding := int(input.padding * 2)
	if input.context == nil || input.context.Width() != width-totalPadding || input.context.Height() != height-totalPadding {
		input.context = gg.NewContext(width-totalPadding, height-totalPadding)
	}

	window := input.window
	context := input.context

	if input.selected {
		window.context.SetHexColor("#e4e4e4")
		context.SetHexColor("#e4e4e4")
		context.Clear()
	} else {
		window.context.SetHexColor("#efefef")
		context.SetHexColor("#efefef")
		context.Clear()
	}

	if input.active {
		window.context.SetHexColor("#fff")
		context.SetHexColor("#fff")
		context.Clear()
	}

	window.context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	window.context.Fill()

	context.SetHexColor("#2f2f2f")
	context.LoadFontFace("roboto.ttf", input.fontSize)
	w, h := context.MeasureString(input.value)

	valueBigggerThanInput := w > float64(width)-input.fontSize
	if valueBigggerThanInput && input.active {
		context.DrawStringAnchored(input.value, float64(width)-input.fontSize, float64(height+totalPadding/2)/2, 1, 0)
	} else {
		context.DrawString(input.value, 0, float64(height+totalPadding/2)/2)
	}

	context.Fill()

	//CURSOR
	if input.active {
		context.SetHexColor("#000")

		if valueBigggerThanInput {
			context.DrawRectangle(float64(width-totalPadding*2), h/4, 1.3, float64(input.fontSize))
		} else {
			context.DrawRectangle(w, h/4, 1.3, float64(input.fontSize))
		}

		context.Fill()
	}

	window.context.DrawImage(context.Image(), left+totalPadding/2, top+totalPadding/2)
	window.context.SetHexColor("#000")
	window.context.SetLineWidth(.4)

	window.context.DrawRectangle(
		float64(left)+1,
		float64(top)+1,
		float64(width)-2,
		float64(height)-2,
	)

	window.context.SetLineJoinRound()
	window.context.Stroke()
	input.needsRepaint = false
}
