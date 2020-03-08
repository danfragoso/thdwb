package mustard

func getCoreWidgets(widgets []interface{}) []*widget {
	var coreWidgets []*widget

	for i := 0; i < len(widgets); i++ {
		switch widgets[i].(type) {
		case *Frame:
			widget := widgets[i].(*Frame)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *LabelWidget:
			widget := widgets[i].(*LabelWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *TextWidget:
			widget := widgets[i].(*TextWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *ImageWidget:
			widget := widgets[i].(*ImageWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *ContextWidget:
			widget := widgets[i].(*ContextWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *ButtonWidget:
			widget := widgets[i].(*ButtonWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *InputWidget:
			widget := widgets[i].(*InputWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		}
	}
	return coreWidgets
}

func calculateChildrenWidgetsLayout(children []*widget, top, left, width, height int, orientation FrameOrientation) []*widget {
	var childrenLayout []*widget

	childrenLen := len(children)
	for i := 0; i < childrenLen; i++ {
		currentLayout := &widget{}

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

func getFixedWidthElements(elements []*widget) ([]*widget, []*widget) {
	var fixedWidth []*widget
	var volatileWidth []*widget

	for _, element := range elements {
		if element.fixedWidth {
			fixedWidth = append(fixedWidth, element)
		} else {
			volatileWidth = append(volatileWidth, element)
		}
	}
	return fixedWidth, volatileWidth
}

func getFixedHeightElements(elements []*widget) ([]*widget, []*widget) {
	var fixedHeight []*widget
	var volatileHeight []*widget

	for _, element := range elements {
		if element.fixedHeight {
			fixedHeight = append(fixedHeight, element)
		} else {
			volatileHeight = append(volatileHeight, element)
		}
	}
	return fixedHeight, volatileHeight
}

func calculateFlexibleWidth(avaiableWidth int, elements []*widget) int {
	for _, el := range elements {
		avaiableWidth = avaiableWidth - el.box.width
	}

	return avaiableWidth
}

func calculateFlexibleHeight(avaiableHeight int, elements []*widget) int {
	for _, el := range elements {
		avaiableHeight = avaiableHeight - el.box.height
	}

	return avaiableHeight
}

func (widget *widget) RequestReflow() {
	if widget.window != nil {
		widget.window.needsReflow = true
	}
}

func (widget *widget) RequestRepaint() {
	widget.needsRepaint = true
}
