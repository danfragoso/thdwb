package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danfragoso/thdwb/ketchup"
	//"github.com/danfragoso/thdwb/mustard"
	"github.com/danfragoso/thdwb/sauce"
)

func main() {
	url := os.Args[1]

	resource := sauce.GetResource(url)
	bodyString := string(resource.Body)

	DOM_Tree := ketchup.ParseHTML(bodyString)
	js, err := json.MarshalIndent(DOM_Tree.Children, "", " ")
	fmt.Println(err)
	fmt.Println(string(js))
	//mustard.RenderDOM(DOM_Tree)
}
