package api

import (
	"github.com/gofiber/fiber/v2"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	hotelStore repository.HotelStore
	roomStore  repository.RoomStore
}

type HotelQueryParams struct {
	Rooms bool
}

func NewHotelHandle(hs repository.HotelStore, rs repository.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelId": hexId}
	rooms, err := h.roomStore.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	hotels, err := h.hotelStore.GetHotels(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
