package mustard

import (
	"image"
	"image/draw"
	"log"
	"os"

	gg "thdwb/gg"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func SetGLFWHints() {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
}

func CreateNewWindow(title string, width int, height int) *Window {
	glw, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window := &Window{
		title:  title,
		width:  width,
		height: height,
		glw:    glw,

		defaultCursor: glfw.CreateStandardCursor(glfw.ArrowCursor),
		pointerCursor: glfw.CreateStandardCursor(glfw.HandCursor),
	}

	window.RecreateContext()
	glw.MakeContextCurrent()

	window.backend = createGLBackend()
	window.addEvents()
	window.generateTexture()
	return window
}

//Show - Show the window
func (window *Window) Show() {
	window.needsReflow = true
	window.visible = true
	window.glw.Show()
}

//SetRootFrame - Sets the window root frame
func (window *Window) SetRootFrame(frame *Frame) {
	window.rootFrame = frame
}

//SetRootFrame - Sets the window root frame
func (window *Window) GetSize() (int, int) {
	return window.width, window.height
}

func (window *Window) brocessFrame() {
	if window.needsReflow {

		window.needsReflow = false
		window.glw.MakeContextCurrent()

		drawRootFrame(window)
		window.generateTexture()

		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		window.glw.SwapBuffers()
	}

	glfw.WaitEvents()
}

func (window *Window) processFrame() {
	if window.needsReflow {
		window.needsReflow = false
		window.glw.MakeContextCurrent()

		drawRootFrame(window)
		window.generateTexture()

		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		window.glw.SwapBuffers()
	} else {
		redrawWidgets(window.rootFrame)
		window.generateTexture()
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		window.glw.SwapBuffers()
	}

	glfw.PollEvents()
}

func (window *Window) RequestReflow() {
	window.needsReflow = true
}

func (window *Window) RecreateContext() {
	window.context = nil
	window.context = gg.NewContext(window.width, window.height)

	window.context.SetRGB(0, 0, 0)
	window.context.Clear()

	window.context.SetRGB(1, 1, 1)
}

func (window *Window) addEvents() {
	window.glw.SetFocusCallback(func(w *glfw.Window, focused bool) {
	})

	window.glw.SetSizeCallback(func(w *glfw.Window, width, height int) {
		window.width, window.height = width, height
		window.RecreateContext()
		//window.RecreateOverlayContext()

		gl.Viewport(0, 0, int32(width), int32(height))
		window.needsReflow = true
	})

	window.glw.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		window.cursorX, window.cursorY = x, y
		window.ProcessPointerPosition()
	})

	window.glw.SetCharCallback(func(w *glfw.Window, char rune) {
		if window.activeInput != nil {
			inputVal, cursorPos := window.activeInput.value, window.activeInput.cursorPosition

			window.activeInput.value = inputVal[:len(inputVal)+cursorPos] + string(char) + inputVal[len(inputVal)+cursorPos:]
			window.activeInput.needsRepaint = true
		}
	})

	window.glw.SetCloseCallback(func(w *glfw.Window) {
		os.Exit(0)
	})

	window.glw.SetKeyCallback(func(w *glfw.Window, key glfw.Key, sc int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyBackspace:
			if action == glfw.Repeat || action == glfw.Release {
				if window.activeInput != nil && len(window.activeInput.value) > 0 {
					if window.activeInput.cursorPosition == 0 {
						window.activeInput.value = window.activeInput.value[:len(window.activeInput.value)-1]
					} else {
						inputVal, cursorPos := window.activeInput.value, window.activeInput.cursorPosition

						if cursorPos+len(inputVal) > 0 {
							window.activeInput.value = inputVal[:len(inputVal)+cursorPos-1] + inputVal[len(inputVal)+cursorPos:]
						}
					}
					window.activeInput.needsRepaint = true
				}
			}
			break
		case glfw.KeyEscape:
			if action == glfw.Release {
				window.DestroyContextMenu()

				if window.activeInput != nil {
					window.activeInput.active = false
					window.activeInput.selected = false
					window.activeInput.needsRepaint = true
					window.activeInput = nil
				}
			}

			break
		case glfw.KeyUp:
			if action == glfw.Release || action == glfw.Repeat {
				window.ProcessArrowKeys("up")
			}
			break
		case glfw.KeyDown:
			if action == glfw.Release || action == glfw.Repeat {
				window.ProcessArrowKeys("down")
			}
			break
		case glfw.KeyLeft:
			if action == glfw.Release || action == glfw.Repeat {
				window.ProcessArrowKeys("left")
			}
			break
		case glfw.KeyRight:
			if action == glfw.Release || action == glfw.Repeat {
				window.ProcessArrowKeys("right")
			}
			break
		case glfw.KeyV:
			if action == glfw.Release && mods == glfw.ModControl {
				if window.activeInput != nil {
					if window.activeInput.cursorPosition == 0 {
						window.activeInput.value = window.activeInput.value + glfw.GetClipboardString()
					} else {
						inputVal, cursorPos := window.activeInput.value, window.activeInput.cursorPosition
						window.activeInput.value = inputVal[:len(inputVal)+cursorPos] + glfw.GetClipboardString() + inputVal[len(inputVal)+cursorPos:]
					}
					window.activeInput.needsRepaint = true
				}
			}
			break
		case glfw.KeyEnter:
			if action == glfw.Release {
				window.ProcessReturnKey()
			}
			break
		}
	})

	window.glw.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Release {
			window.ProcessPointerClick(button)
		}
	})

	window.glw.SetScrollCallback(func(w *glfw.Window, x, y float64) {
		window.ProcessScroll(x, y)
	})
}

