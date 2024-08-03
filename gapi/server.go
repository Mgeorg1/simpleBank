package gapi

import (
	"fmt"
	db "github.com/Mgeorg1/simpleBank/db/sqlc"
	"github.com/Mgeorg1/simpleBank/pb"
	"github.com/Mgeorg1/simpleBank/token"
	"github.com/Mgeorg1/simpleBank/util"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %s", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.config = config
	return server, nil
}
