package mustard

func (window *Window) ProcessPointerPosition() {
	window.glw.SetCursor(window.defaultCursor)
	window.ProcessButtons()
	window.ProcessInputs()
}

func (window *Window) ProcessPointerClick() {
	window.ProcessButtonClick()
	window.ProcessInputActivation()
}

func (window *Window) ProcessReturnKey() {
	if window.activeInput != nil && window.activeInput.active == true {
		window.activeInput.active = false
		window.activeInput.selected = false
		window.RequestRepaint()
		window.activeInput.returnCallback()
		window.activeInput = nil
	}
}

func (window *Window) ProcessButtons() {
	x, y := window.cursorX, window.cursorY

	for _, button := range window.registeredButtons {
		if x > float64(button.box.left)+button.padding &&
			x < float64(button.box.left+button.box.width)-button.padding &&
			y > float64(button.box.top)+button.padding &&
			y < float64(button.box.top+button.box.height)-button.padding {
			button.selected = true
			window.glw.SetCursor(button.cursor)
			break
		} else {
			button.selected = false
		}
	}
}

func (window *Window) ProcessInputs() {
	x, y := window.cursorX, window.cursorY

	for _, input := range window.registeredInputs {
		if x > float64(input.box.left)+input.padding &&
			x < float64(input.box.left+input.box.width)-input.padding &&
			y > float64(input.box.top)+input.padding &&
			y < float64(input.box.top+input.box.height)-input.padding {
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
