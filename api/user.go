package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/joshuandeleva/simplebank/db/sqlc"
	"github.com/joshuandeleva/simplebank/util"
	"github.com/lib/pq"
)


type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=10"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}


type createUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}



func (server *Server) createUser(ctx *gin.Context){
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
		
	}
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr , ok:= err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
				case "unique_violation":
					ctx.JSON(http.StatusConflict, errorResponse(err))
					return
				default:
					log.Println(pqErr.Code.Name())
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, resp)

}