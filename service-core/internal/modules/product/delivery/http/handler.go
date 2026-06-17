package http

import (
	"errors"
	"io"
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appmultipart "service-core/internal/common/http/multipart"
	"service-core/internal/modules/product/usecase"
)

type ProductHandler struct {
	findProducts    *usecase.FindProductsUsecase
	getProduct      *usecase.GetProductUsecase
	createProduct   *usecase.CreateProductUsecase
	addProductImage *usecase.AddProductImagesUsecase
}

func NewProductHandler(
	findProducts *usecase.FindProductsUsecase,
	getProduct *usecase.GetProductUsecase,
	createProduct *usecase.CreateProductUsecase,
	addProductImage *usecase.AddProductImagesUsecase,
) *ProductHandler {
	return &ProductHandler{
		findProducts:    findProducts,
		getProduct:      getProduct,
		createProduct:   createProduct,
		addProductImage: addProductImage,
	}
}

func (h *ProductHandler) FindProducts(w http.ResponseWriter, r *http.Request) error {
	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	name := apphttp.Query(r, "name")
	id := apphttp.Query(r, "id")
	sort := apphttp.Query(r, "sort")

	input := usecase.FindProductsInput{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
	if name != "" {
		input.Name = &name
	}
	if id != "" {
		input.ID = &id
	}

	products, total, err := h.findProducts.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]productCatalogResponse, 0, len(products))
	for _, p := range products {
		availabilities := make([]productAvailabilityResponse, 0, len(p.Availability))
		for _, avail := range p.Availability {
			availabilities = append(availabilities, productAvailabilityResponse{
				Slug:  avail.ShopName,
				Name:  avail.ShopSlug,
				Stock: avail.Stock,
			})
		}

		var result productCatalogResponse

		result = productCatalogResponse{
			productBaseResponse: productBaseResponse{
				ID:     p.Product.ID,
				SKU:    p.Product.SKU,
				Name:   p.Product.Name,
				Slug:   p.Product.Slug,
				Status: string(p.Product.Status),
				IsAvailable: productStatusDTO(p.Product.Status) == ProductStatusActive &&
					p.Inventory.TotalStock > 0,
				Price:      p.Product.Price,
				TotalStock: p.Inventory.TotalStock,
				Banner: productImageResponse{
					Thumbnail: &p.Images.Thumbnail,
				},
			},
			Availability: availabilities,
		}

		if *result.Banner.Thumbnail == "" {
			result.Banner = productImageResponse{
				Thumbnail: nil,
			}
		}

		results = append(results, result)
	}

	response := map[string]interface{}{
		"products": results,
		"page":     page,
		"limit":    limit,
		"total":    total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) error {
	productSlug := apphttp.Param(r, "slug")

	productDetail, err := h.getProduct.
		Execute(r.Context(), productSlug)
	if err != nil {
		return err
	}
	if productDetail == nil {
		return apperrors.NewNotFound("product not found")
	}

	var available int
	inventories := make([]productInventoryView, 0, len(productDetail.ShopInventories))
	for _, inventory := range productDetail.ShopInventories {
		available += inventory.Available()

		inventories = append(inventories, productInventoryView{
			ID:         inventory.ID,
			ShopID:     inventory.ShopID,
			TotalStock: inventory.TotalStock,
			Available:  inventory.Available(),
		})
	}

	images := make([]productImageResponse, 0, len(productDetail.Images))
	for _, img := range productDetail.Images {
		image := productImageResponse{}

		if img.Thumbnail != "" {
			image.Thumbnail = &img.Thumbnail
		}

		if img.Preview != "" {
			image.Preview = &img.Preview
		}

		if img.Detail != "" {
			image.Detail = &img.Detail
		}

		images = append(images, image)
	}

	availabilities := make(
		[]productAvailabilityResponse,
		0,
		len(productDetail.Availability),
	)
	for _, avail := range productDetail.Availability {
		availabilities = append(availabilities, productAvailabilityResponse{
			Slug:  avail.ShopName,
			Name:  avail.ShopSlug,
			Stock: avail.Stock,
		})
	}

	response := productDetailResponse{
		productBaseResponse: productBaseResponse{
			ID:     productDetail.Product.ID,
			SKU:    productDetail.Product.SKU,
			Name:   productDetail.Product.Name,
			Slug:   productDetail.Product.Slug,
			Status: string(productDetail.Product.Status),
			IsAvailable: productStatusDTO(productDetail.Product.Status) == ProductStatusActive &&
				available > 0,
			Price:      productDetail.Product.Price,
			TotalStock: available,
			Banner:     images[0],
		},
		Description:  productDetail.Product.Description,
		Weight:       productDetail.Product.Weight,
		UpdatedAt:    productDetail.Product.UpdatedAt,
		Gallery:      images,
		Availability: availabilities,
	}
	response.IsAvailable =
		productStatusDTO(productDetail.Product.Status) == ProductStatusActive &&
			available > 0

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	var req createProductRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.Name == "" {
		return apperrors.NewBadRequest("invalid name")
	}
	if req.SKU == "" {
		return apperrors.NewBadRequest("invalid sku")
	}
	if req.Price < 0 {
		return apperrors.NewBadRequest("invalid price")
	}
	if !req.Status.isStatusValid() {
		return apperrors.NewBadRequest("invalid status")
	}

	input := usecase.CreateProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Status:      string(req.Status),
		Price:       req.Price,
		Weight:      req.Weight,
	}

	err := h.createProduct.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "product successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *ProductHandler) AddProductImages(w http.ResponseWriter, r *http.Request) error {
	files, err := appmultipart.ParseMultiple(r, "image", 16<<20)
	if err != nil {
		return err
	}

	images := make([]usecase.ProductImageInput, 0, len(files))
	for i, file := range files {
		data, err := io.ReadAll(file.File)
		if err != nil {
			return apperrors.NewInternal(errors.New("failed to read uploaded file"))
		}

		err = file.File.Close()
		if err != nil {
			return apperrors.NewInternal(errors.New("failed to close uploaded file"))
		}

		images = append(images, usecase.ProductImageInput{
			Data:         data,
			OriginalName: file.Filename,
			MIMEType:     file.ContentType,
			SizeBytes:    file.Size,
			IsPrimary:    i == 0,
			DisplayOrder: i,
		})
	}

	productSlug := apphttp.Param(r, "slug")
	input := usecase.AddProductImageInput{
		ProductSlug: productSlug,
		Images:      images,
	}

	err = h.addProductImage.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "product image successfully added",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (s productStatusDTO) isStatusValid() bool {
	switch s {
	case ProductStatusActive,
		ProductStatusInactive,
		ProductStatusArchived:
		return true
	default:
		return false
	}
}
