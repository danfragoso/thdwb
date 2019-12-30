package main

import (
	"github.com/danfragoso/mustard"
	"github.com/tfriedel6/canvas"
)

func main() {
	app := mustard.CreateApp("thdwb")
	mainWindow := app.CreateWindow("thdwb - Inicio", 100, 100, 50, 50)
	rootFrame := mustard.CreateFrame(mustard.HorizontalFrame)
	size := float64(10)

	viewport := mustard.CreateSurfaceWidget(mainWindow, func(surface *canvas.Canvas, top, left, width, height int) {
		viewportRendererCallback(surface, size, top, left, width, height)
	})
	addressBar := mustard.CreateFrame(mustard.VerticalFrame)

	addressBar.SetHeight(40)

	fetchBt := mustard.CreateButtonWidget("+", func() {
		size += 10
	})

	renderBt := mustard.CreateButtonWidget("-", func() {
		size -= 10
	})

	addressBar.AttachWidget(fetchBt)
	addressBar.AttachWidget(renderBt)

	rootFrame.AttachWidget(addressBar)
	rootFrame.AttachWidget(viewport)

	mainWindow.SetRootFrame(rootFrame)
	mainWindow.Show()

	app.Run(func() {})
}

func viewportRendererCallback(surface *canvas.Canvas, size float64, top, left, width, height int) {
	surface.SetFillStyle("#00ff00")
	surface.FillRect(float64(left+width/2)-size/2, float64(top+height/2)-size/2, size, size)
}
