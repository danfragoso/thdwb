package main

import "os"
import "fmt"
import "encoding/json"

import "github.com/danfragoso/thdwb/mustard"
import "github.com/danfragoso/thdwb/ketchup"
import "github.com/danfragoso/thdwb/sauce"

func main() {
	url := os.Args[1]
	
	resource := sauce.GetResource(url)
	bodyString := string(resource.Body)
	
	DOM_Tree := ketchup.ParseHTML(bodyString)
	js, err := json.MarshalIndent(DOM_Tree.Children, "", " ")
	fmt.Println(err)
	fmt.Println(string(js))
	mustard.RenderDOM(DOM_Tree)
}
