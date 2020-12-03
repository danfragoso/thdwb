package mayo

import (
	"regexp"

	hotdog "thdwb/hotdog"
)

var rgba = regexp.MustCompile(`rgba?\([\.?\d?\.?\d?%?\s?,?]+\)`)
var rgbaParams = regexp.MustCompile(`\([\.?\d?\.?\d?%?\s?,?]+\)`)

var colorTable = map[string]*hotdog.ColorRGBA{
	"maroon":         {R: 0.5, G: 0.0, B: 0.0, A: 1.0},
	"red":            {R: 1.0, G: 0.0, B: 0.0, A: 1.0},
	"orange":         {R: 1.0, G: 0.6, B: 0.0, A: 1.0},
	"yellow":         {R: 1.0, G: 1.0, B: 0.0, A: 1.0},
	"olive":          {R: 0.5, G: 0.5, B: 0.0, A: 1.0},
	"green":          {R: 0.0, G: 0.5, B: 0.0, A: 1.0},
	"purple":         {R: 0.5, G: 0.0, B: 0.5, A: 1.0},
	"fuchsia":        {R: 1.0, G: 0.0, B: 1.0, A: 1.0},
	"lime":           {R: 0.0, G: 1.0, B: 0.0, A: 1.0},
	"teal":           {R: 0.0, G: 0.5, B: 0.5, A: 1.0},
	"aqua":           {R: 0.0, G: 1.0, B: 1.0, A: 1.0},
	"blue":           {R: 0.0, G: 0.0, B: 1.0, A: 1.0},
	"navy":           {R: 0.0, G: 0.0, B: 0.5, A: 1.0},
	"black":          {R: 0.0, G: 0.0, B: 0.0, A: 1.0},
	"gray":           {R: 0.5, G: 0.5, B: 0.5, A: 1.0},
	"silver":         {R: 0.7, G: 0.7, B: 0.7, A: 1.0},
	"white":          {R: 1.0, G: 1.0, B: 1.0, A: 1.0},
	"tomato":         {R: 1.0, G: 0.38, B: 0.27, A: 1.0},
	"crimson":        {R: 0.8, G: 0.07, B: 0.2, A: 1.0},
	"coral":          {R: 1.0, G: 0.5, B: 0.31, A: 1.0},
	"cornflowerblue": {R: 0.40, G: 0.58, B: 0.92, A: 1.0},
	"darkgreen":      {R: 0.0, G: 0.40, B: 0.0, A: 1.0},
}

var elementFontTable = map[string]float64{
	"h1": float64(32),
	"h2": float64(28),
	"h3": float64(20),
	"p":  float64(14),
}

var elementColorTable = map[string]*hotdog.ColorRGBA{
	"a": {R: 0.0, G: 0.0, B: 1.0, A: 1.0},
}
