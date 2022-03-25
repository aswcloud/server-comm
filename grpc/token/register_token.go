package token

import (
	"context"
	"encoding/hex"

	pb "github.com/aswcloud/idl/v1/servercomm"
	"github.com/aswcloud/server-comm/database"
)

func (self *TokenServer) CreateRegisterToken(ctx context.Context, data *pb.Void) (*pb.TokenMessage, error) {
	db := database.New()
	db.Connect()
	defer db.Disconnect()

	byte_token := db.RegisterTokenCollection().CreateToken(24 * 60 * 60)
	token := hex.EncodeToString(byte_token[:])

	return &pb.TokenMessage{
		Result: true,
		Token:  &token,
	}, nil

}
