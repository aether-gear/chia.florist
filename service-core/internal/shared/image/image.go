package image

import (
	"image"
	"io"
)

type VariantCreator interface {
	GenerateVariant(
		r io.Reader,
		contentType string,
		resolution int,
	) (io.Reader, int64, error)
}

type ImageTransformer interface {
	Resize(img image.Image, targetWidth int) image.Image

	Decode(r io.Reader, contentType string) (image.Image, error)

	Encode(img image.Image, contentType string) ([]byte, error)
}
