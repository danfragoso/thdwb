package mustard

import (
	"github.com/danfragoso/thdwb/ketchup"
	"github.com/danfragoso/thdwb/sauce"
	"github.com/danfragoso/thdwb/structs"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func attachBrowserWindowEvents(browserWindow *structs.AppWindow) {
	browserWindow.GlfwWindow.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		shouldFireKeyEvent := action != 0

		if shouldFireKeyEvent {
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
		}
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
