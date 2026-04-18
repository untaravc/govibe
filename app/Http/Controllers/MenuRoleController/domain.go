package menurolecontroller

type grantItem struct {
	MenuID uint   `json:"menu_id" validate:"required"`
	Method string `json:"method" validate:"required"`
}

type saveMenuRoleRequest struct {
	RoleID uint        `json:"role_id" validate:"required"`
	Grants []grantItem `json:"grants" validate:"omitempty"`
}
