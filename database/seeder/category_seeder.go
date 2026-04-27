package seeder

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type category struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Section   string         `gorm:"type:varchar(64);not null;index"`
	Name      string         `gorm:"type:varchar(191);not null"`
	Slug      string         `gorm:"type:varchar(191);not null;uniqueIndex"`
	Status    uint8          `gorm:"type:tinyint unsigned;not null;default:1"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (category) TableName() string { return "categories" }

func SeedCategories(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	items := []category{
		{Section: "office", Name: "Cabang", Slug: "cabang", Status: 1},
		{Section: "office", Name: "Gudang", Slug: "gudang", Status: 1},
		{Section: "post", Name: "Berita", Slug: "berita", Status: 1},
	}

	for _, c := range items {
		// If exists but soft-deleted, restore it.
		var existing category
		if err := db.Unscoped().Select("id", "deleted_at").Where("slug = ?", c.Slug).First(&existing).Error; err == nil {
			if existing.DeletedAt.Valid {
				if err := db.Unscoped().Model(&category{}).Where("id = ?", existing.ID).Update("deleted_at", nil).Error; err != nil {
					return err
				}
			}
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "slug"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"section",
				"name",
				"status",
			}),
		}).Create(&c).Error; err != nil {
			return err
		}
	}

	return nil
}
