package api

import (
	"fmt"

	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/token"
	"github.com/doctor12th/simple_bank_new/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves http requests for our banking service
type Server struct{
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config ,store db.Store) (*Server,error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil,fmt.Errorf("cannot create token maker: %w",err)
	}
	server := &Server{
		config:config,
		store:store,
		tokenMaker: tokenMaker,
	}
	

	if v,ok:=binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
	}

	server.setupRouter()
	
	return server,nil
}
func (server *Server) setupRouter(){
	router :=gin.Default()
	router.POST("/users",server.createUser)
	router.POST("/users/login",server.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts",server.createAccount)
	authRoutes.GET("/accounts/:id",server.getAccount)
	authRoutes.GET("/accounts",server.listAccounts)
	authRoutes.POST("/transfers",server.createTransfer)

	
	server.router = router
}
//gin.H (key:value)
func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}
