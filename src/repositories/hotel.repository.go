package repository

import (
	"context"

	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const HotelCollection = "hotel"

type HotelStore interface {
	GetHotelByID(ctx context.Context, id string) (*types.Hotel, error)
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	CreateHotel(ctx context.Context, Hotel *types.Hotel) (*types.Hotel, error)
	DeleteHotelByID(ctx context.Context, id string) error
	UpdateHotel(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(DB_NAME).Collection(HotelCollection),
	}
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil

}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	result, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = result.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
func (s *MongoHotelStore) DeleteHotelByID(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var hotel types.Hotel
	err = s.collection.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&hotel)
	return err
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}
