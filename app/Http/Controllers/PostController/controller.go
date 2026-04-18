package postcontroller

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

type PostController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *PostController {
	return &PostController{db: db}
}

func (ctl *PostController) Index(c *fiber.Ctx) error {
	var posts []models.Post
	if err := ctl.db.Order("id desc").Find(&posts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return response.OK(c, "ok", fiber.Map{
		"posts": posts,
	})
}

func (ctl *PostController) Show(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	p, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}
	return response.OK(c, "ok", fiber.Map{
		"post": p,
	})
}

func (ctl *PostController) Store(c *fiber.Ctx) error {
	var req createPostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Slug = strings.TrimSpace(req.Slug)
	if req.Subtitle != nil {
		trimmed := strings.TrimSpace(*req.Subtitle)
		req.Subtitle = &trimmed
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

	p := models.Post{
		Title:      req.Title,
		Subtitle:   req.Subtitle,
		Slug:       req.Slug,
		Content:    req.Content,
		Status:     req.Status,
		ImageURL:   req.ImageURL,
		UserID:     req.UserID,
		ReleaseAt:  req.ReleaseAt,
		CategoryID: req.CategoryID,
	}

	if err := ctl.db.Create(&p).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "created", fiber.Map{
		"post": p,
	})
}

func (ctl *PostController) Update(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req updatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		req.Title = &trimmed
	}
	if req.Subtitle != nil {
		trimmed := strings.TrimSpace(*req.Subtitle)
		req.Subtitle = &trimmed
	}
	if req.Slug != nil {
		trimmed := strings.TrimSpace(*req.Slug)
		req.Slug = &trimmed
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
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Subtitle != nil {
		updates["subtitle"] = *req.Subtitle
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.UserID != nil {
		updates["user_id"] = *req.UserID
	}
	if req.ReleaseAt != nil {
		updates["release_at"] = *req.ReleaseAt
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}

	if len(updates) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
	}

	tx := ctl.db.Model(&models.Post{}).Where("id = ?", uint(id)).Updates(updates)
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusBadRequest, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "post not found")
	}

	p, err := ctl.getByID(uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, "ok", fiber.Map{
		"post": p,
	})
}

func (ctl *PostController) Destroy(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := ctl.db.Delete(&models.Post{}, uint(id))
	if tx.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "post not found")
	}

	return response.OK(c, "deleted", fiber.Map{
		"deleted": true,
	})
}

func (ctl *PostController) getByID(id uint) (models.Post, error) {
	var p models.Post
	if err := ctl.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Post{}, fiber.NewError(fiber.StatusNotFound, "post not found")
		}
		return models.Post{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return p, nil
}
