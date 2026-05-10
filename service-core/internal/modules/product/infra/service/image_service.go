package service

// import (
// 	"fmt"
// 	"io"

// 	infraStorage "service-core/internal/infra/storage"
// 	"service-core/internal/modules/product/domain"
// 	"service-core/internal/modules/product/repository"
// 	sharedImage "service-core/internal/shared/image"

// 	"github.com/google/uuid"
// )

// type uploadService struct {
// 	fileStore infraStorage.Provider
// 	imageGen  sharedImage.VariantCreator
// }

// func NewProductImageUploadRepository(
// 	fileStore infraStorage.Provider,
// 	imageGen sharedImage.VariantCreator,
// ) repository.ProductImageUploadService {
// 	return &uploadService{
// 		fileStore: fileStore,
// 		imageGen:  imageGen,
// 	}
// }

// func (uS *uploadService) UploadCatalogImage(
// 	params repository.UploadProductImageParams,
// ) (string, error) {
// 	return uS.uploadVariant(
// 		params.ProductID,
// 		params.Metadata,
// 		params.File,
// 		domain.ResolutionCatalog,
// 	)
// }

// func (uS *uploadService) UploadCartImage(
// 	params repository.UploadProductImageParams,
// ) (string, error) {
// 	return uS.uploadVariant(
// 		params.ProductID,
// 		params.Metadata,
// 		params.File,
// 		domain.ResolutionCart,
// 	)
// }

// func (uS *uploadService) UploadDetailImage(
// 	params repository.UploadProductImageParams,
// ) (string, error) {
// 	return uS.uploadVariant(
// 		params.ProductID,
// 		params.Metadata,
// 		params.File,
// 		domain.ResolutionDetail,
// 	)
// }

// func (uS *uploadService) DeleteUploadedImages(
// 	urls []string,
// ) error {
// 	var lastErr error
// 	for _, u := range urls {
// 		if err := uS.fileStore.Delete(u); err != nil {
// 			lastErr = fmt.Errorf("delete image %s failed: %w", u, err)
// 		}
// 	}
// 	return lastErr
// }

// func (uS *uploadService) uploadVariant(
// 	productID uuid.UUID,
// 	metadata domain.ProductImageMetadata,
// 	file io.Reader,
// 	resolution domain.ResolutionType,
// ) (string, error) {
// 	// if err := metadata.Validate(); err != nil {
// 	// 	return "", err
// 	// }

// 	targetWidth := domain.ResolutionWidths[resolution]
// 	resizedStream, size, err := uS.imageGen.GenerateVariant(
// 		file,
// 		string(metadata.MIMEType),
// 		targetWidth,
// 	)
// 	if err != nil {
// 		return "", fmt.Errorf("generate %s variant failed: %w", resolution, err)
// 	}

// 	imageID := uuid.New()

// 	ext, _ := metadata.MIMEType.Extension()
// 	contentType := string(metadata.MIMEType)

// 	// Since imageGen might fallback webp to jpeg, adjust extension and content type
// 	if metadata.MIMEType == domain.MIMEWEBP {
// 		ext = "jpg"
// 		contentType = "image/jpeg"
// 	}

// 	objectKey := fmt.Sprintf("products/%s/%s_%s.%s",
// 		productID.String(),
// 		imageID.String(),
// 		resolution,
// 		ext,
// 	)

// 	object, err := uS.fileStore.Upload(infraStorage.UploadInput{
// 		Key:           objectKey,
// 		File:          resizedStream,
// 		ContentType:   contentType,
// 		ContentLength: size,
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("upload product %s image: %w", resolution, err)
// 	}

// 	return object.Key, nil
// }
