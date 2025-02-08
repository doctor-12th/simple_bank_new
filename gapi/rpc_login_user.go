package gapi

import (
	"context"
	"database/sql"

	"github.com/doctor12th/simple_bank_new/pb"
	"github.com/doctor12th/simple_bank_new/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user,err := server.store.GetUser(ctx,req.Username)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil,status.Errorf(codes.NotFound,"user not found:%s",err)
		}
		return nil,status.Errorf(codes.Internal,"failed to get user:%s",err)
	}
	
	err = util.CheckPassword(req.Password,user.HashedPassword)
	if err != nil{
		return nil,status.Errorf(codes.NotFound,"invalid password:%s",err)
	}
	accessToken,err := server.tokenMaker.CreateToken(user.Username,server.config.AccessTokenDuration)
	if err != nil{
		return nil,status.Errorf(codes.Internal,"failed to create access token:%s",err)
	}
	rsp := &pb.LoginUserResponse{
		User:convertUser(user),
		AccessToken:accessToken,
	}
	return rsp,nil

}