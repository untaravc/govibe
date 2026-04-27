package seeder

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type office struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Type      string         `gorm:"type:varchar(64);not null;index"`
	Name      string         `gorm:"type:varchar(191);not null"`
	Code      string         `gorm:"type:varchar(64);not null;uniqueIndex"`
	Status    uint8          `gorm:"type:tinyint unsigned;not null;default:1"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (office) TableName() string { return "offices" }

func SeedOffices(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	items := []office{
		{Type: "cabang", Code: "C001", Name: "Cabang 1", Status: 1},
		{Type: "cabang", Code: "C002", Name: "Cabang 2", Status: 1},
		{Type: "cabang", Code: "C003", Name: "Cabang 3", Status: 1},
		{Type: "gudang", Code: "G001", Name: "Gudang 1", Status: 1},
	}

	for _, o := range items {
		// If exists but soft-deleted, restore it.
		var existing office
		if err := db.Unscoped().Select("id", "deleted_at").Where("code = ?", o.Code).First(&existing).Error; err == nil {
			if existing.DeletedAt.Valid {
				if err := db.Unscoped().Model(&office{}).Where("id = ?", existing.ID).Update("deleted_at", nil).Error; err != nil {
					return err
				}
			}
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"type",
				"name",
				"status",
			}),
		}).Create(&o).Error; err != nil {
			return err
		}
	}

	return nil
}
