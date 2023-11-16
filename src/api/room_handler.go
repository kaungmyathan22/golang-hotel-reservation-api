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
	bookings, err := roomHandler.store.Book.GetBooking(c.Context(), where)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (roomHandler *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	rawRoomID := c.Params("roomId")
	roomID, err := primitive.ObjectIDFromHex(rawRoomID)
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

	if err != nil {
		return err
	}
	ok, err = roomHandler.isRoomAvailableForBooking(c, roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"type":    "error",
			"message": fmt.Sprintf("room %s already booked.", rawRoomID),
		})
	}
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

func (roomHandler *RoomHandler) isRoomAvailableForBooking(c *fiber.Ctx, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"_id": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := roomHandler.store.Book.GetBooking(c.Context(), where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) > 0
	return ok, nil
}
