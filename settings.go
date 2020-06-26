package main

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	Homepage     string `json:"homepage"`
	WindowWidth  int    `json:"windowWidth"`
	WindowHeight int    `json:"windowHeight"`
}

func LoadSettings(path string) *Settings {
	var settings Settings
	settingsData, err := ioutil.ReadFile(path)

	if err == nil {
		json.Unmarshal(settingsData, &settings)
	}

	if settings.Homepage == "" {
		settings.Homepage = "thdwb://homepage"
	}

	if settings.WindowWidth == 0 {
		settings.WindowWidth = 600
	}

	if settings.WindowHeight == 0 {
		settings.WindowHeight = 600
	}

	return &settings
}
