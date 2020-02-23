package mustard

import (
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateTextWidget - Creates and returns a new Text Widget
func CreateTextWidget(content string) *TextWidget {
	var widgets []interface{}

	return &TextWidget{
		widget: widget{
			top:  0,
			left: 0,

			width:  0,
			height: 0,

			dirty:   true,
			widgets: widgets,

			ref: "text",

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",
		},
		content: content,

		fontSize:  20,
		fontColor: "#000",
	}
}

//AttachWidget - Attaches a new widget to the window
func (text *TextWidget) AttachWidget(widget interface{}) {
	text.widgets = append(text.widgets, widget)
}

//SetWidth - Sets the text width
func (text *TextWidget) SetWidth(width int) {
	text.width = width
	text.fixedWidth = true
}

//SetHeight - Sets the text height
func (text *TextWidget) SetHeight(height int) {
	text.height = height
	text.fixedHeight = true
}

//SetFontSize - Sets the text font size
func (text *TextWidget) SetFontSize(fontSize float64) {
	text.fontSize = fontSize
	text.dirty = true
}

//SetContent - Sets the text content
func (text *TextWidget) SetContent(content string) {
	text.content = content
	text.dirty = true
}

//GetContent - Gets the text content
func (text *TextWidget) GetContent() string {
	return text.content
}

//SetFontColor - Sets the text font color
func (text *TextWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		text.fontColor = fontColor
		text.dirty = true
	}
}

//SetBackgroundColor - Sets the text background color
func (text *TextWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		text.backgroundColor = backgroundColor
		text.dirty = true
	}
}

func drawTextWidget(context *gg.Context, widget *TextWidget, top, left, width, height int) {
	context.LoadFontFace("roboto.ttf", widget.fontSize)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.SetHexColor(widget.backgroundColor)
	context.Fill()

	context.SetHexColor(widget.fontColor)
	context.DrawStringWrapped(widget.content, float64(left)+widget.fontSize/4, float64(top)+widget.fontSize*2/2, 0, 0, float64(width), widget.fontSize*0.15, gg.AlignLeft)
	context.Fill()
	//debugLayout(surface, top, left, width, height)
}
