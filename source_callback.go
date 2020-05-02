package thumbnailer

import "C"
import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

// TODO: Figure out how to pass and return int64 :')
//export GoSourceRead
func GoSourceRead(imageID C.int, buffer unsafe.Pointer, bufSize C.int) (read C.int) {
	sourceMu.RLock()
	src := sources[int(imageID)]
	sourceMu.RUnlock()

	// https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang
	sh := &reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufSize),
		Cap:  int(bufSize),
	}

	buf := *(*[]byte)(unsafe.Pointer(sh))

	n, err := src.reader.Read(buf)
	if errors.Is(err, io.EOF) {
		fmt.Printf("GoSourceRead: EOF [read %d]\n", n)
		return C.int(n)
	} else if err != nil {
		fmt.Printf("GoSourceRead: Error: %v [read %d]\n", err, n)
		return C.int(n)
	}

	fmt.Printf("GoSourceRead: OK [read %d]\n", n)
	return C.int(n)
}

// TODO: Figure out how to return int64 :')
//export GoSourceSeek
func GoSourceSeek(imageID C.int, offset int, whence int) (newOffset C.int) {
	sourceMu.RLock()
	src := sources[int(imageID)]
	sourceMu.RUnlock()

	if src.seeker == nil {
		// Unsupported!
		fmt.Printf("GoSourceSeek: Not supported\n")
		return -1
	}

	n, err := src.seeker.Seek(int64(offset), whence)
	if err != nil {
		fmt.Printf("GoSourceSeek: Error: %v [offset %d]\n", err, n)
		return C.int(n) // FIXME: Perhaps return -1 here?
	}

	fmt.Printf("GoSourceSeek: OK [read %d]\n", n)

	return C.int(n)
}