func (window *Window) generateTexture() {
	gl.DeleteTextures(1, &window.backend.texture)
	window.frameBuffer = window.context.Image().(*image.RGBA)

	if window.hasActiveOverlay {
		nBuffer := image.NewRGBA(window.frameBuffer.Bounds())
		draw.Draw(nBuffer, window.frameBuffer.Bounds(), window.frameBuffer, image.ZP, draw.Over)
		window.frameBuffer = nBuffer

		for _, overlay := range window.overlays {
			draw.Draw(window.frameBuffer, overlay.buffer.Bounds().Add(overlay.position), overlay.buffer, image.ZP, draw.Over)
		}
	}

	gl.GenTextures(1, &window.backend.texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, window.backend.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(window.frameBuffer.Rect.Size().X), int32(window.frameBuffer.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(window.frameBuffer.Pix),
	)

}

func (window *Window) GetCursorPosition() (float64, float64) {
	return window.cursorX, window.cursorY
}

func (window *Window) RegisterButton(button *ButtonWidget, callback func()) {
	button.onClick = callback
	window.registeredButtons = append(window.registeredButtons, button)
}

func (window *Window) RegisterInput(input *InputWidget) {
	window.registeredInputs = append(window.registeredInputs, input)
}

func (window *Window) AttachPointerPositionEventListener(callback func(pointerX, pointerY float64)) {
	window.pointerPositionEventListeners = append(window.pointerPositionEventListeners, callback)
}

func (window *Window) AttachScrollEventListener(callback func(direction int)) {
	window.scrollEventListeners = append(window.scrollEventListeners, callback)
}

func (window *Window) AttachClickEventListener(callback func(MustardKey)) {
	window.clickEventListeners = append(window.clickEventListeners, callback)
}

func (window *Window) SetCursor(cursorType string) {
	switch cursorType {
	case "pointer":
		window.glw.SetCursor(window.pointerCursor)
		break

	default:
		window.glw.SetCursor(window.defaultCursor)
	}
}

func (window *Window) AddOverlay(overlay *Overlay) {
	window.overlays = append(
		window.overlays,
		overlay,
	)

	window.hasActiveOverlay = true
}

func (window *Window) RemoveOverlay(overlay *Overlay) {
	for idx, cOverlay := range window.overlays {
		if cOverlay == overlay {
			window.overlays = append(window.overlays[:idx], window.overlays[idx+1:]...)
			break
		}
	}

	if len(window.overlays) < 1 {
		window.hasActiveOverlay = false
	}
}
