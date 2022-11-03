package main

import (
	"context"
	"fmt"
	pb "grpc-project/infrastructure/service"
	"log"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{Res: req.Req, Len: uint32(len(req.Req))}, nil
}

func main() {
	_ctx := context.Background()
	var eg *errgroup.Group

	listen, err := net.Listen("tcp", ":5050")
	errCh := make(chan error)

	eg, ctx := errgroup.WithContext(_ctx)
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	log.Printf("server listening at %v", listen.Addr())

	go func() {
		defer close(errCh)
		if err = s.Serve(listen); err != nil {
			errCh <- errors.Wrap(err, "failed to serve")
		}
	}()

	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		return
	}
}
