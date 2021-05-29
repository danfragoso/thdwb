package mustard

import (
	"log"

	gg "github.com/danfragoso/thdwb/gg"

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

			widgetType: imageWidget,

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",
		},

		//path: path,
		img: img,
	}
}

//SetWidth - Sets the label width
func (label *ImageWidget) SetWidth(width float64) {
	label.box.width = width
	label.fixedWidth = true
	label.RequestReflow()
}

//SetHeight - Sets the label height
func (label *ImageWidget) SetHeight(height float64) {
	label.box.height = height
	label.fixedHeight = true
	label.RequestReflow()
}

func (im *ImageWidget) draw() {
	top, left, _, _ := im.computedBox.GetCoords()
	im.window.context.DrawImage(im.img, int(left)+15, int(top)+3)

	copyWidgetToBuffer(im, im.window.context.Image())
}
