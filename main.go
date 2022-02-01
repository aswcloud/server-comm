package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/aswcloud/idl/protos/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServer
}

func send() {
	// 서버 연결 셋업
	conn, err := grpc.Dial("localhost:8088", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := c.GetUser(ctx, &pb.GetUserRequest{UserId: "testHello"})

	// GetHello 호출
	if err != nil {
		log.Fatalf("GetHello error: %v", err)
	}
	log.Printf("Person: %v", reply)
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received profile: %v", in.GetUserId())
	return &pb.GetUserResponse{UserMessage: &pb.UserMessage{UserId: in.GetUserId(), Name: "chacha", PhoneNumber: "1111-1111", Age: 30}}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8088")
	log.Print("TEST??")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		),
	)

	go func() {
		time.Sleep(time.Second * 1)
		send()
	}()

	pb.RegisterUserServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
