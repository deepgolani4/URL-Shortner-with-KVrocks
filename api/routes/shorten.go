package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL					string		     
	CustomShortURL		string			 
	Expiry				time.Duration	 
}

type response struct {
	URL    				string			 
	CustomShortURL      string			
	Expiry              time.Duration	 
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(body); err != nil {
		fmt.Println("Error parsing body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if !govalidator.IsURL(body.URL) {	
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid url"})
	}

	var id string
	if body.CustomShortURL == "" {
		uuid7, err := uuid.NewV7()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error generating ID"})
		}
		
		uuidStr := uuid7.String()
		uuidWithoutHyphens := ""
		for _, char := range uuidStr {
			if char != '-' {
				uuidWithoutHyphens += string(char)
			}
		}
		
		id = uuidWithoutHyphens[:12]
	} else {
		id = body.CustomShortURL
	}	

	r := database.CreateClient(0)
	defer r.Close()

	val , _ := r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "custom url already exists"})
	}

	if body.Expiry == 0 {
		body.Expiry = 1 * time.Hour
	}

	var setErr error
	setErr = r.Set(database.Ctx, id, body.URL, body.Expiry).Err()
	if setErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error creating custom url"})
	}

	resp := response{
		URL:             body.URL,
		CustomShortURL:  os.Getenv("DOMAIN") + "/" + id,
		Expiry: 		 body.Expiry,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
