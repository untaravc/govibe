package usercontroller

type createUserRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=191"`
	Email    string  `json:"email" validate:"required,email,max=191"`
	Password string  `json:"password" validate:"required,min=6,max=255"`
	Phone    *string `json:"phone" validate:"omitempty,max=64"`
	RoleID   *uint   `json:"role_id" validate:"omitempty,gte=1"`
}

type updateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=191"`
	Email    *string `json:"email" validate:"omitempty,email,max=191"`
	Password *string `json:"password" validate:"omitempty,min=6,max=255"`
	Phone    *string `json:"phone" validate:"omitempty,max=64"`
	RoleID   *uint   `json:"role_id" validate:"omitempty,gte=1"`
}
