package routes

import (
	"api/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)


func GetURL(c *fiber.Ctx) error {
	url := c.Params("url")

	rdb := database.CreateClient(0)
	defer rdb.Close()

	value, err := rdb.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot resolve short url"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"shortened_url": url,
		"original_url": value,
	})
}
