package thumbnailer

/*
#include <vips.h>
#include "source.h"
*/
import "C"
import (
	"io"
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

	return src, nil
}
