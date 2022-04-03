package kubernetes

import (
	pb "github.com/aswcloud/idl/v1/servercomm"
)

type KubernetesServer struct {
	pb.UnimplementedKubernetesServer
}
