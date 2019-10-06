package mustard

import (
	"runtime"

	"github.com/danfragoso/thdwb/structs"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tfriedel6/canvas"
)

//RenderDocument "Renders the DOM to the screen"
func RenderDocument(document *structs.NodeDOM) {
	html := document.Children[0]
	runtime.LockOSThread()
	glfw.Init()

	browserWindow := createBrowserWindow(html)
	attachBrowserWindowEvents(&browserWindow)
	browserWindow.Viewport = canvas.New(browserWindow.ViewportBackend)
	browserWindow.Addressbar = canvas.New(browserWindow.AddressbarBackend)
	browserWindowMainLoop(&browserWindow)
}
