package thumbnailer

/*
#cgo pkg-config: vips
#include "vips.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"sync"
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

func vipsImageLoad(imageSource *Source) *C.VipsImage {
	return C.vips_new_image_bridge(imageSource.src)
}

func vipsThumbnail(imageSource *Source, width, height int, noAutoRotate, crop bool) (*C.VipsImage, error) {
	var image *C.VipsImage

	noRotate := C.int(boolToInt(noAutoRotate))
	cropParam := C.int(boolToInt(crop))

	err := C.vips_thumbnail_bridge(imageSource.src, &image,
		C.int(width), C.int(height), noRotate, cropParam,
	)
	if int(err) != 0 {
		return nil, vipsError()
	}

	return image, nil
}

func vipsSave(image *C.VipsImage, target *C.VipsTargetCustom, options SaveOptions) error {
	saveErr := C.int(0)
	interlace := C.int(boolToInt(options.Interlace))
	quality := C.int(options.Quality)
	strip := C.int(boolToInt(options.StripMetadata))

	saveErr = C.vips_jpegsave_bridge(image, target, strip, quality, interlace)

	if int(saveErr) != 0 {
		return vipsError()
	}

	C.vips_error_clear()

	return nil
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
