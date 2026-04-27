package postcontroller

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	authmiddleware "govibe/app/Http/Middleware/AuthMiddleware"
	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/app/Parser"
	appvalidator "govibe/app/Validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	title := strings.TrimSpace(c.Query("title"))
	q := ctl.db.Model(&models.Post{})
	if title != "" {
		q = q.Where("title LIKE ?", "%"+title+"%")
	}
	if err := q.Order("id desc").Find(&posts).Error; err != nil {
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

	userID, err := authenticatedUserID(c)
	if err != nil {
		return err
	}

	p := models.Post{
		Title:      req.Title,
		Subtitle:   req.Subtitle,
		Slug:       "pending-" + uuid.NewString(),
		Content:    req.Content,
		Status:     req.Status,
		ImageURL:   req.ImageURL,
		UserID:     userID,
		ReleaseAt:  req.ReleaseAt,
		CategoryID: req.CategoryID,
	}

	if err := ctl.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&p).Error; err != nil {
			return err
		}

		p.Slug = postSlug(p.ID, p.Title)
		return tx.Model(&p).Update("slug", p.Slug).Error
	}); err != nil {
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
		updates["slug"] = postSlug(uint(id), *req.Title)
	}
	if req.Subtitle != nil {
		updates["subtitle"] = *req.Subtitle
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

func authenticatedUserID(c *fiber.Ctx) (uint, error) {
	userID, ok := c.Locals(authmiddleware.LocalUserID).(uint)
	if !ok || userID == 0 {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "authenticated user not found")
	}
	return userID, nil
}

func postSlug(id uint, title string) string {
	suffix := slugifyTitle(title)
	if suffix == "" {
		suffix = "post"
	}
	return fmt.Sprintf("%d-%s", id, suffix)
}

func slugifyTitle(title string) string {
	var b strings.Builder
	lastDash := false

	for _, r := range strings.TrimSpace(title) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}

	return strings.Trim(b.String(), "-")
}
