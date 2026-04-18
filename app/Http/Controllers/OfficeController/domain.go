package officecontroller

type createOfficeRequest struct {
	Type       string  `json:"type" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Code       string  `json:"code" validate:"required"`
	Address    *string `json:"address" validate:"omitempty"`
	Phone      *string `json:"phone" validate:"omitempty"`
	ProvinceID *uint   `json:"province_id" validate:"omitempty"`
	CityID     *uint   `json:"city_id" validate:"omitempty"`
	Status     uint8   `json:"status" validate:"omitempty"`
	ImageURL   *string `json:"image_url" validate:"omitempty"`
}

type updateOfficeRequest struct {
	Type       *string `json:"type" validate:"omitempty"`
	Name       *string `json:"name" validate:"omitempty"`
	Code       *string `json:"code" validate:"omitempty"`
	Address    *string `json:"address" validate:"omitempty"`
	Phone      *string `json:"phone" validate:"omitempty"`
	ProvinceID *uint   `json:"province_id" validate:"omitempty"`
	CityID     *uint   `json:"city_id" validate:"omitempty"`
	Status     *uint8  `json:"status" validate:"omitempty"`
	ImageURL   *string `json:"image_url" validate:"omitempty"`
}
