package authcontroller

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=191"`
}

type registerRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=191"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=191"`
}

type requestResetPasswordRequest struct {
	Email string `json:"email" validate:"omitempty,email,max=191"`
	Phone string `json:"phone" validate:"omitempty,min=6,max=64"`
}

type updatePasswordWithTokenRequest struct {
	EmailToken  string `json:"email_token" validate:"required,min=10,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=191"`
}

type updateProfileRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=191"`
	Phone    *string `json:"phone" validate:"omitempty,min=6,max=64"`
	Email    *string `json:"email" validate:"omitempty,email,max=191"`
	Password *string `json:"password" validate:"omitempty,min=8,max=191"`
}
