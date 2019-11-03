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

func hexToFloatInRange(hex string) float64 {
	number, err := strconv.ParseInt(hex, 16, 0)

	if err != nil {
		panic(err)
	}

	return float64(number / 255)
}

// RGBAToColor - Transforms RGBA color string to *structs.ColorRGBA
// TODO - Fix this spaghetti and parse alpha values
func RGBAToColor(colorString string) *structs.ColorRGBA {
	var color *structs.ColorRGBA

	if rgbaParams.MatchString(colorString) {
		paramString := rgbaParams.FindString(colorString)
		paramString = strings.Trim(paramString, "()")

		params := strings.Split(paramString, ",")
		paramsLen := len(params)

		if paramsLen >= 3 {
			var red float64
			var green float64
			var blue float64
			var alpha float64

			if strings.HasSuffix(params[0], "%") {
				value, _ := strconv.ParseInt(strings.Trim(strings.TrimSpace(params[0]), "%"), 10, 0)
				red = float64(value / 100)
			} else if strings.Index(params[0], ".") != -1 {
				value, _ := strconv.ParseFloat(strings.TrimSpace(params[0]), 64)
				red = value
			} else {
				value, _ := strconv.Atoi(strings.TrimSpace(params[0]))
				red = float64(value / 255)
			}

			if strings.HasSuffix(params[1], "%") {
				value, _ := strconv.ParseInt(strings.Trim(strings.TrimSpace(params[1]), "%"), 10, 0)
				green = float64(value / 100)
			} else if strings.Index(params[1], ".") != -1 {
				value, _ := strconv.ParseFloat(strings.TrimSpace(params[1]), 64)
				green = value
			} else {
				value, _ := strconv.Atoi(strings.TrimSpace(params[1]))
				green = float64(value / 255)
			}

			if strings.HasSuffix(params[2], "%") {
				value, _ := strconv.ParseInt(strings.Trim(strings.TrimSpace(params[2]), "%"), 10, 0)
				blue = float64(value / 100)
			} else if strings.Index(params[2], ".") != -1 {
				value, _ := strconv.ParseFloat(strings.TrimSpace(params[2]), 64)
				blue = value
			} else {
				value, _ := strconv.Atoi(strings.TrimSpace(params[2]))
				blue = float64(value / 255)
			}

			alpha = 1

			return &structs.ColorRGBA{
				R: red,
				G: green,
				B: blue,
				A: alpha,
			}
		}
	}

	return color
}

// HexStringToColor - Transforms hex color string to *structs.ColorRGBA
func HexStringToColor(colorString string) *structs.ColorRGBA {
	colorString = strings.ToLower(colorString)
	colorStringLen := len(colorString)

	switch colorStringLen {
	case 9:
		return &structs.ColorRGBA{
			R: hexToFloatInRange(colorString[1:3]),
			G: hexToFloatInRange(colorString[3:5]),
			B: hexToFloatInRange(colorString[5:7]),
			A: hexToFloatInRange(colorString[7:9]),
		}

	case 7:
		return &structs.ColorRGBA{
			R: hexToFloatInRange(colorString[1:3]),
			G: hexToFloatInRange(colorString[3:5]),
			B: hexToFloatInRange(colorString[5:7]),
			A: 1,
		}

	case 5:
		return &structs.ColorRGBA{
			R: hexToFloatInRange(colorString[1:2] + colorString[1:2]),
			G: hexToFloatInRange(colorString[2:3] + colorString[2:3]),
			B: hexToFloatInRange(colorString[3:4] + colorString[3:4]),
			A: hexToFloatInRange(colorString[4:5] + colorString[4:5]),
		}

	case 4:
		return &structs.ColorRGBA{
			R: hexToFloatInRange(colorString[1:2] + colorString[1:2]),
			G: hexToFloatInRange(colorString[2:3] + colorString[2:3]),
			B: hexToFloatInRange(colorString[3:4] + colorString[3:4]),
			A: 1,
		}

	default:
		return &structs.ColorRGBA{}
	}
}

// MapCSSColor - Transforms css color strings to #structs.ColorRGBA
func MapCSSColor(colorString string) *structs.ColorRGBA {
	if string(colorString[0]) == "#" {
		return HexStringToColor(colorString)
	} else if rgba.MatchString(colorString) {
		return RGBAToColor(colorString)
	}

	return colorTable[colorString]
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
