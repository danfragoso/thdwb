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
			top:  0,
			left: 0,

			width:  0,
			height: 0,

			dirty:   true,
			widgets: widgets,

			ref: "input",

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

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
	input.width = width
	input.fixedWidth = true
}

//SetHeight - Sets the input height
func (input *InputWidget) SetHeight(height int) {
	input.height = height
	input.fixedHeight = true
}

//SetFontSize - Sets the input font size
func (input *InputWidget) SetFontSize(fontSize float64) {
	input.fontSize = fontSize
	input.dirty = true
}

//SetFontColor - Sets the input font color
func (input *InputWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		input.fontColor = fontColor
		input.dirty = true
	}
}

//SetBackgroundColor - Sets the input background color
func (input *InputWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		input.backgroundColor = backgroundColor
		input.dirty = true
	}
}

func drawInputWidget(context *gg.Context, widget *InputWidget, top, left, width, height int) {
	context.SetHexColor(widget.fontColor)
	context.LoadFontFace("roboto.ttf", widget.fontSize)
	context.DrawString("memes", float64(left)+widget.fontSize/4, float64(top)+widget.fontSize*2/2)
	context.Fill()
	//debugLayout(surface, top, left, width, height)
}
