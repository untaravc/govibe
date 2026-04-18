package rolecontroller

type createRoleRequest struct {
	Role   string `json:"role" validate:"required,min=2,max=64"`
	Name   string `json:"name" validate:"required,min=2,max=191"`
	Status uint8  `json:"status" validate:"required,oneof=0 1"`
}

type updateRoleRequest struct {
	Role   *string `json:"role" validate:"omitempty,min=2,max=64"`
	Name   *string `json:"name" validate:"omitempty,min=2,max=191"`
	Status *uint8  `json:"status" validate:"omitempty,oneof=0 1"`
}
