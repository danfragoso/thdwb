package mayo

import (
	"testing"

	"github.com/danfragoso/thdwb/structs"
	"github.com/stretchr/testify/assert"
)

func TestMapCSSColor(t *testing.T) {
	blues := []string{"#0000FF", "blue", "#00f", "rgba(0, 0, 255)", "rgba(0%, 0%, 100%)"}

	for _, blue := range blues {
		result := MapCSSColor(blue)
		expected := &structs.ColorRGBA{R: 0, G: 0, B: 1, A: 1}
		assert.Equal(t, result, expected, "Expecting: &structs.ColorRGBA{0, 0, 1, 1}")
	}
}
