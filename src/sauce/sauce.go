package sauce

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"thdwb/assets"
	pages "thdwb/pages"
	hotdog "thdwb/hotdog"
)

var client = &http.Client{}
var cache = &hotdog.ResourceCache{}
var imageCache = &hotdog.ImgCache{}

// GetResource - Makes an http request and returns a resource struct
func GetResource(URL *url.URL, browser *hotdog.WebBrowser) *hotdog.Resource {
	switch URL.Scheme {
	case "thdwb":
		return fetchInternalPage(URL, browser)
	case "file":
		return &hotdog.Resource{Body: pages.RenderFileBrowser(URL.Path), URL: URL}
	case "":
		URL.Scheme = "http"
		break
	}

	return fetchExternalPage(URL)
}

func fetchInternalPage(URL *url.URL, browser *hotdog.WebBrowser) *hotdog.Resource {
	switch URL.Host {
	case "homepage":
		return &hotdog.Resource{
			Body: string(assets.HomePage()),
			URL:  URL,
		}

	case "history":
		return &hotdog.Resource{
			Body: buildHistoryPage(browser.History),
			URL:  URL,
		}
	case "about":
		return &hotdog.Resource{
			Body: pages.RenderAboutPage(browser.BuildInfo),
			URL:  URL,
		}
	default:
		return &hotdog.Resource{
			Body: string(assets.DefaultPage()),
			URL:  URL,
		}
	}
}

func fetchExternalPage(URL *url.URL) *hotdog.Resource {
	url := URL.String()
	go hotdog.Log("sauce", "Downloading page "+url)

	cachedResource := cache.GetResource(url)
	if cachedResource != nil {
		return cachedResource
	}

	resource := &hotdog.Resource{Key: url, URL: URL}
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
	}

	var img []byte
	if len(imgUrl) >= 22 && imgUrl[:22] == "data:image/png;base64," {
		imgData := imgUrl[strings.IndexByte(imgUrl, ',')+1:]

		decodedData, err := base64.RawStdEncoding.DecodeString(imgData)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to decode base64 data")
		}

		img = decodedData
	} else {
		req, err := http.NewRequest("GET", imgUrl, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "THDWB (The HotDog Web Browser);")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to fetch image")
			return nil
		}
		defer resp.Body.Close()

		img, err = ioutil.ReadAll(resp.Body)
	}

	imageCache.AddImage(imgUrl, img)
	return img
}

func buildHistoryPage(history *hotdog.History) string {
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
