package mustard

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateLabelWidget - Creates and returns a new Label Widget
func CreateLabelWidget(content string) *LabelWidget {
	var widgets []interface{}

	return &LabelWidget{
		widget: widget{

			needsRepaint: true,
			widgets:      widgets,

			ref: "label",

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",
		},
		content: content,

		fontSize:  20,
		fontColor: "#000",
	}
}

//AttachWidget - Attaches a new widget to the window
func (label *LabelWidget) AttachWidget(widget interface{}) {
	label.widgets = append(label.widgets, widget)
}

//SetWidth - Sets the label width
func (label *LabelWidget) SetWidth(width int) {
	label.box.width = width
	label.fixedWidth = true
	label.RequestReflow()
}

//SetHeight - Sets the label height
func (label *LabelWidget) SetHeight(height int) {
	label.box.height = height
	label.fixedHeight = true
	label.RequestReflow()
}

//SetFontSize - Sets the label font size
func (label *LabelWidget) SetFontSize(fontSize float64) {
	label.fontSize = fontSize
	label.needsRepaint = true
}

//SetContent - Sets the label content
func (label *LabelWidget) SetContent(content string) {
	label.content = content
	label.needsRepaint = true
}

//GetContent - Gets the label content
func (label *LabelWidget) GetContent() string {
	return label.content
}

//SetFontColor - Sets the label font color
func (label *LabelWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		label.fontColor = fontColor
		label.needsRepaint = true
	}
}

//SetBackgroundColor - Sets the label background color
func (label *LabelWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		label.backgroundColor = backgroundColor
		label.needsRepaint = true
	}
}

func (label *LabelWidget) draw() {
	context := label.window.context
	top, left, width, height := label.computedBox.GetCoords()

	context.SetHexColor(label.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	context.SetHexColor(label.fontColor)
	context.LoadFontFace("roboto.ttf", label.fontSize)
	context.DrawString(label.content, float64(left)+label.fontSize/4, float64(top)+label.fontSize*2/2)
	context.Fill()

	label.needsRepaint = false
}
