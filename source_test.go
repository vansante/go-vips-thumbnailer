package thumbnailer

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	testFile, err := os.Open("assets/test.jpg")
	assert.NoError(t, err)

	src, err := NewSource(testFile)
	assert.NoError(t, err)

	data, err := src.Thumbnail()
	assert.NoError(t, err)

	err = ioutil.WriteFile("assets/test_thumbnail.jpg", data, os.ModePerm)
	assert.NoError(t, err)

	fmt.Printf("%#v\n\n", src.vipsObj)
}
