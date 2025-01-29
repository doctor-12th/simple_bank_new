package api

import (
	"database/sql"
	"errors"
	// "log"

	"github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	// "github.com/doctor12th/simple_bank/token"
	// "github.com/doctor12th/simple_bank/util"
	"net/http"
)

type createAccountRequest struct{
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}
type getAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}
type listAccountsRequest struct{
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

//插入数据库，并返回request
func (server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload :=ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner: authPayload.Username,
		Currency: req.Currency,
		Balance: 0,
	}
	account,err := server.store.CreateAccount(ctx,arg)
	if err != nil {
		if pqErr,ok :=err.(*pq.Error);ok{
			// log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name(){
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccount(ctx *gin.Context){
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account,err := server.store.GetAccount(ctx,req.ID)
	if err != nil {
		//找不到用户
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload :=ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username{
		err :=errors.New("this account is not belong to this user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) listAccounts(ctx *gin.Context){
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err!=nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload :=ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner: authPayload.Username,
		Offset: (req.PageID-1)*req.PageSize,
		Limit: req.PageSize,
	}
	account,err := server.store.ListAccounts(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}


