package token

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/aswcloud/idl"
	"github.com/aswcloud/server-comm/database"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

type TokenServer struct {
	pb.UnimplementedTokenServer
}

func createToken(uuid string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = uuid
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET_TOKEN")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (self *TokenServer) CreateRefreshToken(ctx context.Context, data *pb.UserLoginMessage) (*pb.RefreshToken, error) {
	db := database.New()
	db.Connect()
	defer db.Disconnect()
	result := pb.RefreshToken{}

	login := db.GetUserCollection().Login(data.UserId, data.UserPassword)
	result.Result = login
	uuid, err := db.GetUserCollection().GetUserUuid(data.UserId)
	if err == nil {
		result.Uuid = &pb.Uuid{Uuid: uuid}
	} else {
		result.Uuid = &pb.Uuid{}
	}

	if login {
		token, err2 := createToken(uuid)

		if err2 != nil {
			return nil, err2
		}
		result.Uuid = &pb.Uuid{Uuid: uuid}
		result.Token = &token
		return &result, nil
	}
	return &result, nil
}

// dolor officia id exercitation

func (self *TokenServer) UpdatehRefreshToken(ctx context.Context, data *pb.Uuid) (*pb.RefreshToken, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")[0]
	fmt.Println(p)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})

	requestUuid := claims["user_id"]

	if err != nil || requestUuid != data.Uuid {
		return &pb.RefreshToken{
			Result: false,
			Uuid:   &pb.Uuid{Uuid: data.Uuid},
			Token:  nil,
		}, nil
	}

	token, _ := createToken(data.Uuid)
	return &pb.RefreshToken{
		Result: true,
		Uuid:   &pb.Uuid{Uuid: data.Uuid},
		Token:  &token,
	}, nil
}

func (self *TokenServer) DeleteRefreshToken(ctx context.Context, data *pb.Uuid) (*pb.LoginTokenMessage, error) {
	return &pb.LoginTokenMessage{}, nil
}

func (self *TokenServer) MakeAccessToken(ctx context.Context, data *pb.Uuid) (*pb.AccessToken, error) {
	return &pb.AccessToken{}, nil
}
