package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/media-informatics/streamrpc/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	service.UnimplementedTemperatureServiceServer
}

func (s server) Subscribe(in *service.Request, srv service.TemperatureService_SubscribeServer) error {
	count := int(in.GetRepeat())
	log.Printf("subscribed - count = %d", count)
	for i := 0; i < count; i++ {
		select {
		case <-srv.Context().Done():
			log.Printf("client canceled!")
			count = i
			break

		default:
			tstamp := timestamppb.Now()
			temperature := float32(rand.NormFloat64()*2.0 + 28.0)
			resp := service.Response{
				Time:        tstamp,
				Temperature: temperature,
			}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
				count = i
				break
			}
			time.Sleep(time.Second)
		}
	}
	log.Printf("finished with %d Temperatures", count)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":56789")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterTemperatureServiceServer(s, server{})

	log.Println("start server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
