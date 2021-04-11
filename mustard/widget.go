package mustard

import "image"

func getCoreWidgets(widgets []Widget) []*baseWidget {
	var coreWidgets []*baseWidget
	for _, widget := range widgets {
		coreWidgets = append(coreWidgets, widget.BaseWidget())
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

	if avaiableWidth < 0 {
		return 0
	}

	return avaiableWidth
}

func calculateFlexibleHeight(avaiableHeight int, elements []*baseWidget) int {
	for _, el := range elements {
		avaiableHeight = avaiableHeight - el.box.height
	}

	if avaiableHeight < 0 {
		return 0
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

func (widget *baseWidget) SetWindow(window *Window) {
	widget.window = window

	for _, childWidget := range widget.widgets {
		childWidget.SetWindow(window)
	}
}

func (widget *baseWidget) Buffer() *image.RGBA {
	return widget.buffer
}

func (widget *baseWidget) Widgets() []Widget {
	return widget.widgets
}

func (widget *baseWidget) AttachWidget(wd Widget) {
	wd.SetWindow(widget.window)
	widget.widgets = append(widget.widgets, wd)

	if widget.window != nil && widget.window.rootFrame != nil {
		widget.window.rootFrame.RequestReflow()
	}
}

func (widget *baseWidget) DetachWidget(wd Widget) Widget {
	var detachedWidget Widget
	var childWidgets []Widget

	for _, childWidget := range widget.widgets {
		if childWidget == wd {
			detachedWidget = childWidget
		} else {
			childWidgets = append(childWidgets, childWidget)
		}
	}

	widget.widgets = childWidgets
	if widget.window != nil && widget.window.rootFrame != nil {
		widget.window.rootFrame.RequestReflow()
	}

	return detachedWidget
}

func (widget *baseWidget) BaseWidget() *baseWidget {
	return widget
}

func (widget *baseWidget) NeedsRepaint() bool {
	return widget.needsRepaint
}

func (widget *baseWidget) SetNeedsRepaint(value bool) {
	widget.needsRepaint = value
}

func (widget *baseWidget) IsPointInside(x, y float64) bool {
	if widget.window.hasActiveOverlay {
		return false
	}

	top, left, width, height := widget.GetRect()
	return x > float64(left) && x < float64(left+width) && y > float64(top) && y < float64(top+height)
}
