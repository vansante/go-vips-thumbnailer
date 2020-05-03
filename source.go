package thumbnailer

/*
#include "vips.h"

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

G_DEFINE_TYPE( VipsSourceGo, vips_source_go, VIPS_TYPE_SOURCE );

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
	object_class->description = _( "Go source" );

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
