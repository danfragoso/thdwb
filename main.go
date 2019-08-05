package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danfragoso/thdwb/ketchup"
	"github.com/danfragoso/thdwb/mustard"
	"github.com/danfragoso/thdwb/sauce"
)

func main() {
	url := os.Args[1]

	resource := sauce.GetResource(url)
	bodyString := string(resource.Body)
	TreeDOM := ketchup.ParseHTML(bodyString)

	js, _ := json.MarshalIndent(TreeDOM.Children, "", " ")
	fmt.Println(string(js))

	mustard.RenderDOM(TreeDOM)
}
