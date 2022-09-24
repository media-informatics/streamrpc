package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/media-informatics/streamrpc/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TemperatureServiceServer struct {
	service.UnimplementedTemperatureServiceServer
}

func (s *TemperatureServiceServer) Subscribe(in *service.Request, srv service.TemperatureService_SubscribeServer) error {
	count := int(in.GetRepeat())
	tick := time.Tick(time.Second)
	ctxt := srv.Context()
	var err error = nil
	log.Printf("subscribed - count = %d", count)
	i := 0
	for i < count {
		select {
		case <-ctxt.Done():
			err = fmt.Errorf("temperature subscription aborted: %w", ctxt.Err())
		case <-tick:
			i++
			tstamp := timestamppb.Now()
			temperature := float32(rand.NormFloat64()*2.0 + 28.0)
			resp := service.Response{
				Time:        tstamp,
				Temperature: temperature,
			}
			if err = srv.Send(&resp); err != nil {
				err = fmt.Errorf("send error %w", err)
			}
		}
		if err != nil {
			log.Printf("server info: %v", err)
			break
		}
	}
	log.Printf("finished with %d Temperatures", i)
	return err
}

var (
	port = flag.Int("port", service.Port, "Service port")
	host = flag.String("host", service.Host, "Address of service machine")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := &TemperatureServiceServer{}
	service.RegisterTemperatureServiceServer(s, server)

	log.Println("start server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
