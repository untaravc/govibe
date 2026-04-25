package regioncontroller

import (
	"strconv"
	"strings"

	"govibe/app/Http/Response"
	"govibe/app/Models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RegionController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RegionController {
	return &RegionController{db: db}
}

func (ctl *RegionController) Provinces(c *fiber.Ctx) error {
	var provinces []models.Province

	q := ctl.db.Model(&models.Province{})
	if name := strings.TrimSpace(c.Query("name")); name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	if err := q.Order("name asc").Find(&provinces).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"provinces": provinces,
	})
}

func (ctl *RegionController) Cities(c *fiber.Ctx) error {
	var cities []models.City

	q := ctl.db.Model(&models.City{})

	if provinceIDRaw := strings.TrimSpace(c.Query("province_id")); provinceIDRaw != "" {
		n, err := strconv.ParseUint(provinceIDRaw, 10, 64)
		if err != nil || n == 0 {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid province_id")
		}
		q = q.Where("province_id = ?", n)
	}

	if name := strings.TrimSpace(c.Query("name")); name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	if err := q.Order("name asc").Find(&cities).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"cities": cities,
	})
}

func (ctl *RegionController) Districts(c *fiber.Ctx) error {
	var districts []models.District

	q := ctl.db.Model(&models.District{})

	if cityIDRaw := strings.TrimSpace(c.Query("city_id")); cityIDRaw != "" {
		n, err := strconv.ParseUint(cityIDRaw, 10, 64)
		if err != nil || n == 0 {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid city_id")
		}
		q = q.Where("city_id = ?", n)
	}

	if name := strings.TrimSpace(c.Query("name")); name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	if err := q.Order("name asc").Find(&districts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"districts": districts,
	})
}

func (ctl *RegionController) Villages(c *fiber.Ctx) error {
	var villages []models.Village

	q := ctl.db.Model(&models.Village{})

	if districtIDRaw := strings.TrimSpace(c.Query("district_id")); districtIDRaw != "" {
		n, err := strconv.ParseUint(districtIDRaw, 10, 64)
		if err != nil || n == 0 {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid district_id")
		}
		q = q.Where("district_id = ?", n)
	}

	if name := strings.TrimSpace(c.Query("name")); name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	if err := q.Order("name asc").Find(&villages).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"villages": villages,
	})
}

