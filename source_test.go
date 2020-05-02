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

	target := NewTarget(outFile)

	err = src.Thumbnail(target)
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
