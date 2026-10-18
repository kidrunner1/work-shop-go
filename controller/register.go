package controller

import (
	"errors"
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
			"message": "ใส่ข้อมูลผิดพลาด: กรุณาส่ง body เป็น JSON และให้ทุกฟิลด์เป็นข้อความ",
		})
	}

	// 2. ตรวจสอบทุกฟิลด์ตาม validate tags ใน struct
	if err := validate.Struct(input); err != nil {
		var fieldErrors validator.ValidationErrors
		if !errors.As(err, &fieldErrors) {
			return err // ข้อผิดพลาดภายใน ส่งต่อให้ ErrorHandler จัดการเป็น 500
		}

		// เก็บข้อความแยกตามชื่อฟิลด์ เพื่อให้รู้ว่าต้องแก้ช่องไหน
		details := map[string]string{}
		for _, fieldError := range fieldErrors {
			switch fieldError.StructField() {
			case "Username":
				details["username"] = "กรุณากรอกชื่อผู้ใช้ 3–20 ตัว โดยใช้ตัวอักษรภาษาอังกฤษหรือตัวเลขเท่านั้น"
			case "Email":
				details["email"] = "กรุณากรอกอีเมลให้ถูกต้อง เช่น student@example.com และยาวไม่เกิน 254 ตัวอักษร"
			case "Password":
				details["password"] = "กรุณากรอกรหัสผ่าน 8–72 ตัวอักษร"
			case "ConfirmPassword":
				details["confirm_password"] = "กรุณากรอกยืนยันรหัสผ่านให้ตรงกับรหัสผ่าน"
			}
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ใส่ข้อมูลผิดพลาด",
			"errors":  details,
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
