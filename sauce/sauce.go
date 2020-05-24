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
	resource.URL = resp.Request.URL
	resource.Body = string(body)
	return resource
}

func ParseURL(link string) *url.URL {
	URL, err := url.Parse(link)
	if err != nil {
		URL = ParseURL("thdwb://error?err=failedToParseURL")
	}

	return URL
}

func GetImage(URL *url.URL) (string, []byte) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "THDWB (The HotDog Web Browser);")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	img, err := ioutil.ReadAll(resp.Body)

	return "a", img
}
