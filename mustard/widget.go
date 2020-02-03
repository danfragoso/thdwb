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
		case *ImageWidget:
			widget := widgets[i].(*ImageWidget)
			coreWidgets = append(coreWidgets, &widget.widget)
		case *ContextWidget:
			widget := widgets[i].(*ContextWidget)
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
				currentLayout.left = childrenLayout[i-1].left + childrenLayout[i-1].width
			} else {
				currentLayout.left = left
			}

			if children[i].fixedWidth {
				currentLayout.width = children[i].width
			} else {
				currentLayout.width = remainingWidth / len(volatileWidthElements)
			}

			currentLayout.top = top
			currentLayout.height = height
		} else {
			fixedHeightElements, volatileHeightElements := getFixedHeightElements(children)
			remainingHeight := calculateFlexibleHeight(height, fixedHeightElements)

			if i > 0 {
				currentLayout.top = childrenLayout[i-1].top + childrenLayout[i-1].height
			} else {
				currentLayout.top = top
			}

			if children[i].fixedHeight {
				currentLayout.height = children[i].height
			} else {
				currentLayout.height = remainingHeight / len(volatileHeightElements)
			}

			currentLayout.left = left
			currentLayout.width = width
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
		avaiableWidth = avaiableWidth - el.width
	}

	return avaiableWidth
}

func calculateFlexibleHeight(avaiableHeight int, elements []*widget) int {
	for _, el := range elements {
		avaiableHeight = avaiableHeight - el.height
	}

	return avaiableHeight
}
