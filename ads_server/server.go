package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"../adspb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
}

type adsItem struct {
	ID        *primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
	UserID    string              `bson:"user_id"`
	Title     string              `bson:"title"`
	Content   string              `bson:"content"`
	Address   string              `bson:"address"`
	Email     string              `bson:"email"`
	Phone     string              `bson:"phone"`
	BannerURL string              `bson:"banner_url"`
}

func main() {

	//we can get the file name and line number problem error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Welcome to Server Ads")

	//connect to mongodb
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("mydb").Collection("ads")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	adspb.RegisterAdsServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Stopping server")
	s.Stop()

	fmt.Println("Closing listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.TODO())
	fmt.Println("End Program")
}
