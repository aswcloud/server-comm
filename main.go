package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/aswcloud/idl/gen/go/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	header := metadata.New(map[string]string{"authen": "barer aaaa"})
	// metadata.NewOutgoingContext(context.Background(), header)
	grpc.SendHeader(ctx, header)
	ctx = metadata.NewOutgoingContext(ctx, header)
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

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		resp, err := handler(ctx, req)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
		}
		data := md["authen"]
		fmt.Println(data)

		return resp, err
	}
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
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New())),
			UnaryServerInterceptor(),
		),
	)

	go func() {
		time.Sleep(time.Second / 10)
		send()
	}()

	pb.RegisterUserServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
