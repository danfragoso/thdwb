package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Settings struct {
	Homepage     string `json:"homepage"`
	WindowWidth  int    `json:"windowWidth"`
	WindowHeight int    `json:"windowHeight"`

	HiDPI bool `json:"hiDPI"`

	// golang stdlib spec compliant html5 parser
	ExperimentalHTML bool `json:"experimentalHTML"`
	// spec-ish compliant layout engine
	ExperimentalLayout bool `json:"experimentalLayout"`
}

func LoadSettings(path string) *Settings {
	var settings Settings
	settingsData, err := ioutil.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(settingsData, &settings)
		if err != nil {
			fmt.Println("Error loading settings from file;", err)
		}
	} else {
		fmt.Println("Unable to read settings file", path)
		fmt.Println("Loading default settings.")
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

	// Default to `false` for HiDPI, experimentalHTML, and experimentalLayout nothing to do here.
	return &settings
}
