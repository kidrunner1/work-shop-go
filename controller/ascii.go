package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TestParams(ctx *fiber.Ctx) error {
	// รับค่า Query Parameter ชื่อ tax_id
	// ตัวอย่าง: /api/v3/krit?tax_id=Hello
	taxID := ctx.Query("tax_id")

	// ถ้าไม่ได้ส่ง tax_id มา หรือส่งมาเป็นค่าว่าง
	// ให้ตอบกลับเป็น 400 Bad Request
	if taxID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณาระบุ tax_id ใน query เช่น ?tax_id=Hello",
		})
	}

	// สร้าง slice เอาไว้เก็บค่า ASCII ของแต่ละตัวอักษร
	// เช่น Hello จะเก็บเป็น ["72", "101", "108", "108", "111"]
	values := []string{}

	// วนอ่านตัวอักษรใน taxID ทีละตัว
	for _, character := range taxID {

		// ASCII มีค่าตั้งแต่ 0 ถึง 127
		// ถ้าค่ามากกว่า 127 เช่น ภาษาไทย ให้ตอบ error
		if character > 127 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "tax_id ต้องเป็นอักขระ ASCII เช่น ตัวอักษรภาษาอังกฤษหรือตัวเลข",
			})
		}

		// แปลงตัวอักษรเป็นตัวเลข ASCII
		// strconv.Itoa แปลงเลข 72 ให้เป็น string "72"
		// จากนั้น append เก็บลงใน values
		values = append(values, strconv.Itoa(int(character)))
	}

	// ส่งผลลัพธ์กลับเป็น JSON
	// เช่น ["72", "101", "108"] -> "72 101 108"
	return ctx.JSON(fiber.Map{
		"tax_id": taxID,
		"ascii":  strings.Join(values, " "),
	})
}