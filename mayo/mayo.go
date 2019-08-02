package mayo

import (
	"fmt"
	"strings"

	"github.com/danfragoso/thdwb/structs"
)

func ParseInlineStylesheet(attributes []*structs.Attribute) *structs.Stylesheet {
	var parsedStylesheet *structs.Stylesheet

	for i := 0; i < len(attributes); i++ {
		attributeName := attributes[i].Name
		if attributeName == "style" {

			styleString := attributes[i].Value
			styleProps := strings.Split(strings.Replace(styleString, " ", "", -1), ";")

			for i := 0; i < len(styleProps); i++ {
				styledProperty := strings.Split(styleProps[i], ":")
				if len(styledProperty) >= 2 {
					fmt.Println(styledProperty)
				}
			}
		}
	}

	return parsedStylesheet
}
