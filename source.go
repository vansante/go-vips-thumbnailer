package thumbnailer

/*
#include <vips/vips.h>

typedef struct _GoSourceArguments {
	int image_id;
} GoSourceArguments;

GoSourceArguments * create_go_source_arguments( int image_id )
{
	GoSourceArguments * source_args;
	source_args = malloc(sizeof(GoSourceArguments));
	source_args->image_id = image_id;

	return source_args;
}

static gint64
go_read ( VipsSourceCustom *source_custom, void *buffer, gint64 length, GoSourceArguments * source_args )
{
    return GoSourceRead(source_args->image_id, buffer, length);
}

static gint64
go_seek ( VipsSourceCustom *source_custom, gint64 offset, int whence, GoSourceArguments * source_args )
{
	return GoSourceSeek(source_args->image_id, offset, whence);
}

VipsSourceCustom *
create_go_custom_source( GoSourceArguments * source_args )
{
	VipsSourceCustom * source_custom = vips_source_custom_new();

	g_signal_connect( source_custom, "read", G_CALLBACK(go_read), source_args );
	g_signal_connect( source_custom, "seek", G_CALLBACK(go_seek), source_args );

	return source_custom;
}
*/
import "C"
import (
	"fmt"
	"io"
	"sync"
	"unsafe"
)

var (
	sourceCtr int
	sources   = make(map[int]*Source)
	sourceMu  = sync.RWMutex{}
)

type Source struct {
	reader io.Reader
	seeker io.Seeker
	src    *C.struct__VipsSourceCustom
	args   *C.struct__GoSourceArguments
}

func NewSource(image io.Reader) *Source {
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

	fmt.Printf("New Source ID: %d\n", id)
	src.args = C.create_go_source_arguments(C.int(id))
	src.src = C.create_go_custom_source(src.args)

	return src
}

func (s *Source) Cleanup() {
	C.free(unsafe.Pointer(s.args))
	//C.free(unsafe.Pointer(s.target))
}

// TODO: Change into *Target later on
func (s *Source) Thumbnail(target *Target) (err error) {
	img, err := vipsThumbnail(s, 50, 50, true, true)
	if err != nil {
		return fmt.Errorf("error generating thumbnail: %w", err)
	}

	err = vipsSave(img, target.target)
	if err != nil {
		return fmt.Errorf("error saving thumbnail: %w", err)
	}
	return nil
}
