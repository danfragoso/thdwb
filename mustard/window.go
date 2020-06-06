package mustard

import (
	"image"
	"image/draw"
	"log"
	"os"

	assets "thdwb/assets"
	gg "thdwb/gg"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/goki/freetype/truetype"
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
	window.RecreateOverlayContext()
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

func (window *Window) RecreateOverlayContext() {
	window.overlayContext = gg.NewContext(window.width, window.height)
	window.hasOverlay = true
}

func (window *Window) ClearOverlayContext() {
	window.hasOverlay = false
	window.overlayContext = nil
}

func (window *Window) ClearMenuEntries() {
	window.menuEntries = nil
}

func (window *Window) AddMenuEntry(entryText string, callback func()) {
	window.menuEntries = append(window.menuEntries, &menuEntry{entryText, callback})
}

func (window *Window) addEvents() {
	window.glw.SetFocusCallback(func(w *glfw.Window, focused bool) {
	})

	window.glw.SetSizeCallback(func(w *glfw.Window, width, height int) {
		window.width, window.height = width, height
		window.RecreateContext()
		window.RecreateOverlayContext()

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
			if window.activeInput != nil && action == glfw.Release {
				window.activeInput.active = false
				window.activeInput.selected = false
				window.activeInput.needsRepaint = true
				window.activeInput = nil
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
		case glfw.KeyEnter:
			if action == glfw.Release {
				window.ProcessReturnKey()
			}
			break
		}
	})

	window.glw.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Release {
			if button == glfw.MouseButtonLeft {
				window.ProcessPointerClick()
				window.DestroyContextMenu()
			} else if button == glfw.MouseButtonRight {
				window.CreateContextMenu()
			}
		}
	})

	window.glw.SetScrollCallback(func(w *glfw.Window, x, y float64) {
		window.ProcessScroll(x, y)
	})
}

func (window *Window) generateTexture() {
	gl.DeleteTextures(1, &window.backend.texture)
	if window.hasOverlay {
		ctx, _ := window.context.Image().(draw.Image)
		overlay, _ := window.overlayContext.Image().(draw.Image)

		window.frameBuffer = image.NewRGBA(ctx.Bounds())

		draw.Draw(window.frameBuffer, ctx.Bounds(), ctx, image.ZP, draw.Over)
		draw.Draw(window.frameBuffer, overlay.Bounds(), overlay, image.ZP, draw.Over)
	} else {
		window.frameBuffer = window.context.Image().(*image.RGBA)
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

func (window *Window) AttachClickEventListener(callback func()) {
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

func (window *Window) CreateContextMenu() {
	window.RecreateOverlayContext()
	ctx := window.overlayContext

	menuTop := window.cursorY
	menuLeft := window.cursorX
	menuWidth := 200.
	menuHeight := float64(len(window.menuEntries) * 20)

	if menuLeft+menuWidth > float64(window.width) {
		menuLeft = float64(window.width) - menuWidth
	}

	if menuTop+menuHeight > float64(window.height) {
		menuTop = float64(window.height) - menuHeight
	}

	ctx.SetHexColor("#eee")
	ctx.DrawRectangle(menuLeft, menuTop, menuWidth, menuHeight)
	ctx.Fill()

	font, _ := truetype.Parse(assets.OpenSans(400))
	ctx.SetHexColor("#222")
	ctx.SetFont(font, 16)

	for idx, entry := range window.menuEntries {
		ctx.DrawString(prepEntry(ctx, entry.entryText, menuWidth), menuLeft, menuTop+16+float64(idx*20))
		ctx.Fill()
	}

	ctx.DrawRectangle(menuLeft, menuTop, menuWidth, menuHeight)
	ctx.SetHexColor("#ddd")
	ctx.Stroke()
}

func prepEntry(ctx *gg.Context, entry string, width float64) string {
	w, _ := ctx.MeasureString(entry)

	if w < width {
		return entry
	}

	for i := 0; i < len(entry); i++ {
		nW, _ := ctx.MeasureString(entry[:len(entry)-i] + "...")

		if nW <= width {
			return entry[:len(entry)-i] + "..."
		}
	}

	return entry
}

func (window *Window) DestroyContextMenu() {
	window.ClearOverlayContext()
}
