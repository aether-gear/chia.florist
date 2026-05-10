package domain

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	MIME           string
	ResolutionType string
)

const (
	MIMEJPEG MIME = "image/jpeg"
	MIMEPNG  MIME = "image/png"
	MIMEWEBP MIME = "image/webp"
)

const (
	DirPrefix   = "products"
	KeyPattern  = "products/{product_id}/{image_id}.{ext}"
	MaxFileSize = 5 << 20
	MaxFileMB   = MaxFileSize / (1 << 20)
)

const (
	ResolutionCatalog ResolutionType = "catalog"
	ResolutionCart    ResolutionType = "cart"
	ResolutionDetail  ResolutionType = "detail"
)

var ResolutionWidths = map[ResolutionType]int{
	ResolutionCatalog: 600,
	ResolutionCart:    150,
	ResolutionDetail:  1200,
}

type ImageVariant struct {
	URL string
	Key string
}

type ProductImage struct {
	ID        uuid.UUID
	ProductID uuid.UUID

	Thumbnail ImageVariant
	Preview   ImageVariant
	Detail    ImageVariant

	IsPrimary    bool
	DisplayOrder int

	CreatedAt time.Time
	DeletedAt *time.Time
}

type ProductImageMetadata struct {
	OriginalName string
	MIMEType     MIME
	SizeBytes    int64
}

func (m MIME) IsAllowed() bool {
	switch m {
	case MIMEJPEG, MIMEPNG, MIMEWEBP:
		return true
	default:
		return false
	}
}

func (m MIME) Extension() (string, error) {
	switch m {
	case MIMEJPEG:
		return "jpg", nil
	case MIMEPNG:
		return "png", nil
	case MIMEWEBP:
		return "webp", nil
	default:
		return "", ErrInvalidProductImageMime
	}
}

func (m ProductImageMetadata) Validate() error {
	if !m.MIMEType.IsAllowed() {
		return ErrInvalidProductImageMime
	}

	if m.SizeBytes <= 0 || m.SizeBytes > MaxFileSize {
		return ErrInvalidProductImageSize
	}

	if !IsValidProductImageFileName(m.OriginalName) {
		return ErrInvalidProductImageFileName
	}

	return nil
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

func BuildProductImageObjectKey(
	productID, imageID uuid.UUID,
	mimeType MIME,
) (string, error) {
	ext, err := mimeType.Extension()
	if err != nil {
		return "", err
	}

	str := path.Join(
		DirPrefix,
		productID.String(),
		fmt.Sprintf("%s.%s", imageID.String(), ext),
	)
	return str, nil
}
