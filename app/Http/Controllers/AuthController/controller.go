package authcontroller

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/app/Service"
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
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
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

	return response.OK(c, "ok", fiber.Map{
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
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
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

	return response.Created(c, "created", fiber.Map{
		"user": u,
	})
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

	return response.OK(c, "ok", fiber.Map{
		"claims": claims,
	})
}

func (ctl *AuthController) RequestResetPassword(c *fiber.Ctx) error {
	var req requestResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	if req.Email == "" && req.Phone == "" {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"email": "is required",
				"phone": "is required",
			},
		})
	}
	if req.Email != "" && req.Phone != "" {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"email": "provide either email or phone",
				"phone": "provide either email or phone",
			},
		})
	}

	// Avoid user enumeration: respond OK even when the account doesn't exist.
	if req.Email != "" {
		var u models.User
		if err := ctl.db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.OK(c, "if the account exists, a reset token has been sent", fiber.Map{})
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		token, err := randomHexString(60)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err := ctl.db.Model(&models.User{}).Where("id = ?", u.ID).Update("email_token", token).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		title := "GoVibe password reset"
		body := fmt.Sprintf("Your reset token:\n\n%s\n", token)
		if err := service.SendEmail(req.Email, title, body); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return response.OK(c, "reset token sent", fiber.Map{})
	}

	var u models.User
	if err := ctl.db.Where("phone = ?", req.Phone).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.OK(c, "if the account exists, a reset token has been sent", fiber.Map{})
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	code, err := randomNumericString(6)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := ctl.db.Model(&models.User{}).Where("id = ?", u.ID).Update("phone_token", code).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := service.SendWhatsAppMessage(req.Phone, fmt.Sprintf("Your reset code:\n\n%s\n", code)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "reset token sent", fiber.Map{})
}

func randomHexString(length int) (string, error) {
	if length <= 0 || length%2 != 0 {
		return "", errors.New("length must be a positive even number")
	}
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func randomNumericString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil) // 10^length
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", length, n.Int64()), nil
}
