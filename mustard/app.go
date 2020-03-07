package mustard

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func CreateNewApp(name string) *App {
	return &App{name: name}
}

func (app *App) Run(callback func()) {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	for _, window := range app.windows {
		if window.visible && !window.glw.ShouldClose() {
			window.processFrame()
		}
	}

	callback()
	app.Run(callback)
}

func (app *App) AddWindow(window *Window) {
	app.windows = append(app.windows, window)

	setWidgetWindow(&window.rootFrame.widget, window)
}

func setWidgetWindow(widget *widget, window *Window) {
	widget.window = window
	widgets := widget.widgets

	for i := 0; i < len(widgets); i++ {
		switch widgets[i].(type) {
		case *Frame:
			widget := widgets[i].(*Frame)
			setWidgetWindow(&widget.widget, window)
		case *LabelWidget:
			widget := widgets[i].(*LabelWidget)
			setWidgetWindow(&widget.widget, window)
		case *TextWidget:
			widget := widgets[i].(*TextWidget)
			setWidgetWindow(&widget.widget, window)
		case *ImageWidget:
			widget := widgets[i].(*ImageWidget)
			setWidgetWindow(&widget.widget, window)
		case *ContextWidget:
			widget := widgets[i].(*ContextWidget)
			setWidgetWindow(&widget.widget, window)
		case *ButtonWidget:
			widget := widgets[i].(*ButtonWidget)
			setWidgetWindow(&widget.widget, window)
		case *InputWidget:
			widget := widgets[i].(*InputWidget)
			setWidgetWindow(&widget.widget, window)
		}
	}
}
