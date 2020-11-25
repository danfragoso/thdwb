// +build bundleAssets

package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
			fileString := "package assets\n"
			childFiles, _ := ioutil.ReadDir(rootPath + rootFile.Name())

			for _, childFile := range childFiles {
				if !childFile.IsDir() {
					fileContent, _ := ioutil.ReadFile(rootPath + rootFile.Name() + "/" + childFile.Name())
					fileString += createFileFunction(childFile.Name(), fileContent)
				}
			}

			ioutil.WriteFile(rootPath+rootFile.Name()+".go", []byte(fileString), 0644)
		}
	}
}

func createFileFunction(fileName string, fileContent []byte) string {
	functionTitle := strings.Title(strings.Split(fileName, ".")[0])
	functionStr := fmt.Sprint("\nfunc ", functionTitle, "() []byte {\n")
	functionStr += "\treturn []byte{"

	suffix := ", "
	for byteIdx, _byte := range fileContent {
		if byteIdx == len(fileContent)-1 {
			suffix = "}"
		}

		functionStr += fmt.Sprint(_byte, suffix)
	}

	functionStr += "\n}\n"
	return functionStr
}
