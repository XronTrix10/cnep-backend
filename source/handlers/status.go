package handlers

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"fmt"
)

var startTime time.Time

func StatusInit() {
	startTime = time.Now()
}

func Status() fiber.Handler {
	return func(c *fiber.Ctx) error {

		elapsedTime := time.Since(startTime)
		days := int64(elapsedTime.Hours() / 24)
		hours := int64(int64(elapsedTime.Hours()) % 24)
		minutes := int64(int64(elapsedTime.Minutes()) % 60)
		seconds := int64(int64(elapsedTime.Seconds()) % 60)

		var timeString string
		if days > 0 {
			timeString += fmt.Sprintf("%d days, ", days)
		}
		if hours > 0 || days > 0 {
			timeString += fmt.Sprintf("%d hours, ", hours)
		}
		if minutes > 0 || hours > 0 || days > 0 {
			timeString += fmt.Sprintf("%d minutes, ", minutes)
		}
		timeString += fmt.Sprintf("%d seconds", seconds)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Server is up and running",
			"uptime":  timeString,
		})
	}
}
