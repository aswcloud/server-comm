package auth

import (
	"context"
	"os"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Authorization(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p := md.Get("authorization")
	if len(p) != 1 {
		return "", status.Errorf(codes.Unauthenticated, "not found auth token")
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(p[0], &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_TOKEN")), nil
	})
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	val, ok := claims["user_id"]

	if ok {
		return val.(string), nil
	} else {
		return "", status.Errorf(codes.Unauthenticated, "invalid auth jwt: %v", p[0])
	}
}
