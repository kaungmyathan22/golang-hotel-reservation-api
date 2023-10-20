package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Kaung Myat",
		LastName:  "Han",
		ID:        "",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Hola"})
}

func HandleCreateUsers(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Hola"})
}
