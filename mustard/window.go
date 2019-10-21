package mustard

import (
	"log"

	"github.com/danfragoso/thdwb/structs"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tfriedel6/canvas/backend/goglbackend"
)

func createBrowserWindow(document *structs.NodeDOM) structs.AppWindow {
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

		ViewportOffset: 0,
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

		case glfw.KeyDown:
			browserWindow.ViewportOffset -= 5

		case glfw.KeyBackspace:
			removeInputChars(browserWindow)
		default:
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
			}

			browserWindow.AddressbarWidth = browserWindow.Width
			browserWindow.AddressbarHeight = 40

			browserWindow.ViewportWidth = browserWindow.Width
			browserWindow.ViewportHeight = browserWindow.Height - browserWindow.AddressbarHeight

			browserWindow.AddressbarBackend.SetBounds(0, browserWindow.ViewportHeight, browserWindow.AddressbarWidth, browserWindow.AddressbarHeight)
			browserWindow.ViewportBackend.SetBounds(0, 0, browserWindow.ViewportWidth, browserWindow.ViewportHeight)

			updateAddressBar(browserWindow)
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

		addressbarText := ""
		if oldAddressBarInput != nil {
			addressbarText = oldAddressBarInput.Text
		}

		addressbarInput := Input("addressbarInput", w, h, browserWindow.Addressbar, addressbarText)

		browserWindow.UIElements = append(browserWindow.UIElements, &addressbarBackground, &addressbarInput)
		browserWindow.Initialized = true
	}
}

func updateViewport(browserWindow *structs.AppWindow) {
	viewport := browserWindow.Viewport
	vO := float64(browserWindow.ViewportOffset)

	w := float64(browserWindow.ViewportWidth)
	h := float64(browserWindow.ViewportHeight)

	x := browserWindow.CursorX
	y := browserWindow.CursorY - float64(browserWindow.AddressbarHeight)

	viewport.SetFillStyle("#A0B")
	viewport.FillRect(0, 0, w, h)

	viewport.SetFillStyle("#111")
	viewport.FillRect(0, vO, 100, 100)

	viewport.SetFillStyle("#222")
	viewport.FillRect(50, 50+vO, 100, 100)

	viewport.SetFillStyle("#888")
	viewport.FillRect(100, 100+vO, 100, 100)

	viewport.SetFillStyle("#333")
	viewport.FillRect(150, 150+vO, 100, 100)

	viewport.SetFillStyle("#444")
	viewport.FillRect(200, 200+vO, 100, 100)

	viewport.SetFillStyle("#555")
	viewport.FillRect(250, 250+vO, 100, 100)

	viewport.SetFillStyle("#666")
	viewport.FillRect(300, 300+vO, 100, 100)

	viewport.SetFillStyle("#777")
	viewport.FillRect(350, 350+vO, 100, 100)

	viewport.SetLineWidth(5)
	viewport.SetStrokeStyle("#0F0")
	viewport.StrokeRect(x-32, y-32, 64, 64)
}
