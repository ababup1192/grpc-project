package main

import (
	"context"
	pb "grpc-project/infrastructure/service"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{Res: req.Req, Len: uint32(len(req.Req))}, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	listen, err := net.Listen("tcp", ":5050")

	if err != nil {
		log.Fatal(err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	log.Printf("server listening at %v", listen.Addr())

	go func() {
		<-ctx.Done()
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		s.GracefulStop()
	}()

	log.Fatal(s.Serve(listen))
}
