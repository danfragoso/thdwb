package mustard

import (
	"runtime"
	"strconv"
	"testing"

	gg "github.com/danfragoso/thdwb/gg"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestMustard(t *testing.T) {
	runtime.LockOSThread()
	glfw.Init()
	gl.Init()

	SetGLFWHints()

	app := CreateNewApp("THDWB")
	window := CreateNewWindow("THDWB", 600, 600)
	rootFrame := CreateFrame(HorizontalFrame)

	appBar := CreateFrame(VerticalFrame)

	titleBar := CreateLabelWidget("THDWB - nil")
	titleBar.SetFontColor("#fff")

	appBar.SetHeight(28)
	appBar.AttachWidget(titleBar)
	appBar.SetBackgroundColor("#5f6368")

	rootFrame.AttachWidget(appBar)

	viewPort := CreateCanvasWidget(func(ctx *gg.Context) {})

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
