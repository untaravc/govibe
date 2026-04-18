package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"type:varchar(191);not null"`
	Icon     *string `json:"icon,omitempty" gorm:"type:varchar(191)"`
	Slug     string  `json:"slug" gorm:"type:varchar(191);not null;uniqueIndex"`
	ParentID *uint   `json:"parent_id,omitempty" gorm:"column:parent_id;index"`
	Order    int     `json:"order" gorm:"column:order;not null;default:0"`
	Link     *string `json:"link,omitempty" gorm:"type:varchar(255)"`
	Status   uint8   `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Menu) TableName() string {
	return "menus"
}
