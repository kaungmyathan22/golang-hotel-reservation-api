package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/api"
)

func main() {
	var PORT = flag.String("port", ":5000", "Listen address of the api server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	app.Get("/foo", handleFoofunc)
	apiv1.Post("/user", api.HandleCreateUsers)
	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*PORT)
}

func handleFoofunc(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Hello from golang"})
}
