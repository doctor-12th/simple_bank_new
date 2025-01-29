package api

import (
	// "database/sql"
	// "log"
	"database/sql"
	"time"

	"github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	// "github.com/doctor12th/simple_bank/token"
	// "github.com/doctor12th/simple_bank/util"
	"net/http"
)

type createUserRequest struct{
	Username string `json:"username" binding:"required,alphanum" `
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
type UserResponse struct{
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt time.Time `json:"created_at"`
}
// type getAccountRequest struct{
// 	ID int64 `uri:"id" binding:"required,min=1"`
// }
// type listAccountsRequest struct{
// 	PageID int32 `form:"page_id" binding:"required,min=1"`
// 	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
// }

//插入数据库，并返回request
func (server *Server) createUser(ctx *gin.Context){
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashPassword,err:=util.HashedPassword(req.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashPassword,
		FullName: req.FullName,
		Email: req.Email,
	}
	user,err := server.store.CreateUser(ctx,arg)
	if err != nil {
		if pqErr,ok :=err.(*pq.Error);ok{
			// log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name(){
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

// func (server *Server) getUser(ctx *gin.Context){
// 	var req getAccountRequest
// 	if err := ctx.ShouldBindUri(&req); err!=nil{
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	account,err := server.store.GetAccount(ctx,req.ID)
// 	if err != nil {
// 		//找不到用户
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, account)
// }
type loginUserRequest struct{
	Username string `json:"username" binding:"required,alphanum" `
	Password string `json:"password" binding:"required,min=6"`
	// FullName string `json:"full_name" binding:"required"`
	// Email string `json:"email" binding:"required,email"`
}

type loginUserResponse struct{
	AccessToken string `json:"access_token"`
	User UserResponse `json:"user"`
}

func newUserResponse(user db.Users) UserResponse{
	return UserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user,err := server.store.GetUser(ctx,req.Username)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	err = util.CheckPassword(req.Password,user.HashedPassword)
	if err != nil{
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken,err := server.tokenMaker.CreateToken(user.Username,server.config.AccessTokenDuration)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}