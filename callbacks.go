package thumbnailer

import "C"
import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

//export goSourceRead
// TODO: Figure out how to pass and return int64 :')
func goSourceRead(imageID int, buffer unsafe.Pointer, bufSize C.int) (read C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		fmt.Printf("goSourceRead: Source [id %d] not found\n", imageID)
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
		fmt.Printf("goSourceRead: EOF [read %d]\n", n)
		return C.int(n)
	} else if err != nil {
		fmt.Printf("goSourceRead: Error: %v [read %d]\n", err, n)
		return -1
	}

	fmt.Printf("goSourceRead: OK [read %d]\n", n)
	return C.int(n)
}

//export goSourceSeek
// TODO: Figure out how to return int64 :')
func goSourceSeek(imageID int, offset int, whence int) (newOffset C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		fmt.Printf("goSourceSeek: Source [id %d] not found\n", imageID)
		return -1
	}

	if src.seeker == nil {
		// Unsupported!
		fmt.Printf("goSourceSeek: Not supported\n")
		return -1
	}

	switch whence {
	case io.SeekStart, io.SeekCurrent, io.SeekEnd:
	default:
		fmt.Printf("goSourceSeek: Invalid whence value [%d]\n", whence)
		return -1
	}

	n, err := src.seeker.Seek(int64(offset), whence)
	if err != nil {
		fmt.Printf("goSourceSeek: Error: %v [offset %d | whence %d]\n", err, n, whence)
		return -1
	}

	fmt.Printf("goSourceSeek: OK [seek %d | whence %d]\n", n, whence)

	return C.int(n)
}

//export goTargetWrite
// TODO: Figure out how to return int64 :')
func goTargetWrite(imageID int, buffer unsafe.Pointer, bufSize C.int) (written C.int) {
	targetMu.RLock()
	target, ok := targets[imageID]
	targetMu.RUnlock()
	if !ok {
		fmt.Printf("goTargetWrite: Target [id %d] not found\n", imageID)
		return -1
	}

	// https://stackoverflow.com/questions/51187973/how-to-create-an-array-or-a-slice-from-an-array-unsafe-pointer-in-golang
	sh := &reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufSize),
		Cap:  int(bufSize),
	}
	buf := *(*[]byte)(unsafe.Pointer(sh))

	n, err := target.writer.Write(buf)
	if err != nil {
		fmt.Printf("goTargetWrite: Error: %v [wrote %d]\n", err, n)
		return C.int(n)
	}

	fmt.Printf("goTargetWrite: OK [wrote %d]\n", n)
	return C.int(n)
}

//export goTargetFinish
func goTargetFinish(imageID int) {
	targetMu.RLock()
	target, ok := targets[imageID]
	targetMu.RUnlock()
	if !ok {
		fmt.Printf("goTargetFinish: Target [id %d] not found\n", imageID)
		return
	}

	target.finish()
}
