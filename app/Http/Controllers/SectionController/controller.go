package sectioncontroller

import (
	"govibe/app/Http/Response"

	"github.com/gofiber/fiber/v2"
)

type SectionController struct{}

func New() *SectionController {
	return &SectionController{}
}

func (ctl *SectionController) Index(c *fiber.Ctx) error {
	return response.OK(c, "ok", fiber.Map{
		"sections": []string{"post", "office"},
	})
}
