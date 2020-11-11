package mustard

import (
	"image"
	"image/draw"
	"log"

	gg "thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateImageWidget - Creates and returns a new Image Widget
func CreateImageWidget(path []byte) *ImageWidget {
	var widgets []Widget

	img, err := gg.LoadAsset(path)
	if err != nil {
		log.Fatal(err)
	}

	return &ImageWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			ref: "image",

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",
		},

		//path: path,
		img: img,
	}
}

//AttachWidget - Attaches a new widget to the window
func (label *ImageWidget) AttachWidget(widget Widget) {
	label.widgets = append(label.widgets, widget)
}

//SetWidth - Sets the label width
func (label *ImageWidget) SetWidth(width int) {
	label.box.width = width
	label.fixedWidth = true
	label.RequestReflow()
}

//SetHeight - Sets the label height
func (label *ImageWidget) SetHeight(height int) {
	label.box.height = height
	label.fixedHeight = true
	label.RequestReflow()
}

func (im *ImageWidget) draw() {
	top, left, width, height := im.computedBox.GetCoords()
	im.window.context.DrawImage(im.img, left+15, top+3)

	if im.buffer == nil || im.buffer.Bounds().Max.X != width && im.buffer.Bounds().Max.Y != height {
		im.buffer = image.NewRGBA(image.Rectangle{
			image.Point{}, image.Point{width, height},
		})
	}

	draw.Draw(im.buffer, image.Rectangle{
		image.Point{},
		image.Point{width, height},
	}, im.window.context.Image(), image.Point{left, top}, draw.Over)
}
