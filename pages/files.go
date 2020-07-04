package pages

import (
	"bytes"
	"html/template"
	"io/ioutil"
)

type FileBrowser struct {
	Path  string
	Dirs  []string
	Files []string
}

func RenderFileBrowser(path string) string {
	renderedTemplate := new(bytes.Buffer)
	fileBrowser := &FileBrowser{Path: path}

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			fileBrowser.Dirs = append(fileBrowser.Dirs, file.Name())
		} else {
			fileBrowser.Files = append(fileBrowser.Files, file.Name())
		}
	}

	t, _ := template.New("fileBrowser").Parse(fileBrowserTemplate())
	t.ExecuteTemplate(renderedTemplate, "fileBrowser", fileBrowser)

	return renderedTemplate.String()
}
