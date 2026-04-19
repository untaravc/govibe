package authcontroller

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
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

	appCfg := configs.LoadAppConfig()

	now := time.Now()
	expiresAt := now.Add(time.Duration(appCfg.AccessTokenPeriod) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     u.ID,
		"email":   u.Email,
		"id":      u.ID,
		"name":    u.Name,
		"image":   u.Image,
		"role_id": u.RoleID,
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
	})

	signed, err := token.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := randomHexString(120)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	refreshTokenExpiresAt := now.Add(time.Duration(appCfg.RefreshTokenPeriod) * time.Minute)

	if err := ctl.db.
		Model(&models.User{}).
		Where("id = ?", u.ID).
		Where("deleted_at IS NULL").
		Updates(map[string]any{
			"refresh_token":            refreshToken,
			"refresh_token_expired_at": refreshTokenExpiresAt,
			"refresh_token_updated_at": now,
		}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	u.RefreshToken = &refreshToken
	u.RefreshTokenExp = &refreshTokenExpiresAt
	u.RefreshTokenUpd = &now

	return response.OK(c, "ok", fiber.Map{
		"access_token":  signed,
		"token_type":    "Bearer",
		"expires_at":    expiresAt.UTC().Format(time.RFC3339),
		"expires_unix":  expiresAt.Unix(),
		"refresh_token": refreshToken,
		"user":          u,
	})
}

func (ctl *AuthController) RefreshToken(c *fiber.Ctx) error {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
	}
	refreshToken := parts[1]

	now := time.Now()

	var u models.User
	if err := ctl.db.
		Where("refresh_token = ?", refreshToken).
		Where("deleted_at IS NULL").
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if u.RefreshTokenExp == nil || !u.RefreshTokenExp.After(now) {
		return fiber.NewError(fiber.StatusUnauthorized, "refresh token expired")
	}

	jwtCfg, err := configs.LoadJWTConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	appCfg := configs.LoadAppConfig()

	accessExpiresAt := now.Add(time.Duration(appCfg.AccessTokenPeriod) * time.Minute)
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     u.ID,
		"email":   u.Email,
		"id":      u.ID,
		"name":    u.Name,
		"image":   u.Image,
		"role_id": u.RoleID,
		"iat":     now.Unix(),
		"exp":     accessExpiresAt.Unix(),
	})
	signedAccess, err := access.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newRefreshToken, err := randomHexString(120)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	newRefreshTokenExpiresAt := now.Add(time.Duration(appCfg.RefreshTokenPeriod) * time.Minute)

	res := ctl.db.
		Model(&models.User{}).
		Where("id = ?", u.ID).
		Where("refresh_token = ?", refreshToken).
		Where("deleted_at IS NULL").
		Updates(map[string]any{
			"refresh_token":            newRefreshToken,
			"refresh_token_expired_at": newRefreshTokenExpiresAt,
			"refresh_token_updated_at": now,
		})
	if res.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	u.RefreshToken = &newRefreshToken
	u.RefreshTokenExp = &newRefreshTokenExpiresAt
	u.RefreshTokenUpd = &now

	return response.OK(c, "ok", fiber.Map{
		"access_token":  signedAccess,
		"token_type":    "Bearer",
		"expires_at":    accessExpiresAt.UTC().Format(time.RFC3339),
		"expires_unix":  accessExpiresAt.Unix(),
		"refresh_token": newRefreshToken,
		"user":          u,
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

	defaultRoleID := uint(1)
	u := models.User{
		Name:     req.Name,
		Email:    req.Email,
		RoleID:   &defaultRoleID,
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

	userID, ok := parseJWTSubToUint(claims["sub"])
	if !ok || userID == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token subject")
	}

	var u models.User
	if err := ctl.db.
		Select("id", "name", "email", "phone", "image").
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"user": fiber.Map{
			"name":  u.Name,
			"email": u.Email,
			"phone": u.Phone,
			"image": u.Image,
		},
	})
}

