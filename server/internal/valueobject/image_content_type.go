package valueobject

import "errors"

type ImageContentType struct {
	LiteralBase[string]
}

func NewImageContentType(s string) (*ImageContentType, error) {
	switch s {
	case "image/jpeg", "image/png", "image/webp":
		return &ImageContentType{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, errors.New("invalid image content type")
	}

}

func (c ImageContentType) Extension() string {
	switch c.Value() {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	}
	// safe fallback
	return ""
}
