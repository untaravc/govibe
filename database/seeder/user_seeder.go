package seeder

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type user struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"type:varchar(191);not null"`
	Email     string         `gorm:"type:varchar(191);not null;uniqueIndex"`
	Status    uint8          `gorm:"type:tinyint unsigned;not null;default:1"`
	RoleID    *uint          `gorm:"column:role_id;index"`
	AuthType  string         `gorm:"type:enum('email','phone');not null;default:'email'"`
	Password  string         `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (user) TableName() string { return "users" }

func SeedUsers(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	pw, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hash := string(pw)

	role1 := uint(1)
	role2 := uint(2)
	role3 := uint(3)

	users := []user{
		{Name: "Superadmin", Email: "superadmin@admin.com", Password: hash, RoleID: &role1, Status: 1, AuthType: "email"},
		{Name: "Admin", Email: "admin@admin.com", Password: hash, RoleID: &role2, Status: 1, AuthType: "email"},
		{Name: "User", Email: "user@user.com", Password: hash, RoleID: &role3, Status: 1, AuthType: "email"},
	}

	for _, u := range users {
		// If the user exists but is soft-deleted, restore it.
		var existing user
		if err := db.Unscoped().Select("id", "deleted_at").Where("email = ?", u.Email).First(&existing).Error; err == nil {
			if existing.DeletedAt.Valid {
				if err := db.Unscoped().Model(&user{}).Where("id = ?", existing.ID).Update("deleted_at", nil).Error; err != nil {
					return err
				}
			}
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"status",
				"role_id",
				"auth_type",
				"password",
			}),
		}).Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}
