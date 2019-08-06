package mayo

import (
	"strconv"
	"strings"

	"github.com/danfragoso/thdwb/structs"
)

var colorTable = map[string]*structs.ColorRGBA{
	"maroon":  &structs.ColorRGBA{R: 0.5, G: 0.0, B: 0.0, A: 1.0},
	"red":     &structs.ColorRGBA{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
	"orange":  &structs.ColorRGBA{R: 1.0, G: 0.6, B: 0.0, A: 1.0},
	"yellow":  &structs.ColorRGBA{R: 1.0, G: 1.0, B: 0.0, A: 1.0},
	"olive":   &structs.ColorRGBA{R: 0.5, G: 0.5, B: 0.0, A: 1.0},
	"green":   &structs.ColorRGBA{R: 0.0, G: 1.0, B: 0.0, A: 1.0},
	"purple":  &structs.ColorRGBA{R: 0.5, G: 0.0, B: 0.5, A: 1.0},
	"fuchsia": &structs.ColorRGBA{R: 1.0, G: 0.0, B: 1.0, A: 1.0},
	"lime":    &structs.ColorRGBA{R: 0.0, G: 1.0, B: 0.0, A: 1.0},
	"teal":    &structs.ColorRGBA{R: 0.0, G: 0.5, B: 0.5, A: 1.0},
	"aqua":    &structs.ColorRGBA{R: 0.0, G: 1.0, B: 1.0, A: 1.0},
	"blue":    &structs.ColorRGBA{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
	"navy":    &structs.ColorRGBA{R: 0.0, G: 0.0, B: 0.5, A: 1.0},
	"black":   &structs.ColorRGBA{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
	"gray":    &structs.ColorRGBA{R: 0.5, G: 0.5, B: 0.5, A: 1.0},
	"silver":  &structs.ColorRGBA{R: 0.7, G: 0.7, B: 0.7, A: 1.0},
	"white":   &structs.ColorRGBA{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
}

var elementFontTable = map[string]float64{
	"h1": float64(32),
	"h2": float64(28),
	"p":  float64(14),
}

func hexValueToFloat(value string) float64 {
	//TODO: round float and fix errors
	n, _ := strconv.ParseInt(value, 16, 0)
	return float64(n) / 15
}

func hexStringToColor(colorString string) *structs.ColorRGBA {
	color := &structs.ColorRGBA{
		R: hexValueToFloat(colorString[1:2]),
		G: hexValueToFloat(colorString[3:4]),
		B: hexValueToFloat(colorString[5:6]),
		A: 1.0,
	}
	return color
}

func mapCSSColor(colorString string) *structs.ColorRGBA {
	var color *structs.ColorRGBA

	if string(colorString[0]) == "#" {
		color = hexStringToColor(colorString)
	} else {
		color = colorTable[colorString]
	}

	return color
}

func mapSizeValue(sizeValue string) float64 {
	valueString := sizeValue[0 : len(sizeValue)-2]
	value, err := strconv.ParseInt(valueString, 10, 0)

	if err != nil {
		return float64(14)
	}

	return float64(value)
}

func mapPropToStylesheet(parsedStyleSheet *structs.Stylesheet, propSlice []string) *structs.Stylesheet {
	propName := propSlice[0]
	propValue := propSlice[1]

	switch propName {
	case "color":
		parsedStyleSheet.Color = mapCSSColor(propValue)
	case "font-size":
		parsedStyleSheet.FontSize = mapSizeValue(propValue)
	}

	return parsedStyleSheet
}

func parseInlineStylesheet(attributes []*structs.Attribute) *structs.Stylesheet {
	parsedStylesheet := &structs.Stylesheet{}

	for i := 0; i < len(attributes); i++ {
		attributeName := attributes[i].Name
		if attributeName == "style" {

			styleString := attributes[i].Value
			styleProps := strings.Split(strings.Replace(styleString, " ", "", -1), ";")

			for i := 0; i < len(styleProps); i++ {
				styledProperty := strings.Split(styleProps[i], ":")
				if len(styledProperty) >= 2 {
					parsedStylesheet = mapPropToStylesheet(parsedStylesheet, styledProperty)
				}
			}
		}
	}

	return parsedStylesheet
}

func hasInlineStyle(attributes []*structs.Attribute) bool {
	inlineStyle := false

	for i := 0; i < len(attributes); i++ {
		attributeName := attributes[i].Name
		if attributeName == "style" {
			inlineStyle = true
		}
	}

	return inlineStyle
}

func GetElementStylesheet(elementName string, attributes []*structs.Attribute) *structs.Stylesheet {
	elementStylesheet := &structs.Stylesheet{&structs.ColorRGBA{0, 0, 0, 0}, 0}

	if hasInlineStyle(attributes) {
		elementStylesheet = parseInlineStylesheet(attributes)
	}

	if elementStylesheet.FontSize == float64(0) {
		fontSize := elementFontTable[elementName]
		if fontSize != float64(0) {
			elementStylesheet.FontSize = fontSize
		} else {
			elementStylesheet.FontSize = float64(14)
		}
	}

	return elementStylesheet
}
