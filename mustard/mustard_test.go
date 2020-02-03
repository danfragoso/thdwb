package mustard

import (
	"log"
	"runtime"
	"strconv"
	"testing"

	"github.com/fogleman/gg"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestMustard(t *testing.T) {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	setGLFWHints()

	app := createNewApp("THDWB")
	window := createNewWindow("THDWB", 600, 600)
	rootFrame := CreateFrame(HorizontalFrame)

	appBar := CreateFrame(VerticalFrame)

	titleBar := CreateLabelWidget("THDWB - nil")
	titleBar.SetFontColor("#fff")

	logo := CreateImageWidget("logo.png")
	logo.SetWidth(20)

	appBar.SetHeight(28)
	appBar.AttachWidget(logo)
	appBar.AttachWidget(titleBar)
	appBar.SetBackgroundColor("#5f6368")

	rootFrame.AttachWidget(appBar)

	viewPort := CreateContextWidget(func(ctx *gg.Context) {
		img, err := gg.LoadImage("logo.png")
		if err != nil {
			log.Fatal(err)
		}
		ctx.DrawImage(img, 0, 0)
	})

	rootFrame.AttachWidget(viewPort)

	statusBar := CreateFrame(HorizontalFrame)
	statusBar.SetBackgroundColor("#babcbe")
	statusBar.SetHeight(20)

	statusLabel := CreateLabelWidget("Processed Events:")
	statusLabel.SetFontSize(16)
	frameEvents := 0

	rootFrame.AttachWidget(statusBar)
	statusBar.AttachWidget(statusLabel)

	window.SetRootFrame(rootFrame)

	app.AddWindow(window)

	window.Show()
	app.Run(func() {
		frameEvents++
		statusLabel.SetContent("Processed Events: " + strconv.Itoa(frameEvents) + "; Resolution: " + strconv.Itoa(window.width) + "X" + strconv.Itoa(window.height))
		//window.RequestRepaint()
	})
}
