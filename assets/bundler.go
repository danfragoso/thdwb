// +build bundleAssets

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootPath := "./assets/"
	rootFiles, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, rootFile := range rootFiles {
		if rootFile.IsDir() {
			fileString := "package assets\n\nimport _ \"embed\"\n"
			childFiles, _ := ioutil.ReadDir(rootPath + rootFile.Name())

			for _, childFile := range childFiles {
				if !childFile.IsDir() {
					fileContent, _ := ioutil.ReadFile(rootPath + rootFile.Name() + "/" + childFile.Name())
					fileString += createFileFunction(rootFile.Name(), childFile.Name(), fileContent)
				}
			}

			ioutil.WriteFile(rootPath+rootFile.Name()+".go", []byte(fileString), 0644)
		}
	}
}

func createFileFunction(folderName string, fileName string, fileContent []byte) string {
	fileTitle := strings.Split(fileName, ".")[0]

	embedStr := "\n//go:embed " + filepath.Join(folderName, fileName) + "\n"
	embedStr += "var " + fileTitle + " []byte\n"

	functionTitle := strings.Title(fileTitle)
	functionStr := fmt.Sprint("func ", functionTitle, "() []byte {\n")
	functionStr += "\treturn " + fileTitle

	functionStr += "\n}\n"
	return embedStr + functionStr
}
