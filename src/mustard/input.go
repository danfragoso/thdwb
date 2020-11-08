package mustard

import (
	"image"
	"image/draw"
	assets "thdwb/assets"
	gg "thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/goki/freetype/truetype"
)

//CreateInputWidget - Creates and returns a new Input Widget
func CreateInputWidget() *InputWidget {
	var widgets []Widget
	font, _ := truetype.Parse(assets.OpenSans(400))

	return &InputWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			ref: "input",

			cursor: glfw.CreateStandardCursor(glfw.IBeamCursor),

			backgroundColor: "#fff",

			font: font,
		},

		fontSize:  20,
		fontColor: "#000",
	}
}

//AttachWidget - Attaches a new widget to the window
func (input *InputWidget) AttachWidget(widget Widget) {
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

func (input *InputWidget) GetCursorPos() int {
	return input.cursorPosition
}

//SetBackgroundColor - Sets the input background color
func (input *InputWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		input.backgroundColor = backgroundColor
		input.needsRepaint = true
	}
}

func (input *InputWidget) draw() {
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
	} else {
		input.cursorPosition = 0
	}

	window.context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	window.context.Fill()

	context.SetHexColor("#2f2f2f")

	context.SetFont(input.font, input.fontSize)
	w, h := context.MeasureString(input.value)

	cursorP := float64(width - totalPadding*2)
	cP, _ := context.MeasureString(input.value[len(input.value)+input.cursorPosition:])
	cursorP = cursorP - cP

	if cursorP > 0 {
		input.cursorFloat = true
	} else {
		input.cursorFloat = false
	}

	valueBigggerThanInput := w > float64(width)-input.fontSize
	if valueBigggerThanInput && input.active {
		if cursorP > 0 {
			context.DrawStringAnchored(input.value, float64(width)-input.fontSize, float64(height+totalPadding/2)/2, 1, 0)
		} else {
			context.DrawStringAnchored(input.value, cP, float64(height+totalPadding/2)/2, 1, 0)
		}
	} else {
		context.DrawString(input.value, 0, float64(height+totalPadding/2)/2)
	}

	context.Fill()

	if input.active {
		context.SetHexColor("#000")

		if valueBigggerThanInput {
			if cursorP > 0 {
				context.DrawRectangle(cursorP, h/4, 1.3, float64(input.fontSize))
			} else {
				context.DrawRectangle(0, h/4, 1.3, float64(input.fontSize))
			}

		} else {
			cursorDefaultPosition, _ := context.MeasureString(input.value[:len(input.value)+input.cursorPosition])
			context.DrawRectangle(cursorDefaultPosition, h/4, 1.3, float64(input.fontSize))
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

	input.buffer = image.NewRGBA(image.Rectangle{
		image.Point{}, image.Point{width, height},
	})

	draw.Draw(input.buffer, image.Rectangle{
		image.Point{},
		image.Point{width, height},
	}, window.context.Image(), image.Point{left, top}, draw.Over)
}
