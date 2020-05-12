package bun

import (
	"thdwb/assets"

	"github.com/goki/freetype/truetype"
)

var sansSerif = map[int]*truetype.Font{
	300: parseFont(300),
	400: parseFont(400),
	600: parseFont(600),
	700: parseFont(700),
	800: parseFont(800),
}

func parseFont(weight int) *truetype.Font {
	font, _ := truetype.Parse(assets.OpenSans(weight))
	return font
}
