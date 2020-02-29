package mustard

func (window *Window) ProcessPointerPosition(x, y float64) {
	window.glw.SetCursor(window.defaultCursor)
	window.ProcessButtons(x, y)
	window.ProcessInputs(x, y)
}

func (window *Window) ProcessPointerClick() {
	window.ProcessButtonClick()
	window.ProcessInputActivation()
}

func (window *Window) ProcessButtons(x, y float64) {
	for _, button := range window.registeredButtons {
		if x > float64(button.left)+button.padding &&
			x < float64(button.left+button.width)-button.padding &&
			y > float64(button.top)+button.padding &&
			y < float64(button.top+button.height)-button.padding {
			button.selected = true
			window.glw.SetCursor(button.cursor)
			break
		} else {
			button.selected = false
		}
	}
}

func (window *Window) ProcessInputs(x, y float64) {
	for _, input := range window.registeredInputs {
		if x > float64(input.left)+input.padding &&
			x < float64(input.left+input.width)-input.padding &&
			y > float64(input.top)+input.padding &&
			y < float64(input.top+input.height)-input.padding {
			input.selected = true
			window.glw.SetCursor(input.cursor)
			break
		} else {
			input.selected = false
		}
	}
}

func (window *Window) ProcessInputActivation() {
	for _, input := range window.registeredInputs {
		if input.selected == true {
			window.activeInput = input
			input.active = true
			window.RequestRepaint()
			return
		}

		if input.active {
			input.active = false
			window.activeInput = nil
			window.RequestRepaint()
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
