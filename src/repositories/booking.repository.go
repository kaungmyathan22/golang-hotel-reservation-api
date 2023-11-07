package repository

import (
	"context"

	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BookingCollection = "bookings"

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBooking(ctx context.Context, filter bson.M) ([]*types.Booking, error)
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DB_NAME).Collection(BookingCollection),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	result, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = result.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBooking(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := s.collection.Find(ctx, filter)
	var bookings []*types.Booking
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
