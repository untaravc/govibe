package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Role      string         `json:"role" gorm:"type:varchar(64);not null;uniqueIndex"`
	Name      string         `json:"name" gorm:"type:varchar(191);not null"`
	Status    uint8          `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Role) TableName() string {
	return "roles"
}
