package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/stneto1/link-store/pkg"
)

func main() {
	pkg.PrepareDB()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", pkg.HandleIndex)
	app.Get("/links", pkg.HandleListLinks)
	app.Post("/create", pkg.HandleNewLink)

	app.Listen(":3000")
}
