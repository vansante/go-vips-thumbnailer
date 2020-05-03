package thumbnailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	testFile, err := os.Open("assets/test.jpg")
	assert.NoError(t, err)
	defer testFile.Close()

	src := NewSource(testFile)
	defer src.Cleanup()

	outFile, err := os.Create("assets/output.jpg")
	assert.NoError(t, err)
	defer outFile.Close()

	target := NewTarget(outFile)

	err = src.Thumbnail(target)
	assert.NoError(t, err)

	//spew.Dump(target.target)
}
