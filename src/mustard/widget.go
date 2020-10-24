package mustard

func getCoreWidgets(widgets []Widget) []*baseWidget {
	var coreWidgets []*baseWidget

	for i := 0; i < len(widgets); i++ {
		switch widgets[i].(type) {
		case *Frame:
			widget := widgets[i].(*Frame)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *LabelWidget:
			widget := widgets[i].(*LabelWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *TextWidget:
			widget := widgets[i].(*TextWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *ImageWidget:
			widget := widgets[i].(*ImageWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *CanvasWidget:
			widget := widgets[i].(*CanvasWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *ButtonWidget:
			widget := widgets[i].(*ButtonWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *InputWidget:
			widget := widgets[i].(*InputWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *ScrollBarWidget:
			widget := widgets[i].(*ScrollBarWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		}
	}
	return coreWidgets
}

func getBoreWidgets(widgets []interface{}) []*baseWidget {
	var coreWidgets []*baseWidget

	for i := 0; i < len(widgets); i++ {
		switch widgets[i].(type) {
		case *Frame:
			widget := widgets[i].(*Frame)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *LabelWidget:
			widget := widgets[i].(*LabelWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *TextWidget:
			widget := widgets[i].(*TextWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *ImageWidget:
			widget := widgets[i].(*ImageWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *CanvasWidget:
			widget := widgets[i].(*CanvasWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *ButtonWidget:
			widget := widgets[i].(*ButtonWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *InputWidget:
			widget := widgets[i].(*InputWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		case *ScrollBarWidget:
			widget := widgets[i].(*ScrollBarWidget)
			coreWidgets = append(coreWidgets, &widget.baseWidget)
		}
	}
	return coreWidgets
}

func calculateChildrenWidgetsLayout(children []*baseWidget, top, left, width, height int, orientation FrameOrientation) []*baseWidget {
	var childrenLayout []*baseWidget

	childrenLen := len(children)
	for i := 0; i < childrenLen; i++ {
		currentLayout := &baseWidget{}

		if orientation == VerticalFrame {
			fixedWidthElemens, volatileWidthElements := getFixedWidthElements(children)
			remainingWidth := calculateFlexibleWidth(width, fixedWidthElemens)

			if i > 0 {
				currentLayout.box.left = childrenLayout[i-1].box.left + childrenLayout[i-1].box.width
			} else {
				currentLayout.box.left = left
			}

			if children[i].fixedWidth {
				currentLayout.box.width = children[i].box.width
			} else {
				currentLayout.box.width = remainingWidth / len(volatileWidthElements)
			}

			currentLayout.box.top = top
			currentLayout.box.height = height
		} else {
			fixedHeightElements, volatileHeightElements := getFixedHeightElements(children)
			remainingHeight := calculateFlexibleHeight(height, fixedHeightElements)

			if i > 0 {
				currentLayout.box.top = childrenLayout[i-1].box.top + childrenLayout[i-1].box.height
			} else {
				currentLayout.box.top = top
			}

			if children[i].fixedHeight {
				currentLayout.box.height = children[i].box.height
			} else {
				currentLayout.box.height = remainingHeight / len(volatileHeightElements)
			}

			currentLayout.box.left = left
			currentLayout.box.width = width
		}

		childrenLayout = append(childrenLayout, currentLayout)
	}

	return childrenLayout
}

func getFixedWidthElements(elements []*baseWidget) ([]*baseWidget, []*baseWidget) {
	var fixedWidth []*baseWidget
	var volatileWidth []*baseWidget

	for _, element := range elements {
		if element.fixedWidth {
			fixedWidth = append(fixedWidth, element)
		} else {
			volatileWidth = append(volatileWidth, element)
		}
	}
	return fixedWidth, volatileWidth
}

func getFixedHeightElements(elements []*baseWidget) ([]*baseWidget, []*baseWidget) {
	var fixedHeight []*baseWidget
	var volatileHeight []*baseWidget

	for _, element := range elements {
		if element.fixedHeight {
			fixedHeight = append(fixedHeight, element)
		} else {
			volatileHeight = append(volatileHeight, element)
		}
	}
	return fixedHeight, volatileHeight
}

func calculateFlexibleWidth(avaiableWidth int, elements []*baseWidget) int {
	for _, el := range elements {
		avaiableWidth = avaiableWidth - el.box.width
	}

	return avaiableWidth
}

func calculateFlexibleHeight(avaiableHeight int, elements []*baseWidget) int {
	for _, el := range elements {
		avaiableHeight = avaiableHeight - el.box.height
	}

	return avaiableHeight
}

func (widget *baseWidget) RequestReflow() {
	if widget.window != nil {
		widget.window.needsReflow = true
	}
}

func (widget *baseWidget) RequestRepaint() {
	widget.needsRepaint = true
}

func (widget *baseWidget) GetRect() (int, int, int, int) {
	return widget.computedBox.top, widget.computedBox.left, widget.computedBox.width, widget.computedBox.height
}

func (widget *baseWidget) GetTop() int {
	return widget.computedBox.top
}

func (widget *baseWidget) GetLeft() int {
	return widget.computedBox.left
}

func (widget *baseWidget) GetWidth() int {
	return widget.computedBox.width
}

func (widget *baseWidget) GetHeight() int {
	return widget.computedBox.height
}

func (widget *baseWidget) ComputedBox() *box {
	return &widget.computedBox
}

func (widget *baseWidget) IsPointInside(x, y float64) bool {
	if widget.window.hasActiveOverlay {
		return false
	}

	top, left, width, height := widget.GetRect()
	return x > float64(left) && x < float64(left+width) && y > float64(top) && y < float64(top+height)
}
