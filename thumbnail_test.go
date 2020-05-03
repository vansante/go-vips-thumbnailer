package thumbnailer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	SetLogger(&TestLogger{t})

	testFile, err := os.Open("assets/test.jpg")
	assert.NoError(t, err)
	defer testFile.Close()

	src := NewSource(testFile)
	defer src.free()

	outFile, err := os.Create("assets/output.jpg")
	assert.NoError(t, err)
	defer outFile.Close()

	image, err := src.Thumbnail(ThumbnailOptions{
		Width:             120,
		Height:            120,
		Crop:              true,
		DisableAutoRotate: true,
	})
	assert.NoError(t, err)

	target := NewTarget(outFile)
	err = image.Save(target, SaveOptions{
		Interlace:     false,
		StripMetadata: true,
		Lossless:      false,
		Quality:       50,
		Compression:   0,
	})
	assert.NoError(t, err)

	//spew.Dump(target.target)
}

type TestLogger struct {
	t *testing.T
}

func (l TestLogger) Debugf(format string, args ...interface{}) {
	l.t.Logf("[DEBUG] "+format, args...)
}

func (l TestLogger) Errorf(format string, args ...interface{}) {
	l.t.Logf("[ERROR] "+format, args...)
}
