package mayo

import (
	"fmt"
	"strings"
)

type Stylesheet struct {
	Color    string
	FontSize int
}

func ParseInlineStyle(styleString string) *Stylesheet {
	styleProps := strings.Split(strings.Replace(styleString, " ", "", -1), ";")

	for i := 0; i < len(styleProps); i++ {
		styledProperty := strings.Split(styleProps[i], ":")
		if len(styledProperty) >= 2 {
			fmt.Println(styledProperty)
		}
	}

	parsedStylesheet := &Stylesheet{}

	return parsedStylesheet
}
