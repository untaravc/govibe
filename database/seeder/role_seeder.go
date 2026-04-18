package seeder

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type role struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Role   string `gorm:"type:varchar(64);not null;uniqueIndex"`
	Name   string `gorm:"type:varchar(191);not null"`
	Status uint8  `gorm:"type:tinyint unsigned;not null;default:1"`
}

func (role) TableName() string { return "roles" }

func SeedRoles(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	roles := []role{
		{ID: 1, Role: "superadmin", Name: "Superadmin", Status: 1},
		{ID: 2, Role: "admin", Name: "Admin", Status: 1},
		{ID: 3, Role: "guest", Name: "Guest", Status: 1},
	}

	for _, r := range roles {
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "role"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"status",
			}),
		}).Create(&r).Error; err != nil {
			return err
		}
	}

	return nil
}
