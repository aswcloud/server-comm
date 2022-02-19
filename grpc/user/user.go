package user

import (
	"context"

	pb "github.com/aswcloud/idl"
	"github.com/aswcloud/server-comm/database"
)

type UserServer struct {
	pb.UnimplementedUserAccountServer
}

func (self *UserServer) CreateUser(ctx context.Context, data *pb.MakeUser) (*pb.Result, error) {
	db := database.New()
	db.Connect()
	defer db.Disconnect()

	email := ""
	if data.User.UserEmail != nil {
		email = *data.User.UserEmail
	}

	uuid := db.GetUserCollection().NewUser(
		data.User.UserId,
		data.User.UserPassword,
		data.User.UserNickname,
		email,
	)
	if uuid == "" {
		errorText := "exists userid : " + data.User.UserId
		return &pb.Result{
			Result: false,
			Error:  &errorText,
		}, nil
	} else {
		result := pb.Result{}
		result.Result = true
		result.Any = append(result.Any, uuid)
		return &result, nil
	}
}

func (self *UserServer) ReadUser(ctx context.Context, data *pb.Uuid) (*pb.UserDetail, error) {
	return &pb.UserDetail{}, nil
}

func (self *UserServer) UpdateUser(ctx context.Context, data *pb.User) (*pb.Result, error) {
	return &pb.Result{}, nil
}

func (self *UserServer) DeleteUser(ctx context.Context, data *pb.Uuid) (*pb.Result, error) {
	return &pb.Result{}, nil
}
