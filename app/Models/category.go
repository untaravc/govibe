package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Section   string         `json:"section" gorm:"type:varchar(64);not null;index"`
	Name      string         `json:"name" gorm:"type:varchar(191);not null"`
	Slug      string         `json:"slug" gorm:"type:varchar(191);not null;uniqueIndex"`
	Status    uint8          `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Category) TableName() string {
	return "categories"
}
