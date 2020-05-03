package main

import (
	"gopkg.in/vansante/go-vips-thumbnailer.v1"
	"io/ioutil"
	"os"
)

func main() {
	testFile, err := os.Open("assets/test.jpg")
	if err != nil {
		panic(err)
	}

	src, err := thumbnailer.NewSource(testFile)
	if err != nil {
		panic(err)
	}

	data, err := src.Thumbnail()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("assets/test_thumbnail.jpg", data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
