package repository

import (
	"context"

	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const RoomCollection = "room"

type RoomStore interface {
	GetRoomByID(ctx context.Context, id string) (*types.Room, error)
	GetRooms(ctx context.Context) ([]*types.Room, error)
	CreateRoom(ctx context.Context, Room *types.Room) (*types.Room, error)
	DeleteRoomByID(ctx context.Context, id string) error
	// UpdateRoom(context.Context, bson.M, types.UpdateRoomPayload) error
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DB_NAME).Collection(RoomCollection),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context) ([]*types.Room, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil

}

func (s *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*types.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var room types.Room
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = result.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}
func (s *MongoRoomStore) DeleteRoomByID(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var room types.Room
	err = s.collection.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&room)
	return err
}