func (ctl *AuthController) UpdateProfile(c *fiber.Ctx) error {
	var raw updateProfileRequest
	if err := c.BodyParser(&raw); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	userIDAny := c.Locals("auth_user_id")
	userID, ok := parseJWTSubToUint(userIDAny)
	if !ok || userID == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
	}

	var req updateProfileRequest
	updates := map[string]any{}

	if raw.Name != nil {
		v := strings.TrimSpace(*raw.Name)
		req.Name = &v
		updates["name"] = v
	}

	if raw.Email != nil {
		v := strings.TrimSpace(*raw.Email)
		req.Email = &v
		updates["email"] = v
	}

	if raw.Phone != nil {
		v := strings.TrimSpace(*raw.Phone)
		if v == "" {
			req.Phone = nil
			updates["phone"] = nil
		} else {
			req.Phone = &v
			updates["phone"] = v
		}
	}

	if raw.Password != nil {
		v := strings.TrimSpace(*raw.Password)
		if v != "" {
			req.Password = &v
		}
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	if raw.Name == nil && raw.Email == nil && raw.Phone == nil && (raw.Password == nil || strings.TrimSpace(*raw.Password) == "") {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"fields": "no fields to update",
			},
		})
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"name": "cannot be empty",
			},
		})
	}
	if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"email": "cannot be empty",
			},
		})
	}

	if req.Email != nil {
		var existing models.User
		if err := ctl.db.
			Select("id").
			Where("email = ?", *req.Email).
			Where("id <> ?", userID).
			First(&existing).Error; err == nil {
			return fiber.NewError(fiber.StatusConflict, "email already registered")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		updates["password"] = string(hash)
	}

	if len(updates) == 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"fields": "no fields to update",
			},
		})
	}

	res := ctl.db.
		Model(&models.User{}).
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		Updates(updates)
	if res.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	var u models.User
	if err := ctl.db.
		Select("id", "name", "email", "phone", "image").
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "profile updated", fiber.Map{
		"user": fiber.Map{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"phone": u.Phone,
			"image": u.Image,
		},
	})
}

func (ctl *AuthController) Logout(c *fiber.Ctx) error {
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

	userID, ok := parseJWTSubToUint(claims["sub"])
	if !ok || userID == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token subject")
	}

	now := time.Now()
	res := ctl.db.
		Model(&models.User{}).
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		Updates(map[string]any{
			"refresh_token":            nil,
			"refresh_token_expired_at": nil,
			"refresh_token_updated_at": &now,
		})
	if res.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "user not found")
	}

	return response.OK(c, "ok", fiber.Map{
		"logged_out": true,
	})
}

func parseJWTSubToUint(v any) (uint, bool) {
	switch t := v.(type) {
	case float64:
		if t <= 0 {
			return 0, false
		}
		return uint(t), true
	case int:
		if t <= 0 {
			return 0, false
		}
		return uint(t), true
	case int64:
		if t <= 0 {
			return 0, false
		}
		return uint(t), true
	case uint:
		if t == 0 {
			return 0, false
		}
		return t, true
	case uint64:
		if t == 0 {
			return 0, false
		}
		return uint(t), true
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0, false
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || n == 0 {
			return 0, false
		}
		return uint(n), true
	default:
		return 0, false
	}
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

		appURL := strings.TrimRight(strings.TrimSpace(os.Getenv("APP_URL")), "/")
		link := fmt.Sprintf("%s/auth/validate-email-token?token=%s", appURL, url.QueryEscape(token))
		if appURL == "" {
			link = fmt.Sprintf("/auth/validate-email-token?token=%s", url.QueryEscape(token))
		}

		title := "GoVibe password reset"
		body := fmt.Sprintf("Reset your password using this link:\n\n%s\n", link)
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

func (ctl *AuthController) ValidateEmailToken(c *fiber.Ctx) error {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": fiber.Map{
				"token": "is required",
			},
		})
	}

	var u models.User
	if err := ctl.db.
		Select("id").
		Where("email_token = ?", token).
		Where("deleted_at IS NULL").
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "token not found", fiber.Map{
				"valid": false,
			})
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"valid": true,
	})
}

func (ctl *AuthController) UpdatePasswordWithToken(c *fiber.Ctx) error {
	var req updatePasswordWithTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.EmailToken = strings.TrimSpace(req.EmailToken)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{
			"errors": errs,
		})
	}

	var u models.User
	if err := ctl.db.
		Select("id").
		Where("email_token = ?", req.EmailToken).
		Where("deleted_at IS NULL").
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "token not found", fiber.Map{
				"updated": false,
			})
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := ctl.db.
		Model(&models.User{}).
		Where("id = ?", u.ID).
		Where("email_token = ?", req.EmailToken).
		Where("deleted_at IS NULL").
		Updates(map[string]any{
			"password":    string(hash),
			"email_token": nil,
		})
	if res.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return response.Error(c, fiber.StatusNotFound, "token not found", fiber.Map{
			"updated": false,
		})
	}

	return response.OK(c, "password updated", fiber.Map{
		"updated": true,
	})
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
