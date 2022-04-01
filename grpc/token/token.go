package token

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/aswcloud/idl/v1/servercomm"
	"github.com/aswcloud/server-comm/database"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

type TokenServer struct {
	pb.UnimplementedTokenServer
}

func createAccessJWT(uuid, nickname string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["type"] = "access"
	atClaims["user_id"] = uuid
	atClaims["name"] = nickname
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(time.Second * 30).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET_TOKEN")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func createRefreshJWT(uuid string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["type"] = "refresh"
	atClaims["user_id"] = uuid
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET_TOKEN")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (self *TokenServer) CreateRefreshToken(ctx context.Context, data *pb.CreateRefreshTokenMessage) (*pb.TokenMessage, error) {
	db := database.New()
	db.Connect()
	defer db.Disconnect()

	login := db.UserCollection().Login(data.UserId, data.UserPassword)
	if login == false {
		return &pb.TokenMessage{
			Result: false,
			Token:  nil,
		}, nil
	} else {
		uuid, err := db.UserCollection().GetUserUuid(data.UserId)
		if err != nil {
			return &pb.TokenMessage{
				Result: false,
				Token:  nil,
			}, nil
		}

		token, err := createRefreshJWT(uuid)
		if err != nil {
			return &pb.TokenMessage{
				Result: false,
				Token:  nil,
			}, nil
		}

		return &pb.TokenMessage{
			Result: true,
			Token:  &token,
		}, nil
	}
}

func (self *TokenServer) ReadRefreshToken(ctx context.Context, data *pb.Void) (*pb.RefreshTokenList, error) {
	return &pb.RefreshTokenList{}, nil
}

// dolor officia id exercitation

func (self *TokenServer) UpdatehRefreshToken(ctx context.Context, data *pb.Uuid) (*pb.TokenMessage, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")[0]
	fmt.Println(p)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})

	claimType := claims["type"]
	requestUuid := claims["user_id"]

	if err != nil || requestUuid != data.Uuid || claimType != "refresh" {
		return &pb.TokenMessage{
			Result: false,
			Token:  nil,
		}, nil
	}

	token, _ := createRefreshJWT(data.Uuid)
	return &pb.TokenMessage{
		Result: true,
		Token:  &token,
	}, nil
}

func (self *TokenServer) DeleteRefreshToken(ctx context.Context, data *pb.Uuid) (*pb.DeleteRefreshTokenMessage, error) {
	return &pb.DeleteRefreshTokenMessage{}, nil
}

func (self *TokenServer) CreateAccessToken(ctx context.Context, data *pb.Uuid) (*pb.TokenMessage, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")[0]
	fmt.Println(p)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})
	db := database.New()
	db.Connect()
	defer db.Disconnect()
	if err != nil {
		return &pb.TokenMessage{
			Result: false,
			Token:  nil,
		}, nil
	}
	requestUuid := claims["user_id"].(string)
	info, err := db.UserCollection().GetUserInfo(requestUuid)

	if err != nil || requestUuid != data.Uuid {
		return &pb.TokenMessage{
			Result: false,
			Token:  nil,
		}, nil
	}

	token, _ := createAccessJWT(info.Uuid, info.Nickname)
	return &pb.TokenMessage{
		Result: true,
		Token:  &token,
	}, nil
}
