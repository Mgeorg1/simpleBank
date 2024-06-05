package api

import (
	"errors"
	db "github.com/Mgeorg1/simpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{store: store}
	router := gin.Default()

	val, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("error while creating custom validator")
	}
	err := val.RegisterValidation("currency", validCurrency)
	if err != nil {
		return nil, err
	}
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.CreateTransfer)
	router.POST("/users", server.createUser)
	server.router = router
	return server, nil
}
