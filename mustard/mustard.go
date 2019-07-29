package mustard

import "github.com/gotk3/gotk3/cairo"
import "github.com/gotk3/gotk3/gdk"
import "github.com/gotk3/gotk3/gtk"
import "github.com/danfragoso/thdwb/ketchup"
import "log"

const KEY_LEFT  uint = 65361
const KEY_UP    uint = 65362
const KEY_RIGHT uint = 65363
const KEY_DOWN  uint = 65364

func RenderDOM(DOM_Tree *ketchup.DOM_Node) {
	gtk.Init(nil)
	
	browserWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	drawingArea, _ := gtk.DrawingAreaNew()
	
	browserWindow.Add(drawingArea)
	
	header, err := gtk.HeaderBarNew()
	if err != nil {
		log.Fatal("Could not create header bar:", err)
	}
	
	header.SetShowCloseButton(true)
	header.SetTitle(DOM_Tree.Content)
	header.SetSubtitle("thdwb")
	
	browserWindow.SetTitlebar(header)
	browserWindow.Connect("destroy", gtk.MainQuit)
	browserWindow.ShowAll()
	
	unitSize := 20.0
	x := 0.0
	y := 0.0
	keyMap := map[uint]func(){
		KEY_LEFT:  func() { x-- },
		KEY_UP:    func() { y-- },
		KEY_RIGHT: func() { x++ },
		KEY_DOWN:  func() { y++ },
	}
	
	drawingArea.Connect("draw", func(drawingArea *gtk.DrawingArea, cr *cairo.Context) {
		cr.SetSourceRGB(0, 0, 0)
		cr.Rectangle(x*unitSize, y*unitSize, unitSize, unitSize)
		cr.Fill()
	})
	
	browserWindow.Connect("key-press-event", func(browserWindow *gtk.Window, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		if move, found := keyMap[keyEvent.KeyVal()]; found {
			move()
			browserWindow.QueueDraw()
		}
	})
	
	gtk.Main()
}