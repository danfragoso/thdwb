package main

import (
	"fmt"
	"image"
	"strings"

	gg "thdwb/gg"
	ketchup "thdwb/ketchup"
	mustard "thdwb/mustard"
	profiler "thdwb/profiler"
	sauce "thdwb/sauce"
	structs "thdwb/structs"
)

func loadDocument(browser *structs.WebBrowser, link string) {
	URL := sauce.ParseURL(link)

	if URL.Scheme == "" && URL.Host == "" {
		if !strings.HasPrefix(URL.Path, "/") {
			URL.Path = "/" + URL.Path
		}

		URL = sauce.ParseURL(browser.ActiveDocument.URL.Scheme + "://" + browser.ActiveDocument.URL.Host + URL.Path)
	}

	resource := sauce.GetResource(URL, browser)
	rawDocument := string(resource.Body)

	switch strings.Split(resource.ContentType, ";")[0] {
	case "text/plain", "application/json":
		browser.ActiveDocument = ketchup.ParsePlainText(rawDocument)
	default:
		browser.ActiveDocument = ketchup.ParseHTML(rawDocument)
	}

	browser.ActiveDocument.URL = resource.URL
	browser.ActiveDocument.ContentType = resource.ContentType

	browser.Window.RemoveStaticOverlay("debugOverlay")

	if browser.History.PageCount() == 0 || browser.History.Last().String() != resource.URL.String() {
		browser.History.Push(resource.URL)
	}
}

func loadDocumentFromUrl(browser *structs.WebBrowser, statusLabel *mustard.LabelWidget, urlInput *mustard.InputWidget, viewPort *mustard.CanvasWidget) {
	statusLabel.SetContent("Loading: " + urlInput.GetValue())
	statusLabel.RequestRepaint()

	loadDocument(browser, urlInput.GetValue())
	viewPort.SetOffset(0)
	viewPort.SetDrawingRepaint(true)
	viewPort.RequestRepaint()
	urlInput.SetValue(browser.ActiveDocument.URL.String())
}

func createStatusLabel(perf *profiler.Profiler) string {
	return "Loaded; " +
		"Render took: " + perf.GetProfile("render").GetElapsedTime().String() + "; "
}

func processPointerPositionEvent(browser *structs.WebBrowser, x, y float64) {
	y -= float64(browser.Viewport.GetOffset())
	selectedElement := browser.ActiveDocument.DOM.CalcPointIntersection(x, y)

	if browser.ActiveDocument.SelectedElement != selectedElement {
		browser.ActiveDocument.SelectedElement = selectedElement

		if browser.ActiveDocument.SelectedElement != nil && browser.ActiveDocument.SelectedElement.Element == "a" {
			browser.Window.SetCursor("pointer")
			browser.StatusLabel.SetContent(browser.ActiveDocument.SelectedElement.Attr("href"))
		} else {
			browser.Window.SetCursor("default")
			browser.StatusLabel.SetContent(createStatusLabel(browser.Profiler))
		}

		if browser.ActiveDocument.DebugFlag &&
			browser.ActiveDocument.SelectedElement != nil &&
			browser.ActiveDocument.SelectedElement.Element != "html" {

			printNodeDebug(browser.ActiveDocument.SelectedElement)
			showDebugOverlay(browser)
		}

		browser.StatusLabel.RequestRepaint()
	}
}

func printNodeDebug(node *structs.NodeDOM) {
	rect := fmt.Sprintf("{%.2f, %.2f, %.2f, %.2f}", node.RenderBox.Top, node.RenderBox.Left, node.RenderBox.Width, node.RenderBox.Height)
	fmt.Printf("%s [\n %s\n]\n\n", node.Element, rect)
}

func showDebugOverlay(browser *structs.WebBrowser) {
	browser.Window.RemoveStaticOverlay("debugOverlay")

	debugEl := browser.ActiveDocument.SelectedElement
	top, left, _, height := debugEl.RenderBox.GetRect()
	ctx := gg.NewContext(int(browser.ActiveDocument.DOM.RenderBox.Width), int(height+20))
	paintDebugRect(ctx, debugEl)

	overlay := mustard.CreateStaticOverlay("debugOverlay", ctx, image.Point{
		int(left), int(top) + browser.Viewport.GetTop() + browser.Viewport.GetOffset(),
	})

	browser.Window.AddStaticOverlay(overlay)
}

func paintDebugRect(ctx *gg.Context, node *structs.NodeDOM) {
	debugString := node.Element + " {" + fmt.Sprint(node.RenderBox.Top, node.RenderBox.Left, node.RenderBox.Width, node.RenderBox.Height) + "}"
	ctx.DrawRectangle(0, 0, node.RenderBox.Width, node.RenderBox.Height)
	ctx.SetRGBA(.2, .8, .4, .3)
	ctx.Fill()

	w, h := ctx.MeasureString(debugString)

	if node.RenderBox.Width < w {
		ctx.DrawRectangle(0, node.RenderBox.Height, w+4, h+4)
		ctx.SetRGB(1, 1, 0)
		ctx.Fill()

		ctx.SetRGB(0, 0, 0)
		ctx.DrawString(debugString, 2, node.RenderBox.Height+h)
		ctx.Fill()
	} else {
		ctx.DrawRectangle(node.RenderBox.Width-w-2, node.RenderBox.Height, w+4, h+4)
		ctx.SetRGB(1, 1, 0)
		ctx.Fill()

		ctx.SetRGB(0, 0, 0)
		ctx.DrawString(debugString, node.RenderBox.Width-w, node.RenderBox.Height+h)
		ctx.Fill()
	}
}
