package thumbnailer

import "C"
import (
	"fmt"
)

type Image struct {
	img *C.struct__VipsImage
}

func (i *Image) Save(target *Target, options SaveOptions) error {
	err := vipsSave(i.img, target.target, options)
	if err != nil {
		return fmt.Errorf("error saving thumbnail: %w", err)
	}
	return nil
}
