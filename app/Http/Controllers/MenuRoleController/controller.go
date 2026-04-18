package menurolecontroller

import (
	"strconv"
	"strings"

	"govibe/app/Http/Response"
	"govibe/app/Models"
	appvalidator "govibe/app/Validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MenuRoleController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *MenuRoleController {
	return &MenuRoleController{db: db}
}

func (ctl *MenuRoleController) Index(c *fiber.Ctx) error {
	roleID := uint(parsePositiveIntQuery(c, "role_id", 0))
	if roleID == 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{"role_id": "is required"},
		})
	}

	var menus []models.Menu
	if err := ctl.db.
		Where("status = ?", 1).
		Where("deleted_at IS NULL").
		Order("parent_id asc").
		Order("`order` asc").
		Order("id asc").
		Find(&menus).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var grants []models.MenuRole
	if err := ctl.db.
		Select("id", "role_id", "menu_id", "method").
		Where("role_id = ?", roleID).
		Where("deleted_at IS NULL").
		Order("id asc").
		Find(&grants).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"menus":  menus,
		"grants": grants,
	})
}

func (ctl *MenuRoleController) Save(c *fiber.Ctx) error {
	var req saveMenuRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	clean := make([]models.MenuRole, 0, len(req.Grants))
	incomingKeys := make(map[string]struct{}, len(req.Grants))
	for _, g := range req.Grants {
		method := strings.ToLower(strings.TrimSpace(g.Method))
		if !isValidMethod(method) {
			return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
				"errors": fiber.Map{"method": "is invalid"},
			})
		}
		key := strings.TrimSpace(strings.Join([]string{itoa(g.MenuID), method}, ":"))
		if _, ok := incomingKeys[key]; ok {
			continue
		}
		incomingKeys[key] = struct{}{}
		clean = append(clean, models.MenuRole{
			RoleID: req.RoleID,
			MenuID: g.MenuID,
			Method: method,
		})
	}

	var existingAll []models.MenuRole
	if err := ctl.db.
		Unscoped().
		Select("id", "menu_id", "method", "deleted_at").
		Where("role_id = ?", req.RoleID).
		Find(&existingAll).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	existingByKey := make(map[string]models.MenuRole, len(existingAll))
	existingActiveKeys := make(map[string]struct{}, len(existingAll))
	for _, g := range existingAll {
		method := strings.ToLower(strings.TrimSpace(g.Method))
		key := strings.TrimSpace(strings.Join([]string{itoa(g.MenuID), method}, ":"))
		existingByKey[key] = g
		if !g.DeletedAt.Valid {
			existingActiveKeys[key] = struct{}{}
		}
	}

	unchanged := len(incomingKeys) == len(existingActiveKeys)
	if unchanged {
		for k := range incomingKeys {
			if _, ok := existingActiveKeys[k]; !ok {
				unchanged = false
				break
			}
		}
	}
	if unchanged {
		return response.OK(c, "no changes", fiber.Map{
			"saved":   true,
			"changed": false,
		})
	}

	if err := ctl.db.Transaction(func(tx *gorm.DB) error {
		// Restore or delete existing rows depending on incoming keys.
		for key, existing := range existingByKey {
			_, shouldBeActive := incomingKeys[key]

			if shouldBeActive {
				if existing.DeletedAt.Valid {
					if err := tx.Unscoped().
						Model(&models.MenuRole{}).
						Where("id = ?", existing.ID).
						Update("deleted_at", nil).Error; err != nil {
						return err
					}
				}
				continue
			}

			if !existing.DeletedAt.Valid {
				if err := tx.Delete(&models.MenuRole{}, existing.ID).Error; err != nil {
					return err
				}
			}
		}

		// Create missing rows (never recreates soft-deleted duplicates because we restore those above).
		toCreate := make([]models.MenuRole, 0, len(clean))
		for _, g := range clean {
			key := strings.TrimSpace(strings.Join([]string{itoa(g.MenuID), g.Method}, ":"))
			if _, ok := existingByKey[key]; ok {
				continue
			}
			toCreate = append(toCreate, g)
		}
		if len(toCreate) > 0 {
			if err := tx.Create(&toCreate).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "saved", fiber.Map{
		"saved":   true,
		"changed": true,
	})
}

func isValidMethod(m string) bool {
	switch m {
	case "get", "create", "update", "delete", "show":
		return true
	default:
		return false
	}
}

func parsePositiveIntQuery(c *fiber.Ctx, key string, fallback int) int {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func itoa(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
