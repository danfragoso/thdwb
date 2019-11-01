package mustard

import (
	"log"

	"github.com/danfragoso/thdwb/ketchup"
	"github.com/danfragoso/thdwb/mayo"
	"github.com/danfragoso/thdwb/sauce"
	"github.com/danfragoso/thdwb/structs"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tfriedel6/canvas/backend/goglbackend"
)

func createBrowserWindow(document *structs.NodeDOM, location string) structs.AppWindow {
	defaultCursor := glfw.CreateStandardCursor(glfw.ArrowCursor)
	pageTile := getPageTitle(document) + " - thdwb"

	browserWindow := structs.AppWindow{
		Width:  900,
		Height: 600,

		ViewportWidth:  0,
		ViewportHeight: 0,

		AddressbarWidth:  0,
		AddressbarHeight: 0,

		CursorX: 0,
		CursorY: 0,

		DefaultCursor: defaultCursor,

		Title:  pageTile,
		Redraw: true,
		Resize: true,
		Reflow: true,

		ViewportOffset: 0,

		Location: location,

		DOM: document,
	}

	glfw.WindowHint(glfw.StencilBits, 8)
	glfw.WindowHint(glfw.DepthBits, 0)

	window, err := glfw.CreateWindow(browserWindow.Width, browserWindow.Height, browserWindow.Title, nil, nil)
	if err != nil {
		log.Fatalf("Error creating window: %v", err)
	}
	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		log.Fatalf("Error initializing GL: %v", err)
	}

	glfw.SwapInterval(1)
	gl.Enable(gl.MULTISAMPLE)

	addressbarBackend, err := goglbackend.New(0, 0, 0, 0, nil)
	if err != nil {
		log.Fatalf("Error loading canvas GL assets: %v", err)
	}

	viewportBackend, err := goglbackend.New(0, 0, 0, 0, nil)
	if err != nil {
		log.Fatalf("Error loading canvas GL assets: %v", err)
	}

	browserWindow.GlfwWindow = window
	browserWindow.AddressbarBackend = addressbarBackend
	browserWindow.ViewportBackend = viewportBackend

	return browserWindow
}

func handleInputChars(key rune, browserWindow *structs.AppWindow) {
	selectedElement := getSelectedUIElement(browserWindow.UIElements)

	if selectedElement != nil && selectedElement.WType == "input" {
		selectedElement.Text += string(key)
	}

	browserWindow.Redraw = true
}

func removeInputChars(browserWindow *structs.AppWindow) {
	selectedElement := getSelectedUIElement(browserWindow.UIElements)

	if selectedElement != nil && selectedElement.WType == "input" {
		inputTextLen := len(selectedElement.Text)

		if inputTextLen > 0 {
			selectedElement.Text = selectedElement.Text[:inputTextLen-1]
		}
	}

	browserWindow.Redraw = true
}

func attachBrowserWindowEvents(browserWindow *structs.AppWindow) {
	browserWindow.GlfwWindow.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyUp:
			browserWindow.ViewportOffset += 5
			break

		case glfw.KeyDown:
			browserWindow.ViewportOffset -= 5
			break

		case glfw.KeyBackspace:
			removeInputChars(browserWindow)
			break

		case glfw.KeyEnter:
			handleEnterKey(browserWindow)
			break
		}

		browserWindow.Redraw = true
	})

	browserWindow.GlfwWindow.SetCharCallback(func(w *glfw.Window, char rune) {
		handleInputChars(char, browserWindow)
	})

	browserWindow.GlfwWindow.SetScrollCallback(func(w *glfw.Window, xOffset float64, yOffset float64) {
		if yOffset < 0 {
			browserWindow.ViewportOffset -= 5
		} else {
			browserWindow.ViewportOffset += 5
		}

		browserWindow.Redraw = true
	})

	browserWindow.GlfwWindow.SetCursorPosCallback(func(w *glfw.Window, x float64, y float64) {
		if y > float64(browserWindow.AddressbarHeight) {
		} else {
			removeUIFocus(browserWindow.UIElements)
			focusedElement := getFocusedUIElement(browserWindow.UIElements, x, y)

			if focusedElement != nil {
				focusedElement.Focused = true
				w.SetCursor(focusedElement.Cursor)
			} else {
				w.SetCursor(browserWindow.DefaultCursor)
			}
		}

		browserWindow.Redraw = true
		browserWindow.CursorX = x
		browserWindow.CursorY = y
	})

	browserWindow.GlfwWindow.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			removeUISelection(browserWindow.UIElements)
			focusedElement := getFocusedUIElement(browserWindow.UIElements, browserWindow.CursorX, browserWindow.CursorY)

			if focusedElement != nil {
				focusedElement.Selected = true
				browserWindow.Redraw = true
			}
		}
	})
}

