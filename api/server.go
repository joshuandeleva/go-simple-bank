package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/joshuandeleva/simplebank/db/sqlc"
)

// server struct to serve all http requests for a banking application

type Server struct {
	store db.Store
	router *gin.Engine
}

// NewServer creates a new server instance and set up all api route

func NewServer(store db.Store) *Server {
	server := &Server{store: store} // it creates  a new server
	router := gin.Default()

	// validators
	if v , ok := binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency", validCurrency)
	}
	// add routes
	router.POST("/accounts", server.createAccount) // create account
	router.GET("/accounts/:id", server.getAccount) // get account by id
	router.GET("/accounts", server.listAccounts) // list all accounts
	router.POST("/transfers", server.createTransfer) // create transfer
	
	server.router = router
	return server
}

// Start runs the http server on a specific address

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// resuable error function

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}