package thumbnailer

/*
#include "vips.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

#include <vips/vips.h>
#include <vips/debug.h>

#define VIPS_TYPE_SOURCE_GO (vips_source_get_type())
#define VIPS_SOURCE_GO( obj ) \
(G_TYPE_CHECK_INSTANCE_CAST( (obj), \
VIPS_TYPE_SOURCE_GO, VipsSourceGo ))
#define VIPS_SOURCE_GO_CLASS( klass ) \
(G_TYPE_CHECK_CLASS_CAST( (klass), \
VIPS_TYPE_SOURCE_GO, VipsSourceGoClass))
#define VIPS_IS_SOURCE_GO( obj ) \
(G_TYPE_CHECK_INSTANCE_TYPE( (obj), VIPS_TYPE_SOURCE_GO ))
#define VIPS_IS_SOURCE_GO_CLASS( klass ) \
(G_TYPE_CHECK_CLASS_TYPE( (klass), VIPS_TYPE_SOURCE_GO ))
#define VIPS_SOURCE_GO_GET_CLASS( obj ) \
(G_TYPE_INSTANCE_GET_CLASS( (obj), \
VIPS_TYPE_SOURCE_GO, VipsSourceGoClass ))

static gint64
vips_source_go_read_real ( VipsSource *source, void *buffer, size_t length )
{
    printf( "GO TEST READ :\n" );
    fflush(stdout);
    return 0;
//    return GoSourceRead(1, buffer, length);
}

static gint64
vips_source_go_seek_real ( VipsSource *source, gint64 offset, int whence )
{
    printf( "GO TEST SEEK :\n" );
    fflush(stdout);
    return 0;
}

static gint64
vips_source_go_read_go ( VipsSourceGo *source, void *buffer, gint64 length )
{
    printf( "GO TEST READ 2:\n" );
    fflush(stdout);
    return 0;
//    return GoSourceRead(source->id, buffer, length);
}

static gint64
vips_source_go_seek_go ( VipsSourceGo *source, gint64 offset, int whence )
{
    printf( "GO TEST SEEK 2:\n" );
    fflush(stdout);
    return 0;
//	return GoSourceSeek(source->id, offset, whence);
}

static void
vips_source_go_class_init ( VipsSourceGoClass *class )
{
    printf( "GO CLASSSS :\n" );
    fflush(stdout);

	VipsObjectClass *object_class = VIPS_OBJECT_CLASS( class );
	VipsSourceClass *source_class = VIPS_SOURCE_CLASS( class );

	object_class->nickname = "go source";
	object_class->description = "Go source";

	class->read = vips_source_go_read_go;
    class->seek = vips_source_go_seek_go;

    source_class->read = vips_source_go_read_real;
    source_class->seek = vips_source_go_seek_real;
}

static void
vips_source_go_init( VipsSourceGo *source_go )
{
	printf( "GO CLASSSS initt:\n" );
    fflush(stdout);
}

G_DEFINE_TYPE( VipsSourceGo, vips_source_go, VIPS_TYPE_SOURCE );

VipsSourceGo * vips_source_go_new ( int id )
{
	VipsSourceGo *source_go;

	printf ( "vips_source_go_new: %d \n", id );
	fflush(stdout);

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
	"fmt"
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

	fmt.Println(C.GoStringN(src.vipsObj.parent_object.parent_object.parent_object.description, 100))

	return src, nil
}

// TODO: Change into *Target later on
func (s *Source) Thumbnail() (thumbnail []byte, err error) {
	img, err := vipsThumbnail(s, 50, 50, false, true)
	if err != nil {
		return nil, fmt.Errorf("error generating thumbnail: %w", err)
	}

	thumbnail, err = vipsSave(img)
	if err != nil {
		return nil, fmt.Errorf("error saving thumbnail: %w", err)
	}
	return thumbnail, nil
}
