package main

import (
	"log"

	"github.com/danfragoso/thdwb/ketchup"
	"github.com/danfragoso/thdwb/mustard"
	"github.com/danfragoso/thdwb/sauce"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	browserWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}

	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	drawingArea, _ := gtk.DrawingAreaNew()

	header, err := gtk.HeaderBarNew()
	if err != nil {
		log.Fatal("Could not create header bar:", err)
	}

	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}

	grid.Add(entry)
	grid.Add(drawingArea)

	drawingArea.SetHExpand(true)
	drawingArea.SetVExpand(true)

	entry.Connect("activate", func() {
		url, _ := entry.GetText()

		entry.SetText("")

		resource := sauce.GetResource(url)
		bodyString := string(resource.Body)
		TreeDOM := ketchup.ParseHTML(bodyString)
		html := TreeDOM.Children[0]

		header.SetTitle(mustard.GetPageTitle(html) + " - THDWB")
		drawingArea.Connect("draw", mustard.DrawDOM(html.Children[1]))

		drawingArea.QueueDraw()
		grid.QueueDraw()
		browserWindow.QueueDraw()
	})

	header.SetShowCloseButton(true)
	browserWindow.SetTitlebar(header)
	browserWindow.Connect("destroy", gtk.MainQuit)

	browserWindow.Add(grid)
	browserWindow.ShowAll()
	gtk.Main()
}
