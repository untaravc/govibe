package usercontroller

type createUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=191"`
	Email    string `json:"email" validate:"required,email,max=191"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type updateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=191"`
	Email    *string `json:"email" validate:"omitempty,email,max=191"`
	Password *string `json:"password" validate:"omitempty,min=6,max=255"`
}
