package usercontroller

import (
	"errors"
	"strings"

	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/app/Parser"
	appvalidator "govibe/app/Validator"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (ctl *UserController) Index(c *fiber.Ctx) error {
	var users []models.User
	if err := ctl.db.Order("id desc").Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return response.OK(c, "ok", fiber.Map{
		"users": users,
	})
}

func (ctl *UserController) Show(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}
	return response.OK(c, "ok", fiber.Map{
		"user": u,
	})
}

func (ctl *UserController) Store(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	u := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}

	if err := ctl.db.Create(&u).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "created", fiber.Map{
		"user": u,
	})
}

func (ctl *UserController) Update(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req updateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Email != nil {
		trimmed := strings.TrimSpace(*req.Email)
		req.Email = &trimmed
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	updates := make(map[string]any, 3)

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if req.Password != nil {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, hashErr.Error())
		}
		updates["password"] = string(hash)
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
	}

	tx := ctl.db.Model(&models.User{}).Where("id = ?", uint(id)).Updates(updates)
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	u, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}
	return response.OK(c, "ok", fiber.Map{
		"user": u,
	})
}

func (ctl *UserController) Destroy(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := ctl.db.Delete(&models.User{}, uint(id))
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return response.OK(c, "deleted", fiber.Map{
		"deleted": true,
	})
}

func (ctl *UserController) getByID(id uint) (models.User, error) {
	var u models.User
	if err := ctl.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return models.User{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return u, nil
}

// param parsing lives in app/Parser
