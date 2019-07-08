package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/naufalziyad/ads-gRPC/adspb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

func (*server) CreateAds(ctx context.Context, req *adspb.CreateAdsRequest) (*adspb.CreateAdsResponse, error) {

	ads := req.GetAds()

	data := adsItem{
		UserID:    ads.GetUserId(),
		Title:     ads.GetTitle(),
		Content:   ads.GetContent(),
		Address:   ads.GetAddress(),
		Email:     ads.GetEmail(),
		Phone:     ads.GetPhone(),
		BannerURL: ads.GetBannerUrl(),
	}

	//this is for connect to mongodb
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal; error: %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}
	fmt.Println("Create Ads Request Success at :", time.Now())
	return &adspb.CreateAdsResponse{
		Ads: &adspb.Ads{
			Id:        oid.Hex(),
			UserId:    ads.GetUserId(),
			Title:     ads.GetTitle(),
			Content:   ads.GetContent(),
			Address:   ads.GetAddress(),
			Email:     ads.GetEmail(),
			Phone:     ads.GetPhone(),
			BannerUrl: ads.GetBannerUrl(),
		},
	}, nil
}

func (*server) ReadAds(ctx context.Context, req *adspb.ReadAdsRequest) (*adspb.ReadAdsResponse, error) {
	fmt.Println("Read Ads Request")

	adsID := req.GetAdsId()
	oid, err := primitive.ObjectIDFromHex(adsID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot Parse ID"),
		)
	}

	//create an empty struct
	data := &adsItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find ads with specified ID"))
	}
	return &adspb.ReadAdsResponse{
		Ads: &adspb.Ads{
			Id:        data.ID.Hex(),
			UserId:    data.UserID,
			Title:     data.Title,
			Content:   data.Content,
			Address:   data.Address,
			Email:     data.Email,
			Phone:     data.Phone,
			BannerUrl: data.BannerURL,
		},
	}, nil

}

func main() {

	//we can get the file name and line number problem error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	//connect to mongodb
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Welcome to Ads Server \nNow Running Ads Server......")
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
