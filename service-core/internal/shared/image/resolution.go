package image

import (
	"bytes"
	"fmt"
	"io"
)

type resolutionGenerator struct {
	processor ImageTransformer
}

func NewResolutionGenerator(processor ImageTransformer) VariantCreator {
	return &resolutionGenerator{
		processor: processor,
	}
}

func (rG *resolutionGenerator) GenerateVariant(
	r io.Reader,
	contentType string,
	targetWidth int,
) (io.Reader, int64, error) {
	// Read the input into memory because image decoding consumes the reader.
	// Callers generating multiple variants must provide a fresh reader for each decode.

	img, err := rG.processor.Decode(r, contentType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode image for variant: %w", err)
	}

	resizedImg := rG.processor.Resize(img, targetWidth)

	// Since we fallback webp to jpeg in processor, we should make sure contentType is jpeg if we did
	if contentType == "image/webp" {
		contentType = "image/jpeg"
	}

	encodedBytes, err := rG.processor.Encode(resizedImg, contentType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to encode resized image: %w", err)
	}

	return bytes.NewReader(encodedBytes), int64(len(encodedBytes)), nil
}
