package main

import (
	"context"
	"fmt"
	"log"

	"../adspb"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Welcome ADS Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Error nih, tidak bisa connect : %v", err)
	}
	defer cc.Close()

	c := adspb.NewAdsServiceClient(cc)

	fmt.Println("Creating New ADS")
	ads := &adspb.Ads{
		UserId:    "Naufal Ziyad",
		Title:     "Beriklan di Adsmart",
		Content:   "Beriklan di Adsmart lebih mudah dan murah",
		Address:   "Jl. Mampang Prapatan kantor Transtv",
		Email:     "naufal.ziyad@detik.com",
		Phone:     "085954586600",
		BannerUrl: "adsmart.detik.com",
	}
	createAdsRes, err := c.CreateAds(context.Background(), &adspb.CreateAdsRequest{Ads: ads})
	if err != nil {
		log.Fatalf("unexepected error : %v", err)
	}
	fmt.Printf("Ads Successfull created ", createAdsRes)
}
