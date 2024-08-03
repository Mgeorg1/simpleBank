package gapi

import (
	"context"
	"errors"
	db "github.com/Mgeorg1/simpleBank/db/sqlc"
	"github.com/Mgeorg1/simpleBank/pb"
	"github.com/Mgeorg1/simpleBank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
			return nil, status.Errorf(codes.Internal, "failed creating user: %s", err)
		}
	}
	resp := &pb.CreateUserResponse{User: convertUser(user)}
	return resp, nil
}
