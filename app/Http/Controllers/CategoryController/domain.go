package categorycontroller

type createCategoryRequest struct {
	Section string `json:"section" validate:"required,min=2,max=64"`
	Name    string `json:"name" validate:"required,min=2,max=191"`
	Slug    string `json:"slug" validate:"required,min=2,max=191"`
	Status  uint8  `json:"status" validate:"required,oneof=0 1"`
}

type updateCategoryRequest struct {
	Section *string `json:"section" validate:"omitempty,min=2,max=64"`
	Name    *string `json:"name" validate:"omitempty,min=2,max=191"`
	Slug    *string `json:"slug" validate:"omitempty,min=2,max=191"`
	Status  *uint8  `json:"status" validate:"omitempty,oneof=0 1"`
}
