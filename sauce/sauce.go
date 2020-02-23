package sauce

import (
	"io/ioutil"
	"log"
	"net/http"

	structs "../structs"
)

// GetResource - Makes an http request and returns a resource struct
func GetResource(url string) *structs.Resource {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "THDWB (The HotDog Web Browser);")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return &structs.Resource{
		Body: string(body),
	}
}
