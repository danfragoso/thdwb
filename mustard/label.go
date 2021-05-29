package mustard

import (
	"image"
	"image/draw"

	assets "github.com/danfragoso/thdwb/assets"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/goki/freetype/truetype"
)

//CreateLabelWidget - Creates and returns a new Label Widget
func CreateLabelWidget(content string) *LabelWidget {
	var widgets []Widget
	font, _ := truetype.Parse(assets.OpenSans(400))

	return &LabelWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			widgetType: labelWidget,

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",

			font: font,
		},
		content: content,

		fontSize:  20,
		fontColor: "#000",
	}
}

//SetWidth - Sets the label width
func (label *LabelWidget) SetWidth(width float64) {
	label.box.width = width
	label.fixedWidth = true
	label.RequestReflow()
}

//SetHeight - Sets the label height
func (label *LabelWidget) SetHeight(height float64) {
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
	context.SetFont(label.font, label.fontSize)
	context.DrawString(label.content, float64(left)+label.fontSize/4, float64(top)+label.fontSize*2/2)
	context.Fill()

	if label.buffer == nil || label.buffer.Bounds().Max.X != int(width) && label.buffer.Bounds().Max.Y != int(height) {
		label.buffer = image.NewRGBA(image.Rectangle{
			image.Point{}, image.Point{int(width), int(height)},
		})
	}

	draw.Draw(label.buffer, image.Rectangle{
		image.Point{},
		image.Point{int(width), int(height)},
	}, context.Image(), image.Point{int(left), int(top)}, draw.Over)
}
