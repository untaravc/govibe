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
