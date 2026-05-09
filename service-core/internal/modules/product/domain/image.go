package domain

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ProductImageDirPrefix       = "products"
	ProductImageKeyPattern      = "products/{product_id}/{image_id}.{ext}"
	ProductImageMaxUploadSize   = 5 << 20
	ProductImageMaxUploadSizeMB = ProductImageMaxUploadSize / (1 << 20)
)

var (
	ErrInvalidProductImageMime     = errors.New("unsupported product image mime type")
	ErrInvalidProductImageFileName = errors.New("malformed product image filename")
	ErrInvalidProductImageSize     = errors.New("product image exceeds max upload size")
)

type ProductImageMIME string

const (
	ProductImageJPEG ProductImageMIME = "image/jpeg"
	ProductImagePNG  ProductImageMIME = "image/png"
	ProductImageWEBP ProductImageMIME = "image/webp"
)

type ProductImage struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	ObjectKey  string
	URL        string
	IsMain     bool
	Metadata   ProductImageMetadata
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	ArchivedAt *time.Time
}

type ProductImageMetadata struct {
	OriginalName string
	MIMEType     ProductImageMIME
	SizeBytes    int64
}

func (m ProductImageMIME) IsAllowed() bool {
	switch m {
	case ProductImageJPEG, ProductImagePNG, ProductImageWEBP:
		return true
	default:
		return false
	}
}

func (m ProductImageMIME) Extension() (string, error) {
	switch m {
	case ProductImageJPEG:
		return "jpg", nil
	case ProductImagePNG:
		return "png", nil
	case ProductImageWEBP:
		return "webp", nil
	default:
		return "", ErrInvalidProductImageMime
	}
}

func (m ProductImageMetadata) Validate() error {
	if !m.MIMEType.IsAllowed() {
		return ErrInvalidProductImageMime
	}

	if m.SizeBytes <= 0 || m.SizeBytes > ProductImageMaxUploadSize {
		return ErrInvalidProductImageSize
	}

	if !IsValidProductImageFileName(m.OriginalName) {
		return ErrInvalidProductImageFileName
	}

	return nil
}

func BuildProductImageObjectKey(productID, imageID uuid.UUID, mimeType ProductImageMIME) (string, error) {
	ext, err := mimeType.Extension()
	if err != nil {
		return "", err
	}

	str := path.Join(ProductImageDirPrefix, productID.String(), fmt.Sprintf("%s.%s", imageID.String(), ext))
	return str, nil
}

func IsValidProductImageFileName(name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return false
	}

	return path.Base(name) == name
}
