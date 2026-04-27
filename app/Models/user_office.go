package models

import (
	"time"

	"gorm.io/gorm"
)

type UserOffice struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"column:user_id;index;not null"`
	OfficeID  uint           `json:"office_id" gorm:"column:office_id;index;not null"`
	Status    uint8          `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`
	Office    *Office        `json:"office,omitempty" gorm:"foreignKey:OfficeID"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (UserOffice) TableName() string {
	return "user_offices"
}
