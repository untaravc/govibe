package userofficecontroller

type createUserOfficeRequest struct {
	OfficeID uint   `json:"office_id" validate:"required,gte=1"`
	Status   *uint8 `json:"status" validate:"omitempty,gte=0,lte=255"`
}
