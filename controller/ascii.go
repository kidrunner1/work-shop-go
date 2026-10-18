package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TestParams(ctx *fiber.Ctx) error {
	taxID := ctx.Query("tax_id")
	if taxID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณาระบุ tax_id ใน query เช่น ?tax_id=Hello",
		})
	}

	values := []string{}
	for _, character := range taxID {
		// ASCII มีค่าตั้งแต่ 0 ถึง 127 ภาษาไทยจึงอยู่นอกช่วงนี้
		if character > 127 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "tax_id ต้องเป็นอักขระ ASCII เช่น ตัวอักษรภาษาอังกฤษหรือตัวเลข",
			})
		}
		values = append(values, strconv.Itoa(int(character)))
	}

	return ctx.JSON(fiber.Map{
		"tax_id": taxID,
		"ascii":  strings.Join(values, " "),
	})
}
