package valueobject

import "errors"

type ImageContentTypeCode string

const (
	ImageContentTypeJPEG ImageContentTypeCode = "image/jpeg"
	ImageContentTypePNG  ImageContentTypeCode = "image/png"
	ImageContentTypeWebP ImageContentTypeCode = "image/webp"
)

var ErrInvalidImageContentType = errors.New("invalid image content type")

type ImageContentType struct {
	LiteralBase[string]
}

func NewImageContentTypeFromString(s string) (*ImageContentType, error) {
	switch ImageContentTypeCode(s) {
	case ImageContentTypeJPEG, ImageContentTypePNG, ImageContentTypeWebP:
		return &ImageContentType{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidImageContentType
	}

}

func (c ImageContentType) Extension() string {
	switch ImageContentTypeCode(c.Value()) {
	case ImageContentTypeJPEG:
		return "jpg"
	case ImageContentTypePNG:
		return "png"
	case ImageContentTypeWebP:
		return "webp"
	default:
		return ""
	}
}
