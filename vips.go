package thumbnailer

/*
#cgo pkg-config: vips
#include "vips.h"

int vips_thumbnail_bridge(VipsSourceCustom *source, VipsImage **out, int width, int height, int no_rotate, int crop) {
	if (crop) {
		return vips_thumbnail_source(source, out, width,
			"height", height,
			"no_rotate", INT_TO_GBOOLEAN(no_rotate),
			"crop", VIPS_INTERESTING_CENTRE,
			NULL
		);
	}
	return vips_thumbnail_source(source, out, width,
		"height", height,
		"no_rotate", INT_TO_GBOOLEAN(no_rotate),
		NULL
	);
}

// TODO: Change to a VipsTargetCustom later on
int vips_jpegsave_bridge(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace) {
	return vips_jpegsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"optimize_coding", TRUE,
		"interlace", INT_TO_GBOOLEAN(interlace),
		NULL
	);
}
*/
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

var (
	vipsMu          sync.Mutex
	vipsInitialized bool
)

func init() {
	InitVips()
}

// Initialize is used to explicitly start libvips in thread-safe way.
// Only call this function if you have previously turned off libvips.
func InitVips() {
	if C.VIPS_MAJOR_VERSION <= 8 && C.VIPS_MINOR_VERSION < 9 {
		panic("Unsupported libvips version. Please use version 8.9.0 or higher")
	}

	vipsMu.Lock()
	defer vipsMu.Unlock()

	if vipsInitialized {
		return
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C.vips_init(C.CString("thumbnailer"))
	if err != 0 {
		panic(fmt.Sprintf("unable to start vips [error code: %d]", err))
	}

	// Define a custom thread concurrency limit in libvips (this may generate thread-unsafe issues)
	// See: https://github.com/jcupitt/libvips/issues/261#issuecomment-92850414
	if os.Getenv("VIPS_CONCURRENCY") == "" {
		C.vips_concurrency_set(1)
	}

	vipsInitialized = true
}

// Shutdown is used to shutdown libvips in a thread-safe way.
// You can call this to drop caches as well.
// If libvips was already initialized, the function is no-op
func ShutdownVips() {
	vipsMu.Lock()
	defer vipsMu.Unlock()

	if vipsInitialized {
		C.vips_shutdown()
		vipsInitialized = false
	}
}

// TODO: Change into writing to a VipsTargetGo
func vipsThumbnail(imageSource *Source, width, height int, autoRotate, crop bool) (*C.VipsImage, error) {
	var image *C.VipsImage

	noRotate := C.int(boolToInt(!autoRotate))
	cropParam := C.int(boolToInt(crop))

	err := C.vips_thumbnail_bridge(imageSource.src, &image,
		C.int(width), C.int(height), noRotate, cropParam,
	)
	if int(err) != 0 {
		return nil, vipsError()
	}

	return image, nil
}

func vipsSave(image *C.VipsImage) ([]byte, error) {
	length := C.size_t(0)
	saveErr := C.int(0)
	interlace := C.int(0)
	quality := C.int(70)
	strip := C.int(1)

	var ptr unsafe.Pointer
	saveErr = C.vips_jpegsave_bridge(image, &ptr, &length, strip, quality, interlace)

	if int(saveErr) != 0 {
		return nil, vipsError()
	}

	buf := C.GoBytes(ptr, C.int(length))

	// Clean up
	C.g_free(C.gpointer(ptr))
	C.vips_error_clear()

	return buf, nil
}

func vipsError() error {
	errStr := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return fmt.Errorf("vips error: %v", errStr)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
