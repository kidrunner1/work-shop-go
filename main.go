package main

import (
	"log"
	"workshop-go/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.Setup(app)
	log.Fatal(app.Listen(":3000"))
}
