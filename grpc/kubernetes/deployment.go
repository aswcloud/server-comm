// // Needs, Authorization
// rpc CreateDeployment(deployment) returns (result);
// // Needs, Authorization
// rpc DeleteDeployment(delete_deployment) returns (result);
// // Needs, Authorization
// rpc ListDeployment(namespace) returns (list_deployment);

package kubernetes

import (
	"context"
	"encoding/json"
	"log"
	"os"

	pbk8s "github.com/aswcloud/idl/v1/kubernetes"
	pb "github.com/aswcloud/idl/v1/servercomm"
	jwtauth "github.com/aswcloud/server-comm/middleware/auth"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func (self *KubernetesServer) CreateDeployment(ctx context.Context, data *pb.Deployment) (*pb.Result, error) {
	_, err := jwtauth.Authorization(ctx)
	if err != nil {
		t := err.Error()
		return &pb.Result{
			Result: false,
			Error:  &t,
		}, nil
	}
	temp := pbk8s.Deployment{}
	bytee, _ := protojson.Marshal(data)
	json.Unmarshal(bytee, &temp)

	k8s_server := os.Getenv("KUBERNETES_SERVER")
	conn, _ := grpc.Dial(k8s_server, grpc.WithInsecure(), grpc.WithBlock())
	log.Println(k8s_server)
	channel := pbk8s.NewKubernetesClient(conn)

	reply, err := channel.CreateDeployment(context.TODO(), &temp)
	log.Println(reply, err)

	return &pb.Result{
		Result: reply.Result,
		Error:  reply.Error,
	}, nil
}
func (self *KubernetesServer) DeleteDeployment(ctx context.Context, data *pb.DeleteDeployment) (*pb.Result, error) {
	_, err := jwtauth.Authorization(ctx)
	if err != nil {
		t := err.Error()
		return &pb.Result{
			Result: false,
			Error:  &t,
		}, nil
	}

	temp := pbk8s.DeleteDeployment{}
	bytee, _ := protojson.Marshal(data)
	json.Unmarshal(bytee, &temp)

	k8s_server := os.Getenv("KUBERNETES_SERVER")
	conn, _ := grpc.Dial(k8s_server, grpc.WithInsecure(), grpc.WithBlock())
	log.Println(k8s_server)
	channel := pbk8s.NewKubernetesClient(conn)

	reply, err := channel.DeleteDeployment(context.TODO(), &temp)
	log.Println(reply, err)

	return &pb.Result{
		Result: reply.Result,
		Error:  reply.Error,
	}, nil
}
func (self *KubernetesServer) ReadDeployment(ctx context.Context, data *pb.Namespace) (*pb.ListDeployment, error) {
	_, err := jwtauth.Authorization(ctx)
	if err != nil {
		return &pb.ListDeployment{}, nil
	}

	k8s_server := os.Getenv("KUBERNETES_SERVER")
	conn, _ := grpc.Dial(k8s_server, grpc.WithInsecure(), grpc.WithBlock())
	log.Println(k8s_server)
	channel := pbk8s.NewKubernetesClient(conn)

	reply, err := channel.ListDeployment(context.TODO(), &pbk8s.Namespace{
		Namespace: data.Name,
	})
	log.Println(reply, err)
	sss := funk.Map(reply.Name, func(name string) *pb.Deployment {
		return &pb.Deployment{
			Name: name,
		}
	}).([]*pb.Deployment)

	return &pb.ListDeployment{
		List: sss,
	}, nil
}
