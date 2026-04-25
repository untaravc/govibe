package shipmentcontroller

type shipmentDetailInput struct {
	ItemName   string  `json:"item_name" validate:"required,min=1,max=191"`
	ItemPrice  float64 `json:"item_price" validate:"omitempty,gte=0"`
	CategoryID *uint   `json:"category_id" validate:"omitempty"`
}

type shipmentTransitInput struct {
	OfficeID uint `json:"office_id" validate:"required,gt=0"`
}

type createShipmentRequest struct {
	Code                string                 `json:"code" validate:"omitempty,min=1,max=64"`
	CustomerName        string                 `json:"customer_name" validate:"required,min=2,max=191"`
	OfficeOriginID      uint                   `json:"office_origin_id" validate:"required,gt=0"`
	OfficeDestinationID uint                   `json:"office_destination_id" validate:"required,gt=0"`
	CustomerPhone       string                 `json:"customer_phone" validate:"required,min=3,max=64"`
	CustomerEmail       *string                `json:"customer_email" validate:"omitempty,email,max=191"`
	Price               float64                `json:"price" validate:"omitempty,gte=0"`
	Wight               float64                `json:"wight" validate:"omitempty,gte=0"`
	Length              float64                `json:"length" validate:"omitempty,gte=0"`
	Width               float64                `json:"width" validate:"omitempty,gte=0"`
	Height              float64                `json:"height" validate:"omitempty,gte=0"`
	PriceType           string                 `json:"price_type" validate:"required,oneof=dimension weight"`
	Status              uint                   `json:"status" validate:"omitempty"`
	Details             []shipmentDetailInput  `json:"details" validate:"omitempty,dive"`
	Transits            []shipmentTransitInput `json:"transits" validate:"omitempty,dive"`
}

type updateShipmentRequest struct {
	Code                *string                 `json:"code" validate:"omitempty,min=1,max=64"`
	CustomerName        *string                 `json:"customer_name" validate:"omitempty,min=2,max=191"`
	OfficeOriginID      *uint                   `json:"office_origin_id" validate:"omitempty,gt=0"`
	OfficeDestinationID *uint                   `json:"office_destination_id" validate:"omitempty,gt=0"`
	CustomerPhone       *string                 `json:"customer_phone" validate:"omitempty,min=3,max=64"`
	CustomerEmail       *string                 `json:"customer_email" validate:"omitempty,email,max=191"`
	Price               *float64                `json:"price" validate:"omitempty,gte=0"`
	Wight               *float64                `json:"wight" validate:"omitempty,gte=0"`
	Length              *float64                `json:"length" validate:"omitempty,gte=0"`
	Width               *float64                `json:"width" validate:"omitempty,gte=0"`
	Height              *float64                `json:"height" validate:"omitempty,gte=0"`
	PriceType           *string                 `json:"price_type" validate:"omitempty,oneof=dimension weight"`
	Status              *uint                   `json:"status" validate:"omitempty"`
	Details             *[]shipmentDetailInput  `json:"details" validate:"omitempty,dive"`
	Transits            *[]shipmentTransitInput `json:"transits" validate:"omitempty,dive"`
}
