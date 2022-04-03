// rpc CreateNamespace(namespace) returns (Result);
// // Needs, Authorization
// rpc ReadNamespace(Void) returns (list_namespace);
// // Needs, Authorization
// rpc UpdateNamespace(update_namespace) returns (Result);
// // Needs, Authorization
// rpc DeleteNamespace(namespace) returns (Result);
package kubernetes

import (
	"context"
	"fmt"
	"log"
	"os"

	pbk8s "github.com/aswcloud/idl/v1/kubernetes"
	pb "github.com/aswcloud/idl/v1/servercomm"
	"github.com/aswcloud/server-comm/database"
	"github.com/golang-jwt/jwt"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (self *KubernetesServer) CreateNamespace(ctx context.Context, data *pb.Namespace) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")[0]
	fmt.Println(p)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})
	if err != nil {
		e := err.Error()
		return &pb.Result{
			Result: false,
			Error:  &e,
		}, nil
	}

	log.Println(err)
	requestUuid := claims["user_id"].(string)

	db := database.New()
	db.Connect()
	defer db.Disconnect()
	_, exists := db.Namespace().CreateNamespace(requestUuid, data.Name)
	if !exists {
		k8s_server := os.Getenv("KUBERNETES_SERVER")
		conn, _ := grpc.Dial(k8s_server, grpc.WithInsecure(), grpc.WithBlock())
		log.Println(k8s_server)
		channel := pbk8s.NewKubernetesClient(conn)
		reply, err := channel.CreateNamespace(context.TODO(), &pbk8s.Namespace{
			Namespace: data.Name,
		})
		log.Println(reply, err)

		return &pb.Result{
			Result: reply.Result,
			Error:  reply.Error,
		}, nil

	} else {
		text := "allready exists"
		return &pb.Result{
			Result: false,
			Error:  &text,
		}, nil
	}
}

func (self *KubernetesServer) DeleteNamespace(ctx context.Context, data *pb.Namespace) (*pb.Result, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")[0]
	fmt.Println(p)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})
	if err != nil {
		e := err.Error()
		return &pb.Result{
			Result: false,
			Error:  &e,
		}, nil
	}

	log.Println(err)
	requestUuid := claims["user_id"].(string)

	db := database.New()
	db.Connect()
	defer db.Disconnect()
	_, exists := db.Namespace().DeleteNamespace(requestUuid, data.Name)
	if exists {
		k8s_server := os.Getenv("KUBERNETES_SERVER")
		conn, _ := grpc.Dial(k8s_server, grpc.WithInsecure(), grpc.WithBlock())
		log.Println(k8s_server)
		channel := pbk8s.NewKubernetesClient(conn)
		reply, err := channel.DeleteNamespace(context.TODO(), &pbk8s.Namespace{
			Namespace: data.Name,
		})
		log.Println(reply, err)

		return &pb.Result{
			Result: reply.Result,
			Error:  reply.Error,
		}, nil

	} else {
		text := "not found namespace"
		return &pb.Result{
			Result: false,
			Error:  &text,
		}, nil
	}
}
func (self *KubernetesServer) ReadNamespace(ctx context.Context, data *pb.Void) (*pb.ListNamespace, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")[0]
	fmt.Println(p)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})
	if err != nil {
		return &pb.ListNamespace{}, err
	}

	log.Println(err)
	requestUuid := claims["user_id"].(string)

	db := database.New()
	db.Connect()
	defer db.Disconnect()
	list, err := db.Namespace().ListNamespace(requestUuid)
	if err != nil {
		return &pb.ListNamespace{}, err
	}

	nameList := funk.Map(list, func(data string) *pb.Namespace {
		return &pb.Namespace{
			Name: data,
		}
	}).([]*pb.Namespace)

	return &pb.ListNamespace{
		List: nameList,
	}, nil
}
