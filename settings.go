package main

type Settings struct {
	homepage     string
	windowWidth  int
	windowHeight int
}

func LoadSettings(path string) *Settings {
	return &Settings{}
}

func (settings *Settings) Homepage() string {
	if settings.homepage == "" {
		return "thdwb://homepage/"
	}

	return settings.homepage
}

func (settings *Settings) WindowWidth() int {
	if settings.windowWidth == 0 {
		return 600
	}

	return settings.windowWidth
}

func (settings *Settings) WindowHeight() int {
	if settings.windowHeight == 0 {
		return 600
	}

	return settings.windowHeight
}
