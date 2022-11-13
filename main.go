package main

import (
	"context"
	"fmt"
	pb "grpc-project/infrastructure/service"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type Server struct {
}

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{Res: req.Req, Len: uint32(len(req.Req))}, nil
}

// LoggedServerがInterfaceのメソッドを満たしているかの、アサーション
type ILogServer interface {
	Log(req interface{}, resp interface{}) error
}

var _ ILogServer = (*LoggedServer)(nil)
type LoggedServer struct {
	Server
}

func (s *LoggedServer) Log(req interface{}, resp interface{}) error {
	fmt.Printf("%+v\n", req)
	fmt.Printf("%+v\n", resp)

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	listen, err := net.Listen("tcp", ":5050")

	if err != nil {
		log.Fatal(err)
		return
	}

	chainUnaryServerInterceptor := grpc_middleware.ChainUnaryServer(UnaryLoggedServerInterceptor())
	s := grpc.NewServer(
		grpc.UnaryInterceptor(chainUnaryServerInterceptor),
	)

	pb.RegisterServiceServer(s, &Server{})

	pb.RegisterLoggedServiceServer(s, &LoggedServer{})

	log.Printf("server listening at %v", listen.Addr())

	go func() {
		<-ctx.Done()
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		s.GracefulStop()
	}()

	log.Fatal(s.Serve(listen))
}

func UnaryLoggedServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		srv, ok := info.Server.(ILogServer)

		if !ok {
			return handler(ctx, req)
		}

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		srv.Log(req, resp)

		return resp, nil
	}
}
