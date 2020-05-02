package thumbnailer

/*
#include "vips.h"

VipsSourceGo * vips_source_go_new ( int id )
{
	VipsSourceGo *source_go;

	VIPS_DEBUG_MSG( "vips_source_go_new:\n" );

	source_go = VIPS_SOURCE_GO( g_object_new( VIPS_TYPE_SOURCE_GO, NULL ) );
	source_go->id = id;

	if( vips_object_build( VIPS_OBJECT( source_go ) ) ) {
		VIPS_UNREF( source_go );
		return( NULL );
	}

	return( source_go );
}
*/
import "C"
import (
	"io"
	"sync"
)

var (
	sourceCtr int
	sources   = make(map[int]*Source)
	sourceMu  = sync.RWMutex{}
)

type Source struct {
	reader  io.Reader
	seeker  io.Seeker
	vipsObj *C.struct__VipsSourceGo
}

func NewSource(image io.Reader) (*Source, error) {
	src := &Source{
		reader: image,
	}

	skr, ok := image.(io.ReadSeeker)
	if ok {
		src.seeker = skr
	}

	sourceMu.Lock()
	sourceCtr++
	id := sourceCtr
	sources[id] = src
	sourceMu.Unlock()

	src.vipsObj = C.vips_source_go_new(C.int(id))

	return src, nil
}
