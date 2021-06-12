package hotdog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Settings struct {
	Homepage string `json:"homepage"`

	WindowWidth  int `json:"windowWidth"`
	WindowHeight int `json:"windowHeight"`

	HiDPI bool `json:"hiDPI"`

	ExperimentalHTML   bool `json:"experimentalHTML"`   // golang stdlib spec compliant html5 parser
	ExperimentalLayout bool `json:"experimentalLayout"` // spec-ish compliant layout engine
}

var defaultSettings Settings = Settings{
	Homepage: "thdwb://homepage",

	WindowWidth:  600,
	WindowHeight: 600,

	HiDPI: true,

	ExperimentalHTML:   false,
	ExperimentalLayout: false,
}

func LoadSettings(path string) *Settings {
	settingsData, err := ioutil.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(settingsData, &defaultSettings)
		if err != nil {
			fmt.Println("Error loading settings from file;", err)
		}
	} else {
		fmt.Println("Unable to read settings file", path)
		fmt.Println("Loading default settings.")
	}

	return &defaultSettings
}
