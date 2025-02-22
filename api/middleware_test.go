package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/doctor12th/simple_bank_new/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	// "github.com/stretchr/testify/require"
)
func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration){
		token ,err :=tokenMaker.CreateToken(username,duration)
		require.NoError(t,err)
		authorizationHeader := fmt.Sprintf("%s %s",authorizationType,token)
		request.Header.Set(authorizationHeaderKey,authorizationHeader)
}
func TestAuthMiddleware(t *testing.T){
	testCases := []struct{
		name string
		setupAuth func(t *testing.T, request *http.Request,tokenMaker token.Maker)
		checkResponse func(t *testing.T,recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth :func(t *testing.T, request *http.Request,tokenMaker token.Maker){
				addAuthorization(t,request,tokenMaker,authorizationTypeBearer,"user",time.Minute)
			},
			checkResponse:func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusOK,recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth :func(t *testing.T, request *http.Request,tokenMaker token.Maker){
				// addAuthorization(t,request,tokenMaker,authorizationTypeBearer,"user",time.Minute)
			},
			checkResponse:func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth :func(t *testing.T, request *http.Request,tokenMaker token.Maker){
				addAuthorization(t,request,tokenMaker,"unsupported","user",time.Minute)
			},
			checkResponse:func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth :func(t *testing.T, request *http.Request,tokenMaker token.Maker){
				addAuthorization(t,request,tokenMaker,"","user",time.Minute)
			},
			checkResponse:func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth :func(t *testing.T, request *http.Request,tokenMaker token.Maker){
				addAuthorization(t,request,tokenMaker,authorizationTypeBearer,"user", -time.Minute)
			},
			checkResponse:func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
	}

	for i := range testCases{
		tc :=testCases[i]
		t.Run(tc.name,func(t *testing.T){
			server :=newTestServer(t,nil)
			authPath := "/auth"
			//GET: 路径，中间件，JSON响应
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func (ctx *gin.Context){
					ctx.JSON(http.StatusOK,gin.H{})
				},
			)
			//模拟一个http响应记录器。
			recorder :=httptest.NewRecorder()
			//模拟发送一个http请求
			requset :=httptest.NewRequest(http.MethodGet,authPath,nil)
			// require.NoError(t,err)
			tc.setupAuth(t,requset,server.tokenMaker)
			server.router.ServeHTTP(recorder,requset)
			tc.checkResponse(t,recorder)

		})
	}
}