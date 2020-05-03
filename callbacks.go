package thumbnailer

import "C"
import (
	"errors"
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
		logger.Errorf("goSourceRead[id %d]: Source not found", imageID)
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
		logger.Debugf("goSourceRead[id %d] EOF [read %d]", imageID, n)
		return C.int(n)
	} else if err != nil {
		logger.Errorf("goSourceRead[id %d]: Error: %v [read %d]", imageID, err, n)
		return -1
	}

	logger.Debugf("goSourceRead[id %d]: OK [read %d]", imageID, n)
	return C.int(n)
}

//export goSourceSeek
// TODO: Figure out how to return int64 :')
func goSourceSeek(imageID int, offset int, whence int) (newOffset C.int) {
	sourceMu.RLock()
	src, ok := sources[imageID]
	sourceMu.RUnlock()
	if !ok {
		logger.Errorf("goSourceSeek[id %d]: Source not found", imageID)
		return -1
	}

	if src.seeker == nil {
		// Unsupported!
		logger.Debugf("goSourceSeek[id %d]: Not supported", imageID)
		return -1
	}

	switch whence {
	case io.SeekStart, io.SeekCurrent, io.SeekEnd:
	default:
		logger.Errorf("goSourceSeek[id %d]: Invalid whence value [%d]", imageID, whence)
		return -1
	}

	n, err := src.seeker.Seek(int64(offset), whence)
	if err != nil {
		logger.Errorf("goSourceSeek[id %d]: Error: %v [offset %d | whence %d]", imageID, err, n, whence)
		return -1
	}

	logger.Debugf("goSourceSeek[id %d]: OK [seek %d | whence %d]", imageID, n, whence)

	return C.int(n)
}

//export goTargetWrite
// TODO: Figure out how to return int64 :')
func goTargetWrite(imageID int, buffer unsafe.Pointer, bufSize C.int) (written C.int) {
	targetMu.RLock()
	target, ok := targets[imageID]
	targetMu.RUnlock()
	if !ok {
		logger.Errorf("goTargetWrite[id %d]: Target not found", imageID)
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
		logger.Errorf("goTargetWrite[id %d]: Error: %v [wrote %d]", imageID, err, n)
		return C.int(n)
	}

	logger.Debugf("goTargetWrite[id %d]: OK [wrote %d]", imageID, n)
	return C.int(n)
}

//export goTargetFinish
func goTargetFinish(imageID int) {
	targetMu.RLock()
	target, ok := targets[imageID]
	targetMu.RUnlock()
	if !ok {
		logger.Errorf("goTargetFinish[id %d]: Target not found", imageID)
		return
	}

	target.finish()
}
