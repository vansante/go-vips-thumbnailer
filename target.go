package thumbnailer

/*
#include <vips/vips.h>

typedef struct _GoTargetArguments {
	int image_id;
} GoTargetArguments;

GoTargetArguments * create_go_target_arguments( int image_id )
{
	GoTargetArguments * target_args;
	target_args = malloc(sizeof(GoTargetArguments));
	target_args->image_id = image_id;

	return target_args;
}

static gint64
go_write ( VipsTargetCustom *target_custom, const void *data, gint64 length, GoTargetArguments * target_args )
{
	return goTargetWrite ( target_args->image_id );
}

static void
go_finish ( VipsTargetCustom *target_custom, GoTargetArguments * target_args )
{
	goTargetFinish ( target_args->image_id );
}

VipsTargetCustom *
create_go_custom_target( GoTargetArguments * target_args )
{
	VipsTargetCustom * target_custom = vips_target_custom_new();

	g_signal_connect( target_custom, "write", G_CALLBACK(go_write), target_args );
	g_signal_connect( target_custom, "finish", G_CALLBACK(go_finish), target_args );

	return target_custom;
}
*/
import "C"
import (
	"io"
	"sync"
	"unsafe"
)

var (
	targetCtr int
	targets   = make(map[int]*Target)
	targetMu  = sync.RWMutex{}
)

type Target struct {
	id          int
	writer      io.Writer
	target      *C.struct__VipsTargetCustom
	args        *C.struct__GoTargetArguments
	CloseWriter bool
}

func NewTarget(writer io.Writer) *Target {
	target := &Target{
		writer: writer,
	}

	targetMu.Lock()
	targetCtr++
	target.id = targetCtr
	targets[target.id] = target
	targetMu.Unlock()

	target.args = C.create_go_target_arguments(C.int(target.id))
	target.target = C.create_go_custom_target(target.args)

	return target
}

func (t *Target) finish() {
	targetMu.Lock()
	delete(targets, t.id)
	targetMu.Unlock()

	logger.Debugf("goTargetFinish[id %d]: Closing [closeWriter: %v]", t.id, t.CloseWriter)

	t.free()

	if !t.CloseWriter {
		return
	}
	closer, ok := t.writer.(io.Closer)
	if ok {
		err := closer.Close()
		if err != nil {
			logger.Errorf("goTargetFinish[id %d]: Error closing: %v", t.id, err)
		}
	}
}

func (t *Target) free() {
	C.free(unsafe.Pointer(t.args))
	//C.free(unsafe.Pointer(s.target))
}
