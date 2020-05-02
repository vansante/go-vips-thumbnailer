package thumbnailer

/*
#cgo pkg-config: vips
#include "source.h"
*/
import "C"

import (
	"fmt"
	"os"
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
