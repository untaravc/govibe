package categorycontroller

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

type CategoryController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *CategoryController {
	return &CategoryController{db: db}
}

func (ctl *CategoryController) Index(c *fiber.Ctx) error {
	var categories []models.Category
	if err := ctl.db.Order("id desc").Find(&categories).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"categories": categories,
	})
}

func (ctl *CategoryController) Show(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	category, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"category": category,
	})
}

func (ctl *CategoryController) Store(c *fiber.Ctx) error {
	var req createCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Section = strings.TrimSpace(req.Section)
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.TrimSpace(req.Slug)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	category := models.Category{
		Section: req.Section,
		Name:    req.Name,
		Slug:    req.Slug,
		Status:  req.Status,
	}

	if err := ctl.db.Create(&category).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "created", fiber.Map{
		"category": category,
	})
}

func (ctl *CategoryController) Update(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req updateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if req.Section != nil {
		trimmed := strings.TrimSpace(*req.Section)
		req.Section = &trimmed
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Slug != nil {
		trimmed := strings.TrimSpace(*req.Slug)
		req.Slug = &trimmed
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	updates := make(map[string]any, 4)
	if req.Section != nil {
		updates["section"] = *req.Section
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
	}

	tx := ctl.db.Model(&models.Category{}).Where("id = ?", uint(id)).Updates(updates)
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "category not found")
	}

	category, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"category": category,
	})
}

func (ctl *CategoryController) Destroy(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := ctl.db.Delete(&models.Category{}, uint(id))
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "category not found")
	}

	return response.OK(c, "deleted", fiber.Map{
		"deleted": true,
	})
}

func (ctl *CategoryController) getByID(id uint) (models.Category, error) {
	var category models.Category
	if err := ctl.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Category{}, fiber.NewError(fiber.StatusNotFound, "category not found")
		}
		return models.Category{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return category, nil
}
