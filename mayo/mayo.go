package mayo

import (
	"strconv"
	"strings"

	structs "../structs"
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
	case "background-color":
		parsedStyleSheet.BackgroundColor = MapCSSColor(propValue)
	case "font-size":
		parsedStyleSheet.FontSize = mapSizeValue(propValue)
	case "display":
		parsedStyleSheet.Display = propValue
	case "postion":
		parsedStyleSheet.Position = propValue
	case "height":
		parsedStyleSheet.Height = mapSizeValue(propValue)
	case "width":
		parsedStyleSheet.Width = mapSizeValue(propValue)
	}

	return parsedStyleSheet
}

func parseInlineStylesheet(attributes []*structs.Attribute, elementStylesheet *structs.Stylesheet) *structs.Stylesheet {
	for i := 0; i < len(attributes); i++ {
		attributeName := attributes[i].Name
		if attributeName == "style" {

			styleString := attributes[i].Value
			styleProps := strings.Split(strings.Replace(styleString, " ", "", -1), ";")

			for i := 0; i < len(styleProps); i++ {
				styledProperty := strings.Split(styleProps[i], ":")
				if len(styledProperty) >= 2 {
					elementStylesheet = mapPropToStylesheet(elementStylesheet, styledProperty)
				}
			}
		}
	}

	return elementStylesheet
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
		Color:           &structs.ColorRGBA{0, 0, 0, 1},
		BackgroundColor: &structs.ColorRGBA{1, 1, 1, 1},
		FontSize:        0,
		Display:         "",
		Position:        "Normal",
	}

	if hasInlineStyle(attributes) {
		elementStylesheet = parseInlineStylesheet(attributes, elementStylesheet)
	}

	if elementStylesheet.FontSize == float64(0) {
		fontSize := elementFontTable[elementName]

		if fontSize != float64(0) {
			elementStylesheet.FontSize = fontSize
		} else {
			elementStylesheet.FontSize = float64(14)
		}
	}

	if elementStylesheet.Display == "" {
		elementStylesheet.Display = getDefaultElementDisplay(elementName)
	}

	return elementStylesheet
}
