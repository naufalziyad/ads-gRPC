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
	fmt.Printf("\nADS Successfull created !\n")

	//THIS SECTION FOR READ ADS FROM SERVER
	adsID := createAdsRes.GetAds().GetId()

	fmt.Println("Reading the Ads")
	fmt.Println("-------------------")

	/*
		for static ads id
		_, err2 := c.ReadAds(context.Background(), &adspb.ReadAdsRequest{AdsId: ""})
		if err2 != nil {
			fmt.Printf("Error while reading : %v \n", err2)
		}
	*/

	readAdsReq := &adspb.ReadAdsRequest{AdsId: adsID}
	readAdsRes, readAdsErr := c.ReadAds(context.Background(), readAdsReq)
	if readAdsErr != nil {
		fmt.Printf("Error while reading : %v \n", readAdsErr)
	}

	fmt.Printf("ADS was Read : %v \n", readAdsRes)

}
