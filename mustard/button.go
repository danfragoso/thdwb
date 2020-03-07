package mustard

import (
	gg "../gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateButtonWidget - Creates and returns a new Button Widget
func CreateButtonWidget(asset []byte) *ButtonWidget {
	var widgets []interface{}

	icon, _ := gg.LoadAsset(asset)
	return &ButtonWidget{
		widget: widget{
			needsRepaint: true,
			widgets:      widgets,

			ref: "button",

			cursor: glfw.CreateStandardCursor(glfw.HandCursor),

			backgroundColor: "#fff",
		},

		icon:      icon,
		fontSize:  20,
		padding:   0,
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
	button.box.width = width
	button.fixedWidth = true
}

//SetHeight - Sets the button height
func (button *ButtonWidget) SetHeight(height int) {
	button.box.height = height
	button.fixedHeight = true
}

//SetFontSize - Sets the button font size
func (button *ButtonWidget) SetFontSize(fontSize float64) {
	button.fontSize = fontSize
	button.needsRepaint = true
}

//SetFontSize - Sets the button font size
func (button *ButtonWidget) SetPadding(padding float64) {
	button.padding = padding
	button.needsRepaint = true
}

//SetContent - Sets the button content
func (button *ButtonWidget) SetContent(content string) {
	button.content = content
	button.needsRepaint = true
}

func (button *ButtonWidget) Click() {
	button.onClick()
}

//GetContent - Gets the button content
func (button *ButtonWidget) GetContent() string {
	return button.content
}

//SetFontColor - Sets the button font color
func (button *ButtonWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		button.fontColor = fontColor
		button.needsRepaint = true
	}
}

//SetBackgroundColor - Sets the button background color
func (button *ButtonWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		button.backgroundColor = backgroundColor
		button.needsRepaint = true
	}
}

func (button *ButtonWidget) draw(context *gg.Context) {
	top, left, width, height := button.computedBox.GetCoords()

	if button.selected {
		context.SetHexColor("#ccc")
	} else {
		context.SetHexColor("#ddd")
	}

	context.DrawRectangle(
		float64(left)+button.padding,
		float64(top)+button.padding,
		float64(width)-(button.padding*2),
		float64(height)-(button.padding*2),
	)

	context.Fill()
	context.SetHexColor(button.fontColor)
	context.LoadFontFace("roboto.ttf", button.fontSize)
	context.DrawImage(button.icon, left+4, top+2)

	context.Fill()
	button.needsRepaint = false
}
