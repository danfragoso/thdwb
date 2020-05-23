package sauce

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"thdwb/assets"
	structs "thdwb/structs"
)

// GetResource - Makes an http request and returns a resource struct
func GetResource(URL *url.URL) *structs.Resource {
	if URL.Scheme == "thdwb" {
		if URL.Host == "homepage" {
			return &structs.Resource{
				Body: string(assets.HomePage()),
			}
		}

		return fetchInternalPage(URL.String())
	}

	return fetchExternalPage(URL.String())
}

func fetchInternalPage(url string) *structs.Resource {
	resource := &structs.Resource{}
	resource.Body = `
		<html>
			<head></head>
			<body>
				<div>thdwb</div>
			</body>
		</html>
	`
	return resource
}

func fetchExternalPage(url string) *structs.Resource {
	client := &http.Client{}
	resource := &structs.Resource{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "THDWB (The HotDog Web Browser);")

	resp, err := client.Do(req)
	if err != nil {
		resource.Body = loadErrorPage(err.Error())
		return resource
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	resource.ContentType = resp.Header.Get("Content-Type")
	resource.Body = string(body)
	return resource
}
