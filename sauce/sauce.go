package sauce

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"thdwb/assets"
	structs "thdwb/structs"
)

var client = &http.Client{}
var cache = &structs.ResourceCache{}
var imageCache = &structs.ImgCache{}

// GetResource - Makes an http request and returns a resource struct
func GetResource(URL *url.URL, browser *structs.WebBrowser) *structs.Resource {
	if URL.Scheme == "thdwb" {
		return fetchInternalPage(URL, browser)
	}

	return fetchExternalPage(URL.String())
}

func fetchInternalPage(URL *url.URL, browser *structs.WebBrowser) *structs.Resource {
	switch URL.Host {
	case "homepage":
		return &structs.Resource{
			Body: string(assets.HomePage()),
			URL:  URL,
		}

	case "history":
		return &structs.Resource{
			Body: buildHistoryPage(browser.History),
			URL:  URL,
		}
	default:
		return &structs.Resource{
			Body: string(assets.DefaultPage()),
			URL:  URL,
		}
	}
}

func fetchExternalPage(url string) *structs.Resource {
	cachedResource := cache.GetResource(url)
	if cachedResource != nil {
		return cachedResource
	} else {
		resource := &structs.Resource{Key: url}
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

		cache.AddResource(resource)
		return resource
	}
}

func ParseURL(link string) *url.URL {
	URL, err := url.Parse(link)
	if err != nil {
		URL = ParseURL("thdwb://error?err=failedToParseURL")
	}

	return URL
}

func GetImage(URL *url.URL) []byte {
	imgUrl := URL.String()
	cachedImage := imageCache.GetImage(imgUrl)

	if cachedImage != nil {
		return cachedImage.Image
	} else {
		req, err := http.NewRequest("GET", imgUrl, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "THDWB (The HotDog Web Browser);")

		resp, _ := client.Do(req)

		defer resp.Body.Close()
		img, err := ioutil.ReadAll(resp.Body)

		imageCache.AddImage(imgUrl, img)
		return img
	}
}

func buildHistoryPage(history *structs.History) string {
	d := `
	<html>
		<head>
			<title>History</title>
		</head>
		<body>
		<h1>History</h1>
		<ul>
	`
	for _, page := range history.AllPages() {
		d += `<li><a href="` + page.String() + `">` + page.String() + `</a></li>`
	}

	d += `
		</ul>
		</body>
	</html>
	`
	return d
}
