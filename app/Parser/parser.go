package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UintParam(c *fiber.Ctx, name string) (uint64, error) {
	raw := strings.TrimSpace(c.Params(name))
	if raw == "" {
		return 0, errors.New("missing " + name)
	}
	return strconv.ParseUint(raw, 10, 64)
}

