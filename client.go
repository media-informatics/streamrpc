package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/media-informatics/streamrpc/service"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":56789", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	temperature := service.NewTemperatureServiceClient(conn)
	const repeat = 10
	count := &service.Request{Repeat: repeat}
	ctxt, cancel := context.WithTimeout(context.Background(), 2*repeat*time.Second)
	defer cancel()
	stream, err := temperature.Subscribe(ctxt, count)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
		tz, err := time.LoadLocation("Local")
		if err != nil {
			log.Printf("invalid time format %v", err)
		}
		time := resp.GetTime().AsTime().In(tz).Format("15:04:05")
		fmt.Printf("Response received: %v: %v\n", time, resp.GetTemperature())
	}
}
