package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"../adspb"

	"google.golang.org/grpc"
)

type server struct {
}

func main() {

	//we can get the file name and line number problem error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Welcome to Server Ads")

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
	fmt.Println("End Program")
}
