package mayo

import (
	"strconv"
	"strings"

	"github.com/danfragoso/thdwb/structs"
)

func getDefaultElementDisplay(element string) string {
	displayType := "block"

	switch element {
	case "script", "style", "meta", "link", "head", "title":
		displayType = "none"
	default:
		displayType = "block"
	}

	return displayType
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

func MapCSSColor(colorString string) *structs.ColorRGBA {
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
		parsedStyleSheet.Color = MapCSSColor(propValue)
	case "font-size":
		parsedStyleSheet.FontSize = mapSizeValue(propValue)
	case "display":
		parsedStyleSheet.Display = propValue
	case "postion":
		parsedStyleSheet.Position = propValue
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
	elementStylesheet := &structs.Stylesheet{
		Color:    &structs.ColorRGBA{0, 0, 0, 0},
		FontSize: 0,
		Display:  "",
		Position: "Normal",
	}

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

		if elementStylesheet.Height == float64(0) {
			elementStylesheet.Height = elementStylesheet.FontSize
		}
	}

	if elementStylesheet.Display == "" {
		elementStylesheet.Display = getDefaultElementDisplay(elementName)
	}

	return elementStylesheet
}
