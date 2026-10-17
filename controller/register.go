package controller

import (
	"workshop-go/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// สร้าง validator ครั้งเดียว แล้วใช้ตรวจสอบข้อมูลในแต่ละ request
var validate = validator.New(validator.WithRequiredStructEnabled())

func Register(ctx *fiber.Ctx) error {
	var input models.RegisterRequest

	// 1. อ่าน JSON จาก body มาใส่ struct
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ใส่ข้อมูลผิดพลาด",
		})
	}

	// 2. ตรวจสอบทุกฟิลด์ตาม validate tags ใน struct
	if err := validate.Struct(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ใส่ข้อมูลผิดพลาด",
		})
	}

	// 3. ตัวอย่างนี้ตรวจสอบข้อมูลเท่านั้น ยังไม่ได้บันทึกลงฐานข้อมูล
	// ส่งกลับเฉพาะ username และ email ไม่ส่งรหัสผ่านกลับ
	return ctx.JSON(fiber.Map{
		"message": "ข้อมูลสมัครสมาชิกถูกต้อง",
		"data": fiber.Map{
			"username": input.Username,
			"email":    input.Email,
		},
	})
}
