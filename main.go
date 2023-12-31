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

	fmt.Println("Pinged. You successfully connected to MongoDB!")

	userStore := repository.NewMongoUserStore(client)
	hotelStore := repository.NewMongoHotelStore(client)
	roomStore := repository.NewMongoRoomStore(client, hotelStore)
	bookingStore := repository.NewMongoBookingStore(client)

	store := &repository.Store{
		User:  userStore,
		Hotel: hotelStore,
		Room:  roomStore,
		Book:  bookingStore,
	}

	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandle(hotelStore, roomStore)
	authHandler := api.NewAuthHandler(userStore)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middlewares.JWTAuthentication(userStore))

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
	//#endregion

	//#region ----- hotel routes
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:hotelId/rooms/:roomId", hotelHandler.HandleGetRoomById)
	//#endregion

	//#region------------------- booking routes
	apiv1.Post("/hotel/:hotelId/rooms/:roomId/book", roomHandler.HandleBookRoom)
	apiv1.Get("/bookings", bookingHandler.GetBookings)
	apiv1.Get("/bookings/me", bookingHandler.GetMyBookings)
	//#endregion
	app.Listen(*PORT)
}
