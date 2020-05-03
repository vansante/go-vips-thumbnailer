package thumbnailer

import "C"
import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

// TODO: Figure out how to pass and return int64 :')
//export GoSourceRead
func GoSourceRead(imageID int, buffer unsafe.Pointer, bufSize C.int) (read C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		fmt.Printf("GoSourceRead: Image [id %d] not found \n", imageID)
		return -1
	}

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
		return -1
	}

	if bufSize > 32 {
		fmt.Println(hex.Dump(buf[0:32]))
	} else {
		fmt.Println(hex.Dump(buf))
	}

	fmt.Printf("GoSourceRead: OK [read %d]\n", n)
	return C.int(n)
}

// TODO: Figure out how to return int64 :')
//export GoSourceSeek
func GoSourceSeek(imageID int, offset int, whence int) (newOffset C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		fmt.Printf("GoSourceSeek: Image [id %d] not found \n", imageID)
		return -1
	}

	if src.seeker == nil {
		// Unsupported!
		fmt.Printf("GoSourceSeek: Not supported\n")
		return -1
	}

	switch whence {
	case io.SeekStart, io.SeekCurrent, io.SeekEnd:
	default:
		fmt.Printf("GoSourceSeek: Invalid whence value [%d]\n", whence)
		return -1
	}

	n, err := src.seeker.Seek(int64(offset), whence)
	if err != nil {
		fmt.Printf("GoSourceSeek: Error: %v [offset %d | whence %d]\n", err, n, whence)
		return -1
	}

	fmt.Printf("GoSourceSeek: OK [seek %d | whence %d]\n", n, whence)

	return C.int(n)
}
