package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	"service-core/internal/infra/storage"
	infraStorage "service-core/internal/infra/storage"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	"service-core/internal/shared/image"
	"service-core/internal/shared/slug"

	"github.com/google/uuid"
)

var specs = []image.VariantSpec{
	{Type: domain.ResolutionThumbnail, Width: 150},
	{Type: domain.ResolutionPreview, Width: 600},
	{Type: domain.ResolutionDetail, Width: 1200},
}

type CreateProductUsecase struct {
	productRepo    repository.ProductRepository
	productImgRepo repository.ProductImageRepository
	slugGen        slug.Generator
	resolutionGen  image.VariantCreator
	fileStore      infraStorage.Provider
}

func NewCreateProductUsecase(
	productRepo repository.ProductRepository,
	productImgRepo repository.ProductImageRepository,
	slugGen slug.Generator,
	resolutionGen image.VariantCreator,
	fileStore infraStorage.Provider,
) *CreateProductUsecase {
	return &CreateProductUsecase{
		productRepo:    productRepo,
		productImgRepo: productImgRepo,
		slugGen:        slugGen,
		resolutionGen:  resolutionGen,
		fileStore:      fileStore,
	}
}

type CreateProductImageInput struct {
	Data         []byte
	OriginalName string
	MIMEType     string
	SizeBytes    int64
	IsPrimary    bool
	DisplayOrder int
}

type CreateProductInput struct {
	SKU         string
	Name        string
	Description *string
	Status      domain.ProductStatus
	Price       int64
	Weight      *float64
	Images      []CreateProductImageInput
}

type productToVariants struct {
	index       int
	variantType image.ResolutionType
}

func (u *CreateProductUsecase) Execute(input CreateProductInput) error {
	var (
		// Input payloads prepared for storage
		// upload operations.
		uploadImages []storage.UploadInput

		// Domain-level product image entities
		// before persistence.
		productImages []domain.ProductImage

		// Index mapping between uploaded results
		// and product images.Used to re-associate
		// generated variant outputs with their
		// source image.
		mappings []productToVariants

		now = time.Now()
	)

	product := &domain.Product{
		ID:          uuid.New(),
		SKU:         input.SKU,
		Name:        input.Name,
		Slug:        u.slugGen.Generate(input.Name),
		Description: input.Description,
		Status:      input.Status,
		Price:       input.Price,
		Weight:      input.Weight,
		CreatedAt:   now,
	}
	if err := product.Validate(); err != nil {
		if errors.Is(err, domain.ErrInvalidProductName) ||
			errors.Is(err, domain.ErrInvalidProductPrice) {
			return appErr.NewInvalidInput(err.Error())
		}

		return err
	}

	for _, img := range input.Images {
		productImage := domain.ProductImage{
			ID:           uuid.New(),
			ProductID:    product.ID,
			Variants:     make(map[image.ResolutionType]domain.ImageVariant),
			IsPrimary:    img.IsPrimary,
			DisplayOrder: img.DisplayOrder,
			Metadata: domain.ProductImageMetadata{
				OriginalName: img.OriginalName,
				MIMEType:     image.MIME(img.MIMEType),
				SizeBytes:    img.SizeBytes,
			},
			CreatedAt: now,
		}

		// UploadMany returns responses in the same order as uploadImages,
		// but each ObjectResponse only contains the stored result data
		// (such as the Key) and does not include domain context such as
		// the variant type or the ProductImage it belongs to.
		//
		// It is necessary to capture the index of the incoming ProductImage
		// before append-ing it so that the location of each upload response
		// can be reconstructed.
		productImageIndex := len(productImages)
		productImages = append(productImages, productImage)

		variants, err := u.resolutionGen.GenerateVariants(
			img.Data,
			image.MIME(img.MIMEType),
			specs,
		)
		if err != nil {
			return fmt.Errorf("failed to create variants of %s: %w",
				img.OriginalName,
				err,
			)
		}

		for _, variant := range variants {
			key, err := productImage.BuildObjectKey()
			if err != nil {
				return fmt.Errorf("failed to build %s key: %w",
					variant.Type,
					err,
				)
			}

			uploadImages = append(uploadImages, storage.UploadInput{
				Bucket:      "public-assets",
				Key:         key,
				File:        bytes.NewReader(variant.Data),
				ContentType: string(variant.MIMEType),
			})

			// Record positional metadata for this upload.
			//
			//   - uploadImages[n] <-> mappings[n]
			//
			// This makes it possible to restore the domain context after
			// UploadMany returns a value, because ObjectResponse doesn't
			// tell which variant/image it has.
			mappings = append(mappings, productToVariants{
				index:       productImageIndex,
				variantType: variant.Type,
			})
		}
	}

	responses, err := u.fileStore.UploadMany(uploadImages)
	if err != nil {
		return fmt.Errorf("failed to upload images: %w", err)
	}

	// Responses maintain the order of UploadMany input.
	// So, responses[i] corresponds to uploadImages[i], and mappings[i]
	// tells exactly which ProductImage and variant slot should receive
	// the uploaded key.
	for i, resp := range responses {
		m := mappings[i]

		productImages[m.index].Variants[m.variantType] =
			domain.ImageVariant{
				Type: m.variantType,
				Key:  resp.Key,
			}
	}

	if err := u.productRepo.CreateProduct(product); err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	if err := u.productImgRepo.Create(productImages); err != nil {
		return fmt.Errorf("failed to save product image: %w", err)
	}

	return nil
}
