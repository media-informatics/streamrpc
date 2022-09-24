package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/media-informatics/streamrpc/service"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", service.Port, "Service port")
	host = flag.String("host", service.Host, "Address of service machine")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect with server %v", err)
	}
	defer conn.Close()

	temperature := service.NewTemperatureServiceClient(conn)
	const repeat = 10
	count := &service.Request{Repeat: repeat}

	ctxt, cancel := context.WithTimeout(context.Background(), 2*repeat*time.Second)
	defer cancel()

	stream, err := temperature.Subscribe(ctxt, count)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	log.Printf("subscription start")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("subscription ended")
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
