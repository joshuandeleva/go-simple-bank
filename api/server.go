package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/joshuandeleva/simplebank/db/sqlc"
	"github.com/joshuandeleva/simplebank/token"
	"github.com/joshuandeleva/simplebank/util"
)

// server struct to serve all http requests for a banking application

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new server instance and set up all api route

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setUpRouter()
	
	return server , nil
}

// set up routes

func (server *Server) setUpRouter() {
	router := gin.Default()

	// add routes
	router.POST("/accounts", server.createAccount)   // create account
	router.GET("/accounts/:id", server.getAccount)   // get account by id
	router.GET("/accounts", server.listAccounts)     // list all accounts
	router.POST("/transfers", server.createTransfer) // create transfer
	router.POST("/users", server.createUser)         // create user
	router.POST("/user/login", server.loginUser)      // login


	server.router = router

}

// Start runs the http server on a specific address

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// resuable error function

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
