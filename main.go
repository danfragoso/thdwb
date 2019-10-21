package main

import (
	"os"

	"github.com/danfragoso/thdwb/ketchup"
	"github.com/danfragoso/thdwb/mustard"
	"github.com/danfragoso/thdwb/sauce"
)

func main() {
	url := os.Args[1]
	resource := sauce.GetResource(url)
	htmlString := string(resource.Body)
	parsedDocument := ketchup.ParseDocument(htmlString)

	//js, _ := json.MarshalIndent(parsedDocument.Children, "", " ")
	//fmt.Println(string(js))

	mustard.RenderDocument(parsedDocument, url)
}
