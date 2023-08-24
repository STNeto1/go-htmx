package pkg

import (
	"github.com/gofiber/fiber/v2"
)

type Link struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Link        string  `json:"link"`
	Thumbnail   *string `json:"thumbnail"`
}

func HandleIndex(c *fiber.Ctx) error {

	links, _ := DB.GetLinks()

	return c.Render("index", fiber.Map{
		"links": links,
		"pool":  true,
	}, "layouts/main")
}

func HandleListLinks(c *fiber.Ctx) error {

	links, _ := DB.GetLinks()

	shouldPool := false
	for _, link := range links {
		if link.Thumbnail == nil {
			shouldPool = true
		}
	}

	return c.Render("partials/list", fiber.Map{
		"links": links,
		"pool":  shouldPool,
	})
}

type CreateLinkRequest struct {
	Title       string  `form:"title"`
	Description *string `form:"description"`
	Link        string  `form:"link"`
}

func HandleNewLink(c *fiber.Ctx) error {
	var req CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		req = CreateLinkRequest{
			Title:       "failed",
			Description: nil,
			Link:        "failed",
		}
	}

	_ = DB.CreateLink(req)
	// if err != nil {
	// 	log.Println(err)
	// }

	links, _ := DB.GetLinks()

	return c.Render("index", fiber.Map{
		"links": links,
		"pool":  true,
	})
}
