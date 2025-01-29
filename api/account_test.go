package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/doctor12th/simple_bank_new/token"

	mockdb "github.com/doctor12th/simple_bank_new/db/mock"
	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)
func TestGetAccountAPI(t *testing.T){
	user,_ := randomUser(t)
	account := randomAccount(user.Username)
    testCases := []struct{
		name string
		accountID int64
		setupAuth func(t *testing.T, request *http.Request,tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request,tokenMaker token.Maker){
				addAuthorization(t,request,tokenMaker,authorizationTypeBearer,user.Username,time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					GetAccount(gomock.Any(),gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// {
		// 	name: "NotFound",
		// 	accountID: account.ID,
		// 	buildStubs: func(store *mockdb.MockStore){
		// 		store.EXPECT().
		// 			GetAccount(gomock.Any(),gomock.Eq(account.ID)).
		// 			Times(1).Return(db.Accounts{}, sql.ErrNoRows)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
		// 	},
		// },
	}
	for i := range testCases{
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).Return(account, nil)
			
			server := newTestServer(t,store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t,request,server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
			// require.Equal(t, http.StatusOK, recorder.Code)
			// requireBodyMatchAccount(t, recorder.Body, account)
		})
	}
	// require.Equal(t, account.ID, recorder.Body.String())

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Accounts){
	data,err:=io.ReadAll(body)
	require.NoError(t,err)
	var gotAccount db.Accounts
	err = json.Unmarshal(data,&gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)

}

func randomAccount(owner string) db.Accounts {
	return db.Accounts{
		ID: util.RandomInt(1,1000),
		Owner: owner,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func randomUser(t *testing.T) (user db.Users, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	user = db.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}