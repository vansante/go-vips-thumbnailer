package thumbnailer

// ImageType represents an image type value.
type ImageType int

const (
	// UNKNOWN represents an unknow image type value.
	UNKNOWN ImageType = iota
	// JPEG represents the JPEG image type.
	JPEG
	// WEBP represents the WEBP image type.
	WEBP
	// PNG represents the PNG image type.
	PNG
	// TIFF represents the TIFF image type.
	TIFF
	// GIF represents the GIF image type.
	GIF
	// PDF represents the PDF type.
	PDF
	// SVG represents the SVG image type.
	SVG
	// MAGICK represents the libmagick compatible genetic image type.
	MAGICK
	// HEIF represents the HEIC/HEIF/HVEC image type
	HEIF
)

// ImageTypes stores as pairs of image types supported and its alias names.
var ImageTypes = map[ImageType]string{
	JPEG:   "jpeg",
	PNG:    "png",
	WEBP:   "webp",
	TIFF:   "tiff",
	GIF:    "gif",
	PDF:    "pdf",
	SVG:    "svg",
	MAGICK: "magick",
	HEIF:   "heif",
}

type SaveOptions struct {
	ImageType     ImageType
	StripMetadata bool
	Interlace     bool
	Lossless      bool
	Quality       int
	Compression   int
	//NoProfile     bool
	//Interlace     bool
}

type ThumbnailOptions struct {
	Width             int
	Height            int
	Crop              bool
	DisableAutoRotate bool
}
