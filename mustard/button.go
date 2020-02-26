package mustard

import (
	gg "../gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateButtonWidget - Creates and returns a new Button Widget
func CreateButtonWidget(content string) *ButtonWidget {
	var widgets []interface{}

	return &ButtonWidget{
		widget: widget{
			top:  0,
			left: 0,

			width:  0,
			height: 0,

			dirty:   true,
			widgets: widgets,

			ref: "button",

			cursor: glfw.CreateStandardCursor(glfw.HandCursor),

			backgroundColor: "#fff",
		},
		content: content,

		fontSize:  20,
		fontColor: "#000",
		selected:  false,
	}

}

//AttachWidget - Attaches a new widget to the window
func (button *ButtonWidget) AttachWidget(widget interface{}) {
	button.widgets = append(button.widgets, widget)
}

//SetWidth - Sets the button width
func (button *ButtonWidget) SetWidth(width int) {
	button.width = width
	button.fixedWidth = true
}

//SetHeight - Sets the button height
func (button *ButtonWidget) SetHeight(height int) {
	button.height = height
	button.fixedHeight = true
}

//SetFontSize - Sets the button font size
func (button *ButtonWidget) SetFontSize(fontSize float64) {
	button.fontSize = fontSize
	button.dirty = true
}

//SetContent - Sets the button content
func (button *ButtonWidget) SetContent(content string) {
	button.content = content
	button.dirty = true
}

//GetContent - Gets the button content
func (button *ButtonWidget) GetContent() string {
	return button.content
}

//SetFontColor - Sets the button font color
func (button *ButtonWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		button.fontColor = fontColor
		button.dirty = true
	}
}

//SetBackgroundColor - Sets the button background color
func (button *ButtonWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		button.backgroundColor = backgroundColor
		button.dirty = true
	}
}

func drawButtonWidget(context *gg.Context, widget *ButtonWidget, top, left, width, height int) {
	if widget.selected {
		context.SetHexColor("#bbb")
	} else {
		context.SetHexColor("#ddd")
	}

	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	context.SetHexColor(widget.fontColor)
	context.LoadFontFace("roboto.ttf", widget.fontSize)

	cW, cH := context.MeasureString(widget.content)
	context.DrawString(
		widget.content,
		float64(left+width/2)-cW/2,
		float64(top+height/2)+cH/2,
	)

	context.Fill()
}
