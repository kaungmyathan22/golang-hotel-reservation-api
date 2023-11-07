package repository

const (
	DB_NAME      = "hotel-reservation"
	TEST_DB_NAME = "test-hotel-reservation"
	DB_URI       = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
	Book  BookingStore
}
