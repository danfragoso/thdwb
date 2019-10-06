package mustard

import (
	"github.com/danfragoso/thdwb/structs"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tfriedel6/canvas"
)

//Input "Creates a new Input element"
func Input(id string, w float64, h float64, cv *canvas.Canvas) structs.UIElement {
	elementCursor := glfw.CreateStandardCursor(glfw.IBeamCursor)
	inputElement := structs.UIElement{ID: id, WType: "input", X: w/2 - w/4, Y: 10, W: w / 2, H: 30, Canvas: cv, Cursor: elementCursor}
	structs.DrawInput(inputElement)
	return inputElement
}

//Box "Creates a new Box element"
func Box(id string, x float64, y float64, w float64, h float64, cv *canvas.Canvas) structs.UIElement {
	elementCursor := glfw.CreateStandardCursor(glfw.ArrowCursor)
	boxElement := structs.UIElement{ID: id, WType: "box", X: x, Y: y, W: w, H: h, Canvas: cv, Cursor: elementCursor}
	structs.DrawBox(boxElement)
	return boxElement
}

func getFocusedUIElement(eList []*structs.UIElement, x float64, y float64) *structs.UIElement {
	var focusedElement *structs.UIElement

	for i := 0; i < len(eList); i++ {
		if x > eList[i].X && x < eList[i].X+eList[i].W && y > eList[i].Y && y < eList[i].Y+eList[i].H {
			focusedElement = eList[i]
		}
	}

	return focusedElement
}