func browserWindowMainLoop(browserWindow *structs.AppWindow) {
	for !browserWindow.GlfwWindow.ShouldClose() {
		if browserWindow.Redraw {
			browserWindow.GlfwWindow.MakeContextCurrent()
			windowWidth, windowHeight := browserWindow.GlfwWindow.GetSize()

			if windowWidth != browserWindow.Width || windowHeight != browserWindow.Height {
				browserWindow.Width = windowWidth
				browserWindow.Height = windowHeight
				browserWindow.Resize = true
				browserWindow.Reflow = true
			}

			browserWindow.AddressbarWidth = browserWindow.Width
			browserWindow.AddressbarHeight = 40

			browserWindow.ViewportWidth = browserWindow.Width
			browserWindow.ViewportHeight = browserWindow.Height - browserWindow.AddressbarHeight

			browserWindow.AddressbarBackend.SetBounds(0, browserWindow.ViewportHeight, browserWindow.AddressbarWidth, browserWindow.AddressbarHeight)
			browserWindow.ViewportBackend.SetBounds(0, 0, browserWindow.ViewportWidth, browserWindow.ViewportHeight)

			updateAddressBar(browserWindow)

			if browserWindow.Reflow {
				browserWindow.DOM.Style.Width = float64(browserWindow.ViewportWidth)
				browserWindow.DOM.Style.Height = float64(browserWindow.ViewportHeight)

				browserWindow.DOM.Style.Top = 0
				browserWindow.DOM.Style.Left = 0

				nodeChildren := getNodeChildren(browserWindow.DOM)
				for i := 0; i < len(nodeChildren); i++ {
					mayo.ReflowNode(nodeChildren[i], nodeChildren[i], 0)
				}

				browserWindow.Reflow = false
			}

			updateViewport(browserWindow)

			browserWindow.Resize = false
			browserWindow.Redraw = false
			browserWindow.GlfwWindow.SwapBuffers()
		}

		glfw.WaitEvents()
	}
}

func updateAddressBar(browserWindow *structs.AppWindow) {
	oldAddressBarInput := getUIElementByID(browserWindow.UIElements, "addressbarInput")

	if browserWindow.Resize {
		browserWindow.Initialized = false
		browserWindow.UIElements = nil
	}

	if browserWindow.Initialized {
		for i := 0; i < len(browserWindow.UIElements); i++ {
			browserWindow.UIElements[i].Redraw()
		}
	} else {
		w := float64(browserWindow.AddressbarWidth)
		h := float64(browserWindow.AddressbarHeight)

		addressbarBackground := Box("addressbarBackground", 0, 0, float64(browserWindow.Width), float64(browserWindow.AddressbarHeight), browserWindow.Addressbar)

		addressbarText := browserWindow.Location
		if oldAddressBarInput != nil {
			addressbarText = oldAddressBarInput.Text
		}

		addressbarInput := Input("addressbarInput", w, h, browserWindow.Addressbar, addressbarText)

		browserWindow.UIElements = append(browserWindow.UIElements, &addressbarBackground, &addressbarInput)
		browserWindow.Initialized = true
	}
}

func renderNode(node *structs.NodeDOM, browserWindow *structs.AppWindow, vOffset float64) {
	sizeStep := node.Style.FontSize

	if node.Style.Display == "block" {
		if node.Style.Color != nil {
			browserWindow.Viewport.SetFillStyle(node.Style.Color.R, node.Style.Color.G, node.Style.Color.B)
		} else {
			browserWindow.Viewport.SetFillStyle("#000")
		}

		browserWindow.Viewport.SetFont("roboto.ttf", sizeStep)
		browserWindow.Viewport.FillText(node.Content, 0, vOffset+node.Style.Top)
	}

	children := getNodeChildren(node)

	for i := 0; i < len(children); i++ {
		if isNodeInsideViewportBounds(browserWindow, children[i], vOffset+children[i].Style.Top) {
			renderNode(children[i], browserWindow, vOffset+children[i].Style.Top)
		}
	}
}

func isNodeInsideViewportBounds(browserWindow *structs.AppWindow, node *structs.NodeDOM, vOffset float64) bool {

	if vOffset > float64(browserWindow.ViewportHeight) {
		return false
	}

	return true
}

func updateViewport(browserWindow *structs.AppWindow) {
	viewport := browserWindow.Viewport
	vO := float64(browserWindow.ViewportOffset)

	w := float64(browserWindow.ViewportWidth)
	h := float64(browserWindow.ViewportHeight)

	viewport.SetFillStyle("#FFF")
	viewport.FillRect(0, 0, w, h)
	renderNode(browserWindow.DOM, browserWindow, vO)
}

func handleEnterKey(browserWindow *structs.AppWindow) {
	addressBarInput := getUIElementByID(browserWindow.UIElements, "addressbarInput")

	if addressBarInput != nil && addressBarInput.Focused == true {
		if addressBarInput.Text != browserWindow.Location {
			url := addressBarInput.Text

			resource := sauce.GetResource(url)
			htmlString := string(resource.Body)
			parsedDocument := ketchup.ParseDocument(htmlString)

			browserWindow.DOM = parsedDocument.Children[0]
			browserWindow.ViewportOffset = 0
			browserWindow.Redraw = true
			browserWindow.Reflow = true
		}
	}
}
