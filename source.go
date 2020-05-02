package thumbnailer

/*
#include <vips.h>
#include "source.h"
*/
import "C"
import (
	"io"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var (
	imageID = int32(0)
)

type Source struct {
	id  int32
	rdr io.ReadSeeker
	src *C.struct__VipsSourceCustom
}

func NewSource(image io.ReadSeeker) (*Source, error) {
	vipsSrc := C.vips_source_custom_new()

	id := atomic.AddInt32(&imageID, 1)

	src := &Source{
		id:  id,
		rdr: image,
		src: vipsSrc,
	}

	return src, nil
}

// TODO: Figure out how to pass and return int64 :')
//export GoSourceRead
func GoSourceRead(buffer unsafe.Pointer, bufSize C.int) (read C.int) {
	// https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang
	sh := &reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufSize),
		Cap:  int(bufSize),
	}

	buf := *(*[]byte)(unsafe.Pointer(sh))

	_ = buf[0] // Do rdr.Read()

	return 0
}

// TODO: Figure out how to return int64 :')
//export GoSourceSeek
func GoSourceSeek(offset int, whence int) (newOffset C.int) {
	return 0
}
