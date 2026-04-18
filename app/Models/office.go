package models

import (
	"time"

	"gorm.io/gorm"
)

type Office struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Type       string         `json:"type" gorm:"type:varchar(64);not null;index"`
	Name       string         `json:"name" gorm:"type:varchar(191);not null"`
	Code       string         `json:"code" gorm:"type:varchar(64);not null;uniqueIndex"`
	Address    *string        `json:"address,omitempty" gorm:"type:text"`
	Phone      *string        `json:"phone,omitempty" gorm:"type:varchar(64)"`
	ProvinceID *uint          `json:"province_id,omitempty" gorm:"column:province_id;index"`
	CityID     *uint          `json:"city_id,omitempty" gorm:"column:city_id;index"`
	Status     uint8          `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`
	ImageURL   *string        `json:"image_url,omitempty" gorm:"column:image_url;type:varchar(512)"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Office) TableName() string {
	return "offices"
}
