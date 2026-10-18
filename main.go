package main

import (
	"log"
	"workshop-go/controller"
	"workshop-go/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: controller.ErrorHandler,
	})
	routes.Setup(app)
	log.Fatal(app.Listen(":3000"))
}
