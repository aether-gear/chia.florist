package http

type courierSelectionRequest struct {
	Code      string `json:"code"`
	IsEnabled bool   `json:"is_enabled"`
}

type configureCourierShopRequest struct {
	Couriers []courierSelectionRequest `json:"couriers"`
}

type updateShopCourierRequest struct {
	Active          bool    `json:"active"`
	Name            *string `json:"name"`
	LocationAddress *string `json:"location_address"`
}

type verifyShopCourierRequest struct {
	Action          string  `json:"action"` // "verify" or "reject"
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

type shopCourierDetailResponse struct {
	ShopID             string  `json:"shop_id"`
	Code               string  `json:"code"`
	BranchName         string  `json:"branch_name"`
	Name               *string `json:"name"`
	LocationAddress    *string `json:"location_address"`
	Active             bool    `json:"active"`
	VerificationStatus string  `json:"verification_status"`
	VerifiedAt         *string `json:"verified_at,omitempty"`
	VerifiedBy         *string `json:"verified_by,omitempty"`
	RejectionReason    *string `json:"rejection_reason,omitempty"`
}
