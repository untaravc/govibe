package rolecontroller

import (
	"errors"
	"strings"

	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/app/Parser"
	appvalidator "govibe/app/Validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RoleController {
	return &RoleController{db: db}
}

func (ctl *RoleController) Index(c *fiber.Ctx) error {
	var roles []models.Role
	if err := ctl.db.Order("id desc").Find(&roles).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return response.OK(c, "ok", fiber.Map{
		"roles": roles,
	})
}

func (ctl *RoleController) Show(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	role, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"role": role,
	})
}

func (ctl *RoleController) Store(c *fiber.Ctx) error {
	var req createRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Role = strings.TrimSpace(req.Role)
	req.Name = strings.TrimSpace(req.Name)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	r := models.Role{
		Role:   req.Role,
		Name:   req.Name,
		Status: req.Status,
	}

	if err := ctl.db.Create(&r).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "created", fiber.Map{
		"role": r,
	})
}

func (ctl *RoleController) Update(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req updateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if req.Role != nil {
		trimmed := strings.TrimSpace(*req.Role)
		req.Role = &trimmed
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	updates := make(map[string]any, 3)
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
	}

	tx := ctl.db.Model(&models.Role{}).Where("id = ?", uint(id)).Updates(updates)
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	role, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"role": role,
	})
}

func (ctl *RoleController) Destroy(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := ctl.db.Delete(&models.Role{}, uint(id))
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	return response.OK(c, "deleted", fiber.Map{
		"deleted": true,
	})
}

func (ctl *RoleController) getByID(id uint) (models.Role, error) {
	var r models.Role
	if err := ctl.db.First(&r, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Role{}, fiber.NewError(fiber.StatusNotFound, "role not found")
		}
		return models.Role{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return r, nil
}
