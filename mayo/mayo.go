package mayo

import (
	"strconv"
	"strings"

	hotdog "github.com/danfragoso/thdwb/hotdog"
)

func getDefaultElementDisplay(element string) string {
	displayType := "block"

	switch element {
	case "script", "style", "meta", "link", "head", "title":
		displayType = "none"
	case "li":
		displayType = "list-item"
	case "html:text", "a", "abbr", "acronym", "b", "bdo", "big", "br",
		"button", "cite", "code", "dfn", "em", "i", "img", "input", "kbd",
		"label", "map", "object", "output", "q", "samp", "select", "small",
		"span", "strong", "sub", "sup", "textarea", "time", "tt", "var", "font":
		displayType = "inline"
	}

	return displayType
}

func getDefaultElementFontWeight(element string) int {
	switch element {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return 600
	}

	return 400
}

func mapSizeValue(sizeValue string) float64 {
	valueString := sizeValue[0 : len(sizeValue)-2]
	value, err := strconv.ParseInt(valueString, 10, 0)

	if err != nil {
		return float64(14)
	}

	return float64(value)
}

func mapPropToStylesheet(parsedStyleSheet *hotdog.Stylesheet, propSlice []string) *hotdog.Stylesheet {
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

func parseInlineStylesheet(attributes []*hotdog.Attribute, elementStylesheet *hotdog.Stylesheet) *hotdog.Stylesheet {
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

func hasInlineStyle(attributes []*hotdog.Attribute) bool {
	inlineStyle := false

	for i := 0; i < len(attributes); i++ {
		attributeName := attributes[i].Name
		if attributeName == "style" {
			inlineStyle = true
		}
	}

	return inlineStyle
}

func GetElementStylesheet(elementName string, attributes []*hotdog.Attribute) *hotdog.Stylesheet {
	elementStylesheet := &hotdog.Stylesheet{
		BackgroundColor: &hotdog.ColorRGBA{1, 1, 1, 0},
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

	if elementStylesheet.Color == nil {
		color := elementColorTable[elementName]
		if color != nil {
			elementStylesheet.Color = color
		} else {
			elementStylesheet.Color = &hotdog.ColorRGBA{0, 0, 0, 1}
		}
	}

	if elementStylesheet.FontWeight == 0 {
		elementStylesheet.FontWeight = getDefaultElementFontWeight(elementName)
	}

	if elementStylesheet.Display == "" {
		elementStylesheet.Display = getDefaultElementDisplay(elementName)
	}

	return elementStylesheet
}
