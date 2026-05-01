package http

type CreateShopRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsActive    string  `json:"is_active"`
}
