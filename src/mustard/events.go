package mustard

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func (window *Window) ProcessPointerPosition() {
	if window.hasActiveOverlay {
		window.ProcessContextMenu()
	} else {
		window.selectWidget(window.rootFrame)
		if window.selectedWidget != nil &&
			window.selectedWidget.BaseWidget().widgetType == frameWidget {
			window.SetCursor(DefaultCursor)
		}

		window.ProcessButtons()
		window.ProcessTrees()
		window.ProcessInputs()
		window.firePointerPositionEvents()
	}
}

func (window *Window) selectWidget(widget Widget) {
	if widget.BaseWidget().IsPointInside(window.cursorX, window.cursorY) {
		window.selectedWidget = widget

		for _, childWidget := range widget.Widgets() {
			window.selectWidget(childWidget)
		}
	}
}

func (window *Window) ProcessPointerClick(button glfw.MouseButton) {
	if window.hasActiveOverlay {
		window.ProcessContextMenuClick()
	} else {
		if button == glfw.MouseButtonLeft {
			window.ProcessButtonClick()
			window.ProcessTreeClick()
			window.ProcessInputActivation()
			window.fireClickEvents(MouseLeft)
		} else if button == glfw.MouseButtonRight {
			window.fireClickEvents(MouseRight)
		}
	}
}

func (window *Window) ProcessScroll(x, y float64) {
	window.fireScrollEvents(x, y)
}

func (window *Window) ProcessReturnKey() {
	if window.activeInput != nil && window.activeInput.active == true {
		window.activeInput.active = false

		window.activeInput.needsRepaint = true
		window.activeInput.returnCallback()
		window.activeInput = nil
	}
}

func (window *Window) ProcessArrowKeys(arrowKey string) {
	if window.activeInput != nil && window.activeInput.active == true {
		if arrowKey == "left" && (window.activeInput.cursorPosition+len(window.activeInput.value)) > 0 {
			window.activeInput.cursorPosition--
			window.activeInput.cursorDirection = true
			window.activeInput.needsRepaint = true
		} else if arrowKey == "right" && window.activeInput.cursorPosition < 0 {
			window.activeInput.cursorPosition++
			window.activeInput.cursorDirection = false
			window.activeInput.needsRepaint = true
		}
	}
}

func (window *Window) ProcessTrees() {
	x, y := window.cursorX, window.cursorY

	for _, tree := range window.registeredTrees {
		if x > float64(tree.computedBox.left) &&
			x < float64(tree.computedBox.left+tree.computedBox.width) &&
			y > float64(tree.computedBox.top) &&
			y < float64(tree.computedBox.top+tree.computedBox.height) {
			tree.selected = true
			tree.needsRepaint = true
			window.glw.SetCursor(tree.cursor)
			break
		} else {
			tree.selected = false
			tree.needsRepaint = true
		}
	}
}

func (window *Window) ProcessButtons() {
	x, y := window.cursorX, window.cursorY

	for _, button := range window.registeredButtons {
		if x > float64(button.computedBox.left)+button.padding &&
			x < float64(button.computedBox.left+button.computedBox.width)-button.padding &&
			y > float64(button.computedBox.top)+button.padding &&
			y < float64(button.computedBox.top+button.computedBox.height)-button.padding {
			button.selected = true
			button.needsRepaint = true
			window.glw.SetCursor(button.cursor)
			break
		} else {
			button.selected = false
			button.needsRepaint = true
		}
	}
}

func (window *Window) ProcessInputs() {
	x, y := window.cursorX, window.cursorY

	for _, input := range window.registeredInputs {
		if x > float64(input.computedBox.left)+input.padding &&
			x < float64(input.computedBox.left+input.computedBox.width)-input.padding &&
			y > float64(input.computedBox.top)+input.padding &&
			y < float64(input.computedBox.top+input.computedBox.height)-input.padding {
			input.selected = true
			input.needsRepaint = true
			window.glw.SetCursor(input.cursor)
			break
		} else {
			input.selected = false
			input.needsRepaint = true
		}
	}
}

func (window *Window) ProcessInputActivation() {
	for _, input := range window.registeredInputs {
		if input.selected == true {
			window.activeInput = input
			input.active = true
			input.needsRepaint = true
			return
		}

		if input.active {
			input.active = false
			input.needsRepaint = true
			window.activeInput = nil
		}
	}
}

func (window *Window) ProcessButtonClick() {
	for _, button := range window.registeredButtons {
		if button.selected == true {
			button.onClick()
			button.needsRepaint = true
			return
		}
	}
}

func (window *Window) ProcessTreeClick() {
	for _, tree := range window.registeredTrees {
		if tree.selected == true {
			tree.Click()
			tree.needsRepaint = true
			return
		}
	}
}

func (window *Window) firePointerPositionEvents() {
	for _, eventCallback := range window.pointerPositionEventListeners {
		eventCallback(window.cursorX, window.cursorY)
	}
}

func (window *Window) fireScrollEvents(x, y float64) {
	for _, eventCallback := range window.scrollEventListeners {
		eventCallback(int(y))
	}
}

func (window *Window) fireClickEvents(key MustardKey) {
	for _, eventCallback := range window.clickEventListeners {
		eventCallback(key)
	}
}

func (window *Window) ProcessContextMenu() {
	if len(window.contextMenu.entries) > 0 {
		x, y := window.GetCursorPosition()

		if x > window.contextMenu.overlay.left &&
			x < window.contextMenu.overlay.left+window.contextMenu.overlay.width &&
			y > window.contextMenu.overlay.top &&
			y < window.contextMenu.overlay.top+window.contextMenu.overlay.height {

			for _, entry := range window.contextMenu.entries {
				if entry.PointIntersects(x, y) {
					if entry != window.contextMenu.selectedEntry {
						window.SelectEntry(entry)
					}

					break
				}
			}

		} else {
			window.DeselectEntries()
		}
	}
}

func (window *Window) ProcessContextMenuClick() {
	if window.contextMenu.selectedEntry != nil {
		window.contextMenu.selectedEntry.action()
	}

	window.DestroyContextMenu()
}
