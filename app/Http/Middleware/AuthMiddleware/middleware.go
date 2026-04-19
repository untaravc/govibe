package authmiddleware

import (
	"errors"
	"strconv"
	"strings"

	"govibe/app/Models"
	"govibe/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	LocalUserID    = "auth_user_id"
	LocalUserEmail = "auth_user_email"
)

func New(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Allow CORS preflight and public auth endpoints to pass through.
		if c.Method() == fiber.MethodOptions || isPublicPath(c.Path()) {
			return c.Next()
		}

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

		c.Locals(LocalUserID, userID)

		if email, ok := claims["email"].(string); ok {
			email = strings.TrimSpace(email)
			if email != "" {
				c.Locals(LocalUserEmail, email)
			}
		}

		// Optional: when a DB handle is provided, ensure the user still exists.
		if db != nil {
			var u models.User
			if err := db.
				Select("id", "email").
				Where("id = ?", userID).
				Where("deleted_at IS NULL").
				First(&u).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fiber.NewError(fiber.StatusUnauthorized, "user not found")
				}
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			c.Locals(LocalUserEmail, u.Email)
		}

		return c.Next()
	}
}

func isPublicPath(path string) bool {
	switch strings.TrimSpace(path) {
	case "/api/health",
		"/api/register",
		"/api/login",
		"/api/refresh-token",
		"/api/request-reset-password":
		return true
	default:
		return false
	}
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
