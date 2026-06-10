package http

type addMerchantAccountRequest struct {
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type createMerchantRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logo_url"`
	BannerUrl   *string `json:"banner_url"`
}
