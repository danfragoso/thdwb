package mustard

func (window *Window) ProcessPointerPosition(x, y float64) {
	window.ProcessButtons(x, y)
}

func (window *Window) ProcessPointerClick() {
	window.ProcessButtonClick()
}

func (window *Window) ProcessButtons(x, y float64) {
	for _, button := range window.registeredButtons {
		if x > float64(button.left) && x < float64(button.left+button.width) && y > float64(button.top) && y < float64(button.top+button.height) {
			button.selected = true
			window.glw.SetCursor(button.cursor)
		} else {
			button.selected = false
			window.glw.SetCursor(window.defaultCursor)
		}
	}
}

func (window *Window) ProcessButtonClick() {
	for _, button := range window.registeredButtons {
		if button.selected == true {
			button.onClick()
			window.RequestRepaint()
			return
		}
	}
}
