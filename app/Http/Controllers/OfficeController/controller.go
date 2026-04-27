package officecontroller

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

type OfficeController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *OfficeController {
	return &OfficeController{db: db}
}

func (ctl *OfficeController) Index(c *fiber.Ctx) error {
	var offices []models.Office
	name := strings.TrimSpace(c.Query("name"))
	q := ctl.db.Model(&models.Office{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if err := q.Order("id desc").Find(&offices).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"offices": offices,
	})
}

func (ctl *OfficeController) Show(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	office, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"office": office,
	})
}

func (ctl *OfficeController) Store(c *fiber.Ctx) error {
	var req createOfficeRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Type = strings.TrimSpace(req.Type)
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	if req.Address != nil {
		trimmed := strings.TrimSpace(*req.Address)
		req.Address = &trimmed
	}
	if req.Phone != nil {
		trimmed := strings.TrimSpace(*req.Phone)
		req.Phone = &trimmed
	}
	if req.ImageURL != nil {
		trimmed := strings.TrimSpace(*req.ImageURL)
		req.ImageURL = &trimmed
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	office := models.Office{
		Type:       req.Type,
		Name:       req.Name,
		Code:       req.Code,
		Address:    req.Address,
		Phone:      req.Phone,
		ProvinceID: req.ProvinceID,
		CityID:     req.CityID,
		Status:     req.Status,
		ImageURL:   req.ImageURL,
	}

	if err := ctl.db.Create(&office).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "created", fiber.Map{
		"office": office,
	})
}

func (ctl *OfficeController) Update(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req updateOfficeRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if req.Type != nil {
		trimmed := strings.TrimSpace(*req.Type)
		req.Type = &trimmed
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Code != nil {
		trimmed := strings.TrimSpace(*req.Code)
		req.Code = &trimmed
	}
	if req.Address != nil {
		trimmed := strings.TrimSpace(*req.Address)
		req.Address = &trimmed
	}
	if req.Phone != nil {
		trimmed := strings.TrimSpace(*req.Phone)
		req.Phone = &trimmed
	}
	if req.ImageURL != nil {
		trimmed := strings.TrimSpace(*req.ImageURL)
		req.ImageURL = &trimmed
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	updates := make(map[string]any, 10)
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.ProvinceID != nil {
		updates["province_id"] = *req.ProvinceID
	}
	if req.CityID != nil {
		updates["city_id"] = *req.CityID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
	}

	tx := ctl.db.Model(&models.Office{}).Where("id = ?", uint(id)).Updates(updates)
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "office not found")
	}

	office, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"office": office,
	})
}

func (ctl *OfficeController) Destroy(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := ctl.db.Delete(&models.Office{}, uint(id))
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "office not found")
	}

	return response.OK(c, "deleted", fiber.Map{
		"deleted": true,
	})
}

func (ctl *OfficeController) getByID(id uint) (models.Office, error) {
	var office models.Office
	if err := ctl.db.First(&office, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Office{}, fiber.NewError(fiber.StatusNotFound, "office not found")
		}
		return models.Office{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return office, nil
}
