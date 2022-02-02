package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/aswcloud/idl/gen/go/v1"
	"github.com/aswcloud/server-comm/middleware/auth"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServer
}

func send() {
	// 서버 연결 셋업
	conn, err := grpc.Dial("localhost:8088", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserClient(conn)
	// header := metadata.New(map[string]string{"authen": "barer aaaa"})
	// metadata.NewOutgoingContext(context.Background(), header)
	// grpc.SendHeader(ctx, header)
	// ctx = metadata.NewOutgoingContext(ctx, header)
	reply, err := c.GetUser(context.Background(), &pb.GetUserRequest{UserId: "testHello"})

	// GetHello 호출
	if err != nil {
		stat, ok := status.FromError(err)
		if ok {
			fmt.Println(stat.Code())
			fmt.Println(stat.Details())
			fmt.Println(stat.Err())
			fmt.Println(stat.Message())
		}
		log.Fatalf("GetHello error: %v", err)
	} else {
		log.Printf("Person: %v", reply)
	}

	return
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received profile: %v", in.GetUserId())
	ctx, err := auth.Authorization(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Length of `Name` cannot be more than 10 characters")
	}

	return &pb.GetUserResponse{UserMessage: &pb.UserMessage{UserId: in.GetUserId(), Name: "chacha", PhoneNumber: "1111-1111", Age: 30}}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8088")
	log.Print("TEST??")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	go func() {
		time.Sleep(time.Second * 5)
		send()
	}()

	pb.RegisterUserServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("ERROR!333")
}
