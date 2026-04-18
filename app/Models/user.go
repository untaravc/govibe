package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(191);not null"`
	Email     string         `json:"email" gorm:"type:varchar(191);not null;uniqueIndex"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "users"
}
