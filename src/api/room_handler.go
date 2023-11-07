package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate    time.Time `json:"fromDate"`
	TillDate    time.Time `json:"tillDate"`
	NumOfGuests int       `json:"numOfGuests"`
}

type RoomHandler struct {
	store *repository.Store
}

func NewRoomHandler(store *repository.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (roomHandler *RoomHandler) HandleGetBookings(c *fiber.Ctx) error {
	where := bson.M{}
	bookings, err := roomHandler.store.Book.GetBooking(c.Context(),where)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (roomHandler *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomID, err := primitive.ObjectIDFromHex(c.Params("roomId"))
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
			"message": "User not found.",
		})
	}
	var payload BookRoomParams
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	fmt.Println(payload)
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   payload.FromDate,
		TillDate:   payload.TillDate,
		NumPersons: payload.NumOfGuests,
	}
	result, err := roomHandler.store.Book.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	fmt.Println(booking)
	return c.JSON(result)
}
