package routes

import (
	c "workshop-go/controller" // ตั้งชื่อย่อ controller เป็น c

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func Setup(app *fiber.App) {
	// วางก่อน routes เพื่อให้ทุก API ต้องผ่าน Basic Auth
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"gofiber": "21022566",
		},
	}))

	v1 := app.Group("/api/v1")
	v1.Get("/fact/:number", c.Factorial)
	v1.Post("/register", c.Register)

	v3 := app.Group("/api/v3")
	v3.Get("/joe", c.TestParams) // เปลี่ยน joe เป็นชื่อเล่นของตัวเอง
}
