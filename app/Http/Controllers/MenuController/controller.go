package menucontroller

import (
	"errors"
	"strconv"
	"strings"

	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type MenuController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *MenuController {
	return &MenuController{db: db}
}

type MenuNode struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Icon     *string     `json:"icon,omitempty"`
	Slug     string      `json:"slug"`
	ParentID *uint       `json:"parent_id,omitempty"`
	Order    int         `json:"order"`
	Link     *string     `json:"link,omitempty"`
	Status   uint8       `json:"status"`
	Children []*MenuNode `json:"children"`
}

func buildMenuTree(rows []models.Menu) []*MenuNode {
	nodes := make(map[uint]*MenuNode, len(rows))
	roots := make([]*MenuNode, 0, len(rows))

	for _, m := range rows {
		nodes[m.ID] = &MenuNode{
			ID:       m.ID,
			Name:     m.Name,
			Icon:     m.Icon,
			Slug:     m.Slug,
			ParentID: m.ParentID,
			Order:    m.Order,
			Link:     m.Link,
			Status:   m.Status,
			Children: []*MenuNode{},
		}
	}

	for _, m := range rows {
		n := nodes[m.ID]
		if n == nil {
			continue
		}

		if m.ParentID == nil || *m.ParentID == 0 || *m.ParentID == m.ID {
			roots = append(roots, n)
			continue
		}

		parent := nodes[*m.ParentID]
		if parent == nil {
			roots = append(roots, n)
			continue
		}
		parent.Children = append(parent.Children, n)
	}

	return roots
}

// List returns all menus as a tree structure (not role-filtered).
func (ctl *MenuController) List(c *fiber.Ctx) error {
	var rows []models.Menu
	if err := ctl.db.
		Where("deleted_at IS NULL").
		Order("parent_id asc").
		Order("`order` asc").
		Order("id asc").
		Find(&rows).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"menus": buildMenuTree(rows),
	})
}

func (ctl *MenuController) Index(c *fiber.Ctx) error {
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
		Select("id", "role_id").
		Where("id = ?", userID).
		Where("deleted_at IS NULL").
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if u.RoleID == nil || *u.RoleID == 0 {
		return response.OK(c, "ok", fiber.Map{
			"menus": []*MenuNode{},
		})
	}

	var allowedIDs []uint
	if err := ctl.db.
		Model(&models.MenuRole{}).
		Distinct("menu_id").
		Where("role_id = ?", *u.RoleID).
		Where("deleted_at IS NULL").
		Where("method IN ?", []string{"get", "show"}).
		Pluck("menu_id", &allowedIDs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if len(allowedIDs) == 0 {
		return response.OK(c, "ok", fiber.Map{
			"menus": []*MenuNode{},
		})
	}

	var rows []models.Menu
	if err := ctl.db.
		Where("status = ?", 1).
		Where("deleted_at IS NULL").
		Order("parent_id asc").
		Order("`order` asc").
		Order("id asc").
		Find(&rows).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	allowedSet := make(map[uint]struct{}, len(allowedIDs))
	for _, id := range allowedIDs {
		allowedSet[id] = struct{}{}
	}

	parentByID := make(map[uint]*uint, len(rows))
	for _, m := range rows {
		parentByID[m.ID] = m.ParentID
	}

	// Add ancestors of allowed menus so tree structure stays intact.
	for _, id := range allowedIDs {
		cur := id
		for i := 0; i < 32; i++ { // safety bound for bad/cyclic data
			p := parentByID[cur]
			if p == nil || *p == 0 {
				break
			}
			if _, ok := allowedSet[*p]; ok {
				cur = *p
				continue
			}
			allowedSet[*p] = struct{}{}
			cur = *p
		}
	}

	nodes := make(map[uint]*MenuNode, len(rows))
	roots := make([]*MenuNode, 0, len(rows))

	for _, m := range rows {
		if _, ok := allowedSet[m.ID]; !ok {
			continue
		}
		nodes[m.ID] = &MenuNode{
			ID:       m.ID,
			Name:     m.Name,
			Icon:     m.Icon,
			Slug:     m.Slug,
			ParentID: m.ParentID,
			Order:    m.Order,
			Link:     m.Link,
			Status:   m.Status,
			Children: []*MenuNode{},
		}
	}

	for _, m := range rows {
		n := nodes[m.ID]
		if n == nil {
			continue
		}

		if m.ParentID == nil || *m.ParentID == 0 {
			roots = append(roots, n)
			continue
		}

		parent := nodes[*m.ParentID]
		if parent == nil {
			roots = append(roots, n)
			continue
		}
		parent.Children = append(parent.Children, n)
	}

	return response.OK(c, "ok", fiber.Map{
		"menus": roots,
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
