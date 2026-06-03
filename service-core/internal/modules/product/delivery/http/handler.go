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

	input := usecase.FindProductsInput{
		Page:  page,
		Limit: limit,
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
		result := productCatalogResponse{
			ID:         p.Product.ID,
			SKU:        p.Product.SKU,
			Name:       p.Product.Name,
			Slug:       p.Product.Slug,
			Price:      p.Product.Price,
			TotalStock: p.Inventory.TotalStock,
			Image: productImageResponse{
				Thumbnail: &p.Images.Thumbnail,
			},
		}

		result.IsAvailable =
			productStatusDTO(p.Product.Status) == ProductStatusActive &&
				(p.Inventory.TotalStock-p.Inventory.ReservedStock) > 0

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
	productID, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	product, err := h.getProduct.Execute(r.Context(), productID)
	if err != nil {
		return err
	}
	if product == nil {
		return apperrors.NewNotFound("product not found")
	}

	var available int
	inventories := make([]productInventoryView, 0, len(product.ShopInventories))
	for _, inventory := range product.ShopInventories {
		available += inventory.Available()

		inventories = append(inventories, productInventoryView{
			ID:         inventory.ID,
			ShopID:     inventory.ShopID,
			TotalStock: inventory.TotalStock,
			Available:  inventory.Available(),
		})
	}

	images := make([]productImageResponse, 0, len(product.Images))
	for _, img := range product.Images {
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

	response := productDetailResponse{
		ID:          product.Product.ID,
		SKU:         product.Product.SKU,
		Name:        product.Product.Name,
		Slug:        product.Product.Slug,
		Description: product.Product.Description,
		Price:       product.Product.Price,
		Weight:      product.Product.Weight,
		TotalStock:  available,
		UpdatedAt:   product.Product.UpdatedAt,
		Images:      images,
	}
	response.IsAvailable =
		productStatusDTO(product.Product.Status) == ProductStatusActive &&
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

	productID, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	input := usecase.AddProductImageInput{
		ProductID: productID,
		Images:    images,
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
