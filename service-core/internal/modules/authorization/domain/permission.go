package domain

import "github.com/google/uuid"

type Permission struct {
	Id   uuid.UUID
	Code string
}

const (
	PermissionShopView          = "shop:view"
	PermissionShopUpdate        = "shop:update"
	PermissionProductCreate     = "product:create"
	PermissionProductUpdate     = "product:update"
	PermissionProductDelete     = "product:delete"
	PermissionInventoryManage   = "inventory:manage"
	PermissionOrderRead         = "order:read"
	PermissionOrderUpdateStatus = "order:update_status"
	PermissionCourierManage     = "courier:manage"
	PermissionAddressManage     = "address:manage"
)

