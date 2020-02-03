package mustard

import "github.com/go-gl/gl/v4.1-core/gl"

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
}
