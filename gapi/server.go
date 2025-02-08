package gapi

import (
	"fmt"

	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/pb"
	"github.com/doctor12th/simple_bank_new/token"
	"github.com/doctor12th/simple_bank_new/util"
)

// Server serves gRPC requests for our banking service
type Server struct{
	pb.UnimplementedSimpleBankNewServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config ,store db.Store) (*Server,error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil,fmt.Errorf("cannot create token maker: %w",err)
	}
	server := &Server{
		config:config,
		store:store,
		tokenMaker: tokenMaker,
	}
	return server,nil
}