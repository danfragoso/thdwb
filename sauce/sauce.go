package sauce

import "io/ioutil"
import "net/http"

type Resource struct {
	Body string
	Code int
}

func GetResource(url string) *Resource {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	
	return &Resource{
		Body: string(body),
	}
}