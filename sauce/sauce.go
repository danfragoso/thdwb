package sauce

import (
	"io/ioutil"
	"net/http"

	structs "../structs"
)

// GetResource - Makes an http request and returns a resource struct
func GetResource(url string) *structs.Resource {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return &structs.Resource{
		Body: string(body),
	}
}
