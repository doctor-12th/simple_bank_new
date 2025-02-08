package gapi

import (
	"context"

	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/pb"
	"github.com/doctor12th/simple_bank_new/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashPassword,err:=util.HashedPassword(req.GetPassword())
	if err != nil{
		return nil,status.Errorf(codes.Internal,"failed to hash password:%s",err)
	}
	arg := db.CreateUserParams{
		Username: req.GetUsername(),
		HashedPassword: hashPassword,
		FullName: req.GetFullName(),
		Email: req.GetEmail(),
	}
	user,err := server.store.CreateUser(ctx,arg)
	if err != nil {
		if pqErr,ok :=err.(*pq.Error);ok{
			// log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name(){
			case "unique_violation":
				return nil,status.Errorf(codes.AlreadyExists,"username already exists:%s",err)
			}
		}
		return nil,status.Errorf(codes.Internal,"failed to create user:%s",err)
	}
	rsp := &pb.CreateUserResponse{
		User:convertUser(user),
	}

	return rsp,nil
}