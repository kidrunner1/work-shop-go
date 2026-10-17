package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Factorial(ctx *fiber.Ctx) error {
	// แปลงข้อความจาก path เป็นจำนวนเต็ม
	number, err := strconv.Atoi(ctx.Params("number"))
	if err != nil || number < 0 || number > 20 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณาใส่จำนวนเต็มตั้งแต่ 0 ถึง 20",
		})
	}

	// int64 เก็บได้ถึง 20! และ 0! มีค่าเท่ากับ 1
	var result int64 = 1
	for i := 1; i <= number; i++ {
		result *= int64(i)
	}

	return ctx.JSON(fiber.Map{
		"number":    number,
		"factorial": result,
	})
}
