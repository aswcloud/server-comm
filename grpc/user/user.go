package user

import (
	pb "github.com/aswcloud/idl/gen/go/v1"
)

type UserServer struct {
	pb.UnimplementedUserAccountServer
}
