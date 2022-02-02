package auth

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func Authorization(ctx context.Context) (context.Context, error) {
	token, error := grpc_auth.AuthFromMD(ctx, "bearer")
	if error != nil {
		return nil, error
	}

	if token != "customToken" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", error)
	}
	newCtx := context.WithValue(ctx, "token", token)

	return newCtx, nil
}
