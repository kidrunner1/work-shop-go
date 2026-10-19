package routes

import (
	g "workshop-go/controller" // ตั้งชื่อย่อ controller เป็น c

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App) {
	// ดัก panic แล้วส่งต่อให้ ErrorHandler ตอบกลับเป็น HTTP 500
	app.Use(recover.New())

	// วางก่อน routes เพื่อให้ทุก API ต้องผ่าน Basic Auth
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"gofiber": "21022566",
		},
		Unauthorized: func(ctx *fiber.Ctx) error {
			ctx.Set("WWW-Authenticate", `Basic realm="Workshop"`)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "ยืนยันตัวตนไม่สำเร็จ กรุณาส่ง Basic Auth ด้วย username และ password ที่ถูกต้อง",
			})
		},
	}))

	v1 := app.Group("/api/v1")
	v1.Get("/fact/:number", g.Factorial)
	v1.Post("/register", g.Register)

	v3 := app.Group("/api/v3")
	v3.Get("/krit", g.TestParams) // เปลี่ยน krit เป็นชื่อเล่นของตัวเอง
}
