package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/joshuandeleva/simplebank/db/sqlc"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance : 0,
	}
	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqErr , ok:= err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
				case "unique_violation":
					ctx.JSON(http.StatusConflict, errorResponse(err))
					return
				case "check_violation":
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return
				case "foreign_key_violation":
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return
				default:
					log.Println(pqErr.Code.Name())
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context){
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		// if user with such id does not exist
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context){
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts , err := server.store.ListAccounts(ctx,arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)

}