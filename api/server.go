package api

import (
	"errors"
	"fmt"
	db "github.com/Mgeorg1/simpleBank/db/sqlc"
	"github.com/Mgeorg1/simpleBank/token"
	"github.com/Mgeorg1/simpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.CreateTransfer)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/tokens/renew_access", server.renewAccessToken)

	server.router = router
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

	val, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("error while creating custom validator")
	}
	err = val.RegisterValidation("currency", validCurrency)
	if err != nil {
		return nil, err
	}

	server.config = config
	server.setupRouter()
	return server, nil
}
