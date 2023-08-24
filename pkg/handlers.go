package pkg

import (
	"log"

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

	rows, err := DB.Query("SELECT * FROM links")
	if err != nil {
		return c.Render("index", fiber.Map{
			"error": err.Error(),
		})
	}

	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Description, &link.Link, &link.Thumbnail)
		if err != nil {
			log.Println(err)
			continue
		}

		links = append(links, link)
	}

	return c.Render("index", fiber.Map{
		"links": links,
		"pool":  true,
	}, "layouts/main")
}

func HandleListLinks(c *fiber.Ctx) error {

	rows, err := DB.Query("SELECT * FROM links")
	if err != nil {
		return c.Render("index", fiber.Map{
			"error": err.Error(),
		})
	}

	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Description, &link.Link, &link.Thumbnail)
		if err != nil {
			log.Println(err)
			continue
		}

		links = append(links, link)
	}

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
		log.Println(err)
		req = CreateLinkRequest{
			Title:       "failed",
			Description: nil,
			Link:        "failed",
		}
	}

	_, err := DB.Exec("INSERT INTO links (title, description, link) VALUES (?, ?, ?)", req.Title, req.Description, req.Link)
	if err != nil {
		log.Println(err)
	}

	return c.Render("partials/form", fiber.Map{})
}
