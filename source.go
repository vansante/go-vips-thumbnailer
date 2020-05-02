package thumbnailer

/*
#include <vips.h>
#include "source.c"
*/
import "C"
import (
	"fmt"
	"io"
)

var (
	signalRead = C.CString("read")
	signalSeek = C.CString("seek")
)

type Source struct {
	rdr io.ReadSeeker
	src *C.struct__VipsSourceCustom
}

func NewSource(image io.ReadSeeker) (*Source, error) {
	vipsSrc := C.vips_source_custom_new()

	src := &Source{
		rdr: image,
		src: vipsSrc,
	}

	gobj := vipsSrc.parent_object.parent_object.parent_object.parent_instance

	fmt.Printf("%#v\n\n", gobj)

	C.g_signal_connect_closure(vipsSrc, signalRead, nil, 0)

	return src, nil
}
