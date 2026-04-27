package userofficecontroller

import (
	"errors"

	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/app/Parser"
	appvalidator "govibe/app/Validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserOfficeController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserOfficeController {
	return &UserOfficeController{db: db}
}

func (ctl *UserOfficeController) Index(c *fiber.Ctx) error {
	userID, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var items []models.UserOffice
	if err := ctl.db.Preload("Office").Where("user_id = ?", uint(userID)).Order("id desc").Find(&items).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"user_offices": items,
	})
}

func (ctl *UserOfficeController) Store(c *fiber.Ctx) error {
	userID, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req createUserOfficeRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	status := uint8(1)
	if req.Status != nil {
		status = *req.Status
	}

	var existing models.UserOffice
	findErr := ctl.db.Unscoped().
		Where("user_id = ? AND office_id = ?", uint(userID), req.OfficeID).
		First(&existing).Error
	if findErr == nil {
		tx := ctl.db.Unscoped().
			Model(&models.UserOffice{}).
			Where("id = ?", existing.ID).
			Updates(map[string]any{
				"status":     status,
				"deleted_at": nil,
			})
		if tx.Error != nil {
			return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
		}

		if err := ctl.db.Preload("Office").First(&existing, existing.ID).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return response.OK(c, "ok", fiber.Map{
			"user_office": existing,
		})
	}
	if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, findErr.Error())
	}

	item := models.UserOffice{
		UserID:   uint(userID),
		OfficeID: req.OfficeID,
		Status:   status,
	}
	if err := ctl.db.Create(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := ctl.db.Preload("Office").First(&item, item.ID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.Created(c, "created", fiber.Map{
		"user_office": item,
	})
}

func (ctl *UserOfficeController) Destroy(c *fiber.Ctx) error {
	userID, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	officeID, err := parser.UintParam(c, "office_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := ctl.db.Where("user_id = ? AND office_id = ?", uint(userID), uint(officeID)).Delete(&models.UserOffice{})
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user office not found")
	}

	return response.OK(c, "deleted", fiber.Map{
		"deleted": true,
	})
}
