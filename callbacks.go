package thumbnailer

import "C"
import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

//export GoSourceRead
// TODO: Figure out how to pass and return int64 :')
func GoSourceRead(imageID int, buffer unsafe.Pointer, bufSize C.int) (read C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		fmt.Printf("GoSourceRead: Source [id %d] not found\n", imageID)
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

	fmt.Printf("GoSourceRead: OK [read %d]\n", n)
	return C.int(n)
}

//export GoSourceSeek
// TODO: Figure out how to return int64 :')
func GoSourceSeek(imageID int, offset int, whence int) (newOffset C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		fmt.Printf("GoSourceSeek: Source [id %d] not found\n", imageID)
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

//export GoTargetWrite
// TODO: Figure out how to return int64 :')
func GoTargetWrite(imageID int, buffer unsafe.Pointer, bufSize C.int) (written C.int) {
	targetMu.RLock()
	target, ok := targets[imageID]
	targetMu.RUnlock()
	if !ok {
		fmt.Printf("GoTargetWrite: Target [id %d] not found\n", imageID)
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
		fmt.Printf("GoTargetWrite: Error: %v [wrote %d]\n", err, n)
		return C.int(n)
	}

	fmt.Printf("GoTargetWrite: OK [wrote %d]\n", n)
	return C.int(n)
}

//export GoTargetFinish
func GoTargetFinish(imageID int) {
	targetMu.RLock()
	target, ok := targets[imageID]
	targetMu.RUnlock()
	if !ok {
		fmt.Printf("GoTargetFinish: Target [id %d] not found\n", imageID)
		return
	}

	targetMu.Lock()
	delete(targets, imageID)
	targetMu.Unlock()

	fmt.Printf("GoTargetFinish: Closing [id %d]\n", imageID)

	defer target.Cleanup()

	if !target.CloseWriter {
		return
	}

	closer, ok := target.writer.(io.Closer)
	if ok {
		err := closer.Close()
		if err != nil {
			fmt.Printf("GoTargetFinish: Error closing [id %d]: %v", imageID, err)
		}
	}
}
