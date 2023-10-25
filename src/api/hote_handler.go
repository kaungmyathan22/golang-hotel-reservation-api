package api

import (
	"github.com/gofiber/fiber/v2"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore repository.HotelStore
	roomStore  repository.RoomStore
}

func NewHotelHandle(hs repository.HotelStore, rs repository.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
