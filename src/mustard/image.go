package mustard

import (
	"log"

	gg "thdwb/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//CreateImageWidget - Creates and returns a new Image Widget
func CreateImageWidget(path []byte) *ImageWidget {
	var widgets []interface{}

	img, err := gg.LoadAsset(path)
	if err != nil {
		log.Fatal(err)
	}

	return &ImageWidget{
		widget: widget{

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
func (label *ImageWidget) AttachWidget(widget interface{}) {
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

func (image *ImageWidget) draw() {
	top, left, _, _ := image.computedBox.GetCoords()
	image.window.context.DrawImage(image.img, left+15, top+3)
	image.needsRepaint = false
}
