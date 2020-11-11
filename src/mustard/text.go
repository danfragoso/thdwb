package mustard

import (
	"image"
	"image/draw"
	assets "thdwb/assets"
	gg "thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/goki/freetype/truetype"
)

//CreateTextWidget - Creates and returns a new Text Widget
func CreateTextWidget(content string) *TextWidget {
	var widgets []Widget
	font, _ := truetype.Parse(assets.OpenSans(400))

	return &TextWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			ref: "text",

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",

			font: font,
		},
		content: content,

		fontSize:  20,
		fontColor: "#000",
	}
}

//AttachWidget - Attaches a new widget to the window
func (text *TextWidget) AttachWidget(widget Widget) {
	text.widgets = append(text.widgets, widget)
}

//SetWidth - Sets the text width
func (text *TextWidget) SetWidth(width int) {
	text.box.width = width
	text.fixedWidth = true
	text.RequestReflow()
}

//SetHeight - Sets the text height
func (text *TextWidget) SetHeight(height int) {
	text.box.height = height
	text.fixedHeight = true
	text.RequestReflow()
}

//SetFontSize - Sets the text font size
func (text *TextWidget) SetFontSize(fontSize float64) {
	text.fontSize = fontSize
	text.needsRepaint = true
}

//SetContent - Sets the text content
func (text *TextWidget) SetContent(content string) {
	text.content = content
	text.needsRepaint = true
	text.RequestReflow()
}

//GetContent - Gets the text content
func (text *TextWidget) GetContent() string {
	return text.content
}

//SetFontColor - Sets the text font color
func (text *TextWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		text.fontColor = fontColor
		text.needsRepaint = true
	}
}

//SetBackgroundColor - Sets the text background color
func (text *TextWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		text.backgroundColor = backgroundColor
		text.needsRepaint = true
	}
}

func (text *TextWidget) draw() {
	context := text.window.context
	top, left, width, height := text.computedBox.GetCoords()

	context.SetFont(text.font, text.fontSize)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.SetHexColor(text.backgroundColor)
	context.Fill()

	context.SetHexColor(text.fontColor)
	context.DrawStringWrapped(text.content, float64(left)+text.fontSize/4, float64(top)+text.fontSize*2/2, 0, 0, float64(width), text.fontSize*0.15, gg.AlignLeft)
	context.Fill()

	if text.buffer == nil || text.buffer.Bounds().Max.X != width && text.buffer.Bounds().Max.Y != height {
		text.buffer = image.NewRGBA(image.Rectangle{
			image.Point{}, image.Point{width, height},
		})
	}

	draw.Draw(text.buffer, image.Rectangle{
		image.Point{},
		image.Point{width, height},
	}, context.Image(), image.Point{left, top}, draw.Over)
}
