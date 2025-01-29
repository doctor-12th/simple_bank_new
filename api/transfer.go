package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/doctor12th/simple_bank_new/token"

	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)
type transferRequest struct{
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}


//插入数据库，并返回request
func (server *Server) createTransfer(ctx *gin.Context){
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount,valid :=server.validAccount(ctx,req.FromAccountID, req.Currency)
	
	if !valid{
		return
	}
	_,valid =server.validAccount(ctx,req.ToAccountID, req.Currency)
	if !valid{
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username{
		err :=errors.New("from account does not belong to this user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}
	result,err := server.store.TransferTx(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Accounts,bool) {
	account,err:=server.store.GetAccount(ctx,accountID)
	if err != nil {
		//找不到用户
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account,false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account,false
	}
	if account.Currency != currency{
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s",account.ID,account.Currency,currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account,false
	}
	return account,true
}