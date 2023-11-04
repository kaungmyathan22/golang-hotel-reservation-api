package main

import (
	"context"
	"fmt"
	"log"

	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("============================== Seeding Start ==============================")
	ctx := context.Background()
	opts := options.Client().ApplyURI(repository.DB_URI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	client.Database(repository.DB_NAME).Drop(ctx)
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	hotelStore := repository.NewMongoHotelStore(client)
	roomStore := repository.NewMongoRoomStore(client, hotelStore)
	hotel := types.Hotel{
		Name:     "Belluia",
		Location: "France",
		Rooms:    []primitive.ObjectID{},
		Rating:   3,
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
		{
			Size:  "king",
			Price: 399.9,
		},
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
	fmt.Println(insertedHotel)
	fmt.Println("============================== Seeding Done ==============================")
}
