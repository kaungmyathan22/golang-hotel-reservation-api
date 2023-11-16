package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *repository.Store
}

func NewBookingHandler(store *repository.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (bookingHandler *BookingHandler) GetBookings(c *fiber.Ctx) error {
	bookings, err := bookingHandler.store.Book.GetBooking(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (bookingHandler *BookingHandler) GetMyBookings(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		c.Status(http.StatusUnauthorized).JSON(types.GenericResponse{Message: "Unauthorized access need to login firrst."})
	}
	bookings, err := bookingHandler.store.Book.GetBooking(c.Context(), bson.M{"userID": user.ID})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}
