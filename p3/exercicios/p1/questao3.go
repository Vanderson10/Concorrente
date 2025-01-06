package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	path := os.Args[1]
	chFilter := make(chan string)

	go findWrap(chFilter, path)
	selectFiles(chFilter)

}

func find(chOut chan<- string, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			find(chOut, path+"/"+file.Name())
		} else {
			chOut <- path + "/" + file.Name()
		}
	}
}

func findWrap(chOut chan<- string, path string) {
	find(chOut, path)
	close(chOut)
}

func selectFiles(chIn <-chan string) {
	for fileName := range chIn {
		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		if bytes[0]%2 == 0 {
			pathName := strings.Split(fileName, "/")
			fmt.Println(pathName[len(pathName)-1])
		}
	}
}
