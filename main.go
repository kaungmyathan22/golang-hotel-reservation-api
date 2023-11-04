package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/api"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/middlewares"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var PORT = flag.String("port", ":5000", "Listen address of the api server")
	flag.Parse()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := err.Error()
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				message = e.Message
			}
			err = ctx.Status(code).JSON(map[string]any{"code": code, "message": message})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
			return nil
		},
	})

	opts := options.Client().ApplyURI(repository.DB_URI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	mongoUserStore := repository.NewMongoUserStore(client)
	hotelStore := repository.NewMongoHotelStore(client)
	roomStore := repository.NewMongoRoomStore(client, hotelStore)
	userHandler := api.NewUserHandler(mongoUserStore)
	hotelHandler := api.NewHotelHandle(hotelStore, roomStore)
	authHandler := api.NewAuthHandler(mongoUserStore)

	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middlewares.JWTAuthentication)

	//#region------------------- auth
	auth.Post("/authentication/login", authHandler.HandleLogin)
	auth.Post("/authentication/register", userHandler.HandleCreateUsers)
	//#endregion

	//#region ----- user routes
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandleCreateUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Patch("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUsers)

	//#region ----- hotel routes
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*PORT)
}
