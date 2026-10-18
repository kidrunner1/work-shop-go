package controller

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler รับ error ที่ controller หรือ middleware ส่งกลับมา
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "เกิดข้อผิดพลาดภายในระบบ กรุณาลองใหม่อีกครั้ง"

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		status = fiberError.Code
		if status < 500 {
			message = fiberError.Message
		}
	}

	if status >= 500 {
		// เก็บรายละเอียดไว้ที่ server ไม่ส่งข้อมูลภายในกลับให้ผู้ใช้
		log.Printf("request error: %v", err)
	}
	if status == fiber.StatusNotFound {
		message = "ไม่พบ API ที่เรียก กรุณาตรวจสอบ URL และ HTTP method"
	}

	return ctx.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
