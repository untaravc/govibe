package models

import (
	"time"

	"gorm.io/gorm"
)

type Shipment struct {
	ID                  uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	Code                string           `json:"code" gorm:"type:varchar(64);not null;uniqueIndex"`
	CustomerName        string           `json:"customer_name" gorm:"column:customer_name;type:varchar(191);not null"`
	OfficeOriginID      uint             `json:"office_origin_id" gorm:"column:office_origin_id;not null;index"`
	OfficeDestinationID uint             `json:"office_destination_id" gorm:"column:office_destination_id;not null;index"`
	CustomerPhone       string           `json:"customer_phone" gorm:"column:customer_phone;type:varchar(64);not null"`
	CustomerEmail       *string          `json:"customer_email,omitempty" gorm:"column:customer_email;type:varchar(191)"`
	Price               float64          `json:"price" gorm:"type:decimal(15,2);not null;default:0.00"`
	UserID              uint             `json:"user_id" gorm:"column:user_id;not null;index"`
	Wight               float64          `json:"wight" gorm:"column:wight;type:decimal(12,3);not null;default:0.000"`
	Length              float64          `json:"length" gorm:"column:length;type:decimal(12,3);not null;default:0.000"`
	Width               float64          `json:"width" gorm:"column:width;type:decimal(12,3);not null;default:0.000"`
	Height              float64          `json:"height" gorm:"column:height;type:decimal(12,3);not null;default:0.000"`
	PriceType           string           `json:"price_type" gorm:"column:price_type;type:enum('dimension','weight');not null;default:'weight'"`
	Status              uint             `json:"status" gorm:"type:mediumint unsigned;not null;default:0;index"`
	Details             []ShipmentDetail `json:"details,omitempty" gorm:"foreignKey:ShipmentID"`
	Logs                []ShipmentLog    `json:"logs,omitempty" gorm:"foreignKey:ShipmentID"`
	CreatedAt           time.Time        `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt           time.Time        `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt           gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Shipment) TableName() string { return "shipments" }

type ShipmentDetail struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipmentID uint           `json:"shipment_id" gorm:"column:shipment_id;not null;index"`
	ItemName   string         `json:"item_name" gorm:"column:item_name;type:varchar(191);not null"`
	ItemPrice  float64        `json:"item_price" gorm:"column:item_price;type:decimal(15,2);not null;default:0.00"`
	CategoryID *uint          `json:"category_id,omitempty" gorm:"column:category_id;index"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (ShipmentDetail) TableName() string { return "shipment_details" }

type ShipmentLog struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipmentID    uint           `json:"shipment_id" gorm:"column:shipment_id;not null;index"`
	OfficeID      uint           `json:"office_id" gorm:"column:office_id;not null;index"`
	ArrivalTime   *time.Time     `json:"arrival_time,omitempty" gorm:"column:arrival_time;index"`
	DepartureTime *time.Time     `json:"departure_time,omitempty" gorm:"column:departure_time;index"`
	UserID        uint           `json:"user_id" gorm:"column:user_id;not null;index"`
	Note          *string        `json:"note,omitempty" gorm:"column:note;type:text"`
	Status        uint           `json:"status" gorm:"type:mediumint unsigned;not null;default:0;index"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (ShipmentLog) TableName() string { return "shipment_logs" }
