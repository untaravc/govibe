package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string         `json:"name" gorm:"type:varchar(191);not null"`
	Email           string         `json:"email" gorm:"type:varchar(191);not null;uniqueIndex"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty" gorm:"column:email_verified_at"`
	Phone           *string        `json:"phone,omitempty" gorm:"type:varchar(64)"`
	PhoneVerifiedAt *time.Time     `json:"phone_verified_at,omitempty" gorm:"column:phone_verified_at"`
	EmailToken      *string        `json:"-" gorm:"column:email_token;type:varchar(255)"`
	PhoneToken      *string        `json:"-" gorm:"column:phone_token;type:varchar(255)"`
	Image           *string        `json:"image,omitempty" gorm:"type:varchar(512)"`
	Status          uint8          `json:"status" gorm:"type:tinyint unsigned;not null;default:1"`
	RoleID          *uint          `json:"role_id,omitempty" gorm:"column:role_id;index"`
	AuthType        string         `json:"auth_type" gorm:"type:enum('email','phone');not null;default:'email'"`
	RefreshToken    *string        `json:"-" gorm:"column:refresh_token;type:varchar(255)"`
	RefreshTokenExp *time.Time     `json:"refresh_token_expired_at,omitempty" gorm:"column:refresh_token_expired_at"`
	RefreshTokenUpd *time.Time     `json:"refresh_token_updated_at,omitempty" gorm:"column:refresh_token_updated_at"`
	Password        string         `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "users"
}
