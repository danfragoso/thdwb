package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danfragoso/thdwb/ketchup"
	"github.com/danfragoso/thdwb/mustard"
	"github.com/danfragoso/thdwb/sauce"
	"github.com/danfragoso/thdwb/structs"
)

func removeParents(TreeDOM *structs.NodeDOM) {
	nodeChildren := TreeDOM.Children
	TreeDOM.Parent = nil

	for i := 0; i < len(nodeChildren); i++ {
		removeParents(nodeChildren[i])
	}
}

func main() {
	url := os.Args[1]

	resource := sauce.GetResource(url)
	bodyString := string(resource.Body)
	TreeDOM := ketchup.ParseHTML(bodyString)
	mustard.RenderDOM(TreeDOM)

	removeParents(TreeDOM)
	js, err := json.MarshalIndent(TreeDOM.Children, "", " ")
	fmt.Println(err)
	fmt.Println(string(js))

}
