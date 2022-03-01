package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/aswcloud/idl"
	"github.com/aswcloud/server-comm/database"
	"github.com/aswcloud/server-comm/grpc/organization"
	"github.com/aswcloud/server-comm/grpc/token"
	"github.com/aswcloud/server-comm/grpc/user"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"

	"google.golang.org/grpc"
)

func main() {
	gotenv.Load()

	lis, err := net.Listen("tcp", ":8088")
	log.Print("TEST??")
	if !database.New().Connect() {
		log.Fatal("DB ERROR!")
	}

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

	pb.RegisterOrganizationAccountServer(s, &organization.OrganizationServer{})
	pb.RegisterTokenServer(s, &token.TokenServer{})
	pb.RegisterUserAccountServer(s, &user.UserServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Println("ERROR!333")
}
