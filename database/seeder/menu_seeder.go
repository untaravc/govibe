package seeder

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menu struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"type:varchar(191);not null"`
	Icon     *string
	Slug     string `gorm:"type:varchar(191);not null;uniqueIndex"`
	ParentID *uint  `gorm:"column:parent_id"`
	Order    int    `gorm:"column:order;not null;default:0"`
	Link     *string
	Status   uint8 `gorm:"type:tinyint unsigned;not null;default:1"`
}

func (menu) TableName() string { return "menus" }

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func SeedMenus(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	// Parents
	if err := upsertMenu(db, menu{
		Name:   "Dashboard",
		Icon:   strPtr("mdi:view-dashboard-outline"),
		Slug:   "dashboard",
		Order:  1,
		Link:   strPtr("/admin"),
		Status: 1,
	}); err != nil {
		return err
	}

	if err := upsertMenu(db, menu{
		Name:   "Config",
		Icon:   strPtr("mdi:cog-outline"),
		Slug:   "config",
		Order:  2,
		Link:   nil,
		Status: 1,
	}); err != nil {
		return err
	}

	var config menu
	if err := db.Select("id").Where("slug = ?", "config").First(&config).Error; err != nil {
		return err
	}

	// Children under Config
	children := []menu{
		{
			Name:     "User",
			Icon:     strPtr("mdi:account-multiple-outline"),
			Slug:     "user",
			ParentID: &config.ID,
			Order:    1,
			Link:     strPtr("/admin/users"),
			Status:   1,
		},
		{
			Name:     "Role",
			Icon:     strPtr("mdi:shield-account-outline"),
			Slug:     "role",
			ParentID: &config.ID,
			Order:    2,
			Link:     strPtr("/admin/roles"),
			Status:   1,
		},
		{
			Name:     "Menu Role",
			Icon:     strPtr("mdi:shield-key-outline"),
			Slug:     "menu-role",
			ParentID: &config.ID,
			Order:    3,
			Link:     strPtr("/admin/menu-roles"),
			Status:   1,
		},
	}

	for _, m := range children {
		if err := upsertMenu(db, m); err != nil {
			return err
		}
	}

	return nil
}

func upsertMenu(db *gorm.DB, m menu) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"icon",
			"parent_id",
			"order",
			"link",
			"status",
		}),
	}).Create(&m).Error
}
