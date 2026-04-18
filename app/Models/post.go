package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string         `json:"title" gorm:"type:varchar(191);not null"`
	Subtitle   *string        `json:"subtitle,omitempty" gorm:"type:varchar(255)"`
	Slug       string         `json:"slug" gorm:"type:varchar(191);not null;uniqueIndex"`
	Content    *string        `json:"content,omitempty" gorm:"type:longtext"`
	Status     uint8          `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`
	ImageURL   *string        `json:"image_url,omitempty" gorm:"column:image_url;type:varchar(512)"`
	UserID     uint           `json:"user_id" gorm:"column:user_id;not null;index"`
	ReleaseAt  *time.Time     `json:"release_at,omitempty" gorm:"column:release_at;index"`
	CategoryID *uint          `json:"category_id,omitempty" gorm:"column:category_id;index"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Post) TableName() string {
	return "posts"
}
