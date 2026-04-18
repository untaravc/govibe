package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuRole struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID    uint           `json:"role_id" gorm:"column:role_id;not null;index"`
	MenuID    uint           `json:"menu_id" gorm:"column:menu_id;not null;index"`
	Method    string         `json:"method" gorm:"type:enum('get','update','create','delete','show');not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (MenuRole) TableName() string {
	return "menu_roles"
}
