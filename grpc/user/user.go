package user

import (
	"context"
	"fmt"

	pb "github.com/aswcloud/idl/v1/servercomm"
	"github.com/aswcloud/server-comm/database"
)

type UserServer struct {
	pb.UnimplementedUserAccountServer
}

func (self *UserServer) CreateUser(ctx context.Context, data *pb.MakeUser) (*pb.Result, error) {
	db := database.New()
	db.Connect()
	defer db.Disconnect()

	// 회원가입 토큰이 정상적인지 검증함.
	result, err := db.RegisterTokenCollection().ExistsToken(data.Token)
	if !result {
		err_text := err.Error()
		return &pb.Result{
			Result: false,
			Error:  &err_text,
		}, nil
	}

	// email 데이터 분리
	email := ""
	if data.User.UserEmail != nil {
		email = *data.User.UserEmail
	}

	uuid, err := db.UserCollection().NewUser(
		data.User.UserId,
		data.User.UserPassword,
		data.User.UserNickname,
		email,
	)

	if err != nil {
		err_text := err.Error()
		return &pb.Result{
			Result: false,
			Error:  &err_text,
		}, nil
	}

	return &pb.Result{
		Result: true,
		Any:    []string{uuid},
	}, nil
}

func (self *UserServer) ReadUser(ctx context.Context, data *pb.Uuid) (*pb.UserDetail, error) {
	return &pb.UserDetail{}, nil
}

func (self *UserServer) UpdateUser(ctx context.Context, data *pb.User) (*pb.Result, error) {
	fmt.Println("HI!!!!")
	return &pb.Result{}, nil
}

func (self *UserServer) DeleteUser(ctx context.Context, data *pb.Uuid) (*pb.Result, error) {
	return &pb.Result{}, nil
}
