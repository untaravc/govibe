package authcontroller

import (
	"errors"
	"strings"
	"time"

	"govibe/app/Models"
	appvalidator "govibe/app/Validator"
	"govibe/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

func (ctl *AuthController) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "validation error",
			"errors":  errs,
		})
	}

	var u models.User
	if err := ctl.db.Where("email = ?", req.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	jwtCfg, err := configs.LoadJWTConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(jwtCfg.TTLMin) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"iat":   now.Unix(),
		"exp":   expiresAt.Unix(),
	})

	signed, err := token.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"token":        signed,
		"token_type":   "Bearer",
		"expires_at":   expiresAt.UTC().Format(time.RFC3339),
		"expires_unix": expiresAt.Unix(),
		"user":         u,
	})
}

func (ctl *AuthController) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "validation error",
			"errors":  errs,
		})
	}

	var existing models.User
	if err := ctl.db.Select("id").Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return fiber.NewError(fiber.StatusConflict, "email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

	return c.Status(fiber.StatusCreated).JSON(u)
}

func (ctl *AuthController) Profile(c *fiber.Ctx) error {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
	}
	rawToken := parts[1]

	jwtCfg, err := configs.LoadJWTConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		func(t *jwt.Token) (any, error) { return []byte(jwtCfg.Secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || token == nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	return c.JSON(fiber.Map{
		"claims": claims,
	})
}
