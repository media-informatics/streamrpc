package main

import (
    "context"
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
	ctx := srv.Context()
	var err error = nil
	log.Printf("subscribed - count = %d", count)
	for i:= 0; i < count; i++ {
		err = sendTemperature(ctx, srv.Send, tick)
		if err != nil {
			log.Printf("server info: %v", err)
			count = i
			break
		}
	}
	log.Printf("finished with %d Temperatures", count)
	return err
}

type SendFunc func(r *service.Response) error

func sendTemperature(ctx context.Context, send SendFunc, tick <-chan time.Time) error {
    var err error = nil
    select {
	case <-ctx.Done():
		err = fmt.Errorf("temperature subscription aborted: %w", ctx.Err())
		
	case <-tick:
		resp := service.Response{
			Time:        timestamppb.Now(),
			Temperature: float32(rand.NormFloat64()*2.0 + 28.0),
		}
		if err = send(&resp); err != nil {
			err = fmt.Errorf("send error %w", err)
		}
	}
	return err
}

var (
	port = flag.Int("port", service.Port, "Service port")
	host = flag.String("host", service.Host, "Address of service machine")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := &TemperatureServiceServer{}
	service.RegisterTemperatureServiceServer(s, server)

	log.Printf("start server at %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
