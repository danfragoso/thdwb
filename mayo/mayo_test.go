package mayo

import (
	"testing"

	"github.com/danfragoso/thdwb/structs"
	"github.com/stretchr/testify/assert"
)

func TestHexStringToColor(t *testing.T) {
	validColors := []string{"#0000FF", "#00f", "#00ff", "#0000ffff"}
	invalidColors := []string{"#00000", "#d0"}

	for _, validColor := range validColors {
		result := HexStringToColor(validColor)
		expected := &structs.ColorRGBA{R: 0, G: 0, B: 1, A: 1}
		assert.Equal(t, expected, result, "Expecting: &structs.ColorRGBA{0, 0, 1, 1}")
	}

	for _, invalidColor := range invalidColors {
		result := HexStringToColor(invalidColor)
		expected := &structs.ColorRGBA{}
		assert.Equal(t, expected, result, "Expecting: &structs.ColorRGBA{}")
	}
}

func TestRGBAToColor(t *testing.T) {
	blues := []string{"rgba(0, 0, 255)", "rgba(0%, 0%, 100%)",
		"rgb(0, 0, 255, 1)", "rgb(0%, 0%, 100%, 100%)"}

	for _, blue := range blues {
		result := RGBAToColor(blue)
		expected := &structs.ColorRGBA{R: 0, G: 0, B: 1, A: 1}
		assert.Equal(t, expected, result, "Expecting: &structs.ColorRGBA{0, 0, 1, 1}")
	}
}

func TestMapCSSColor(t *testing.T) {
	blues := []string{"#0000FF", "blue", "#00f", "rgba(0, 0, 255)", "rgba(0%, 0%, 100%)",
		"rgb(0, 0, 255, 1)", "rgb(0%, 0%, 100%)"}

	for _, blue := range blues {
		result := MapCSSColor(blue)
		expected := &structs.ColorRGBA{R: 0, G: 0, B: 1, A: 1}
		assert.Equal(t, expected, result, "Expecting: &structs.ColorRGBA{0, 0, 1, 1}")
	}
}
