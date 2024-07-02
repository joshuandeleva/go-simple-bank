package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/joshuandeleva/simplebank/db/mock"
	db "github.com/joshuandeleva/simplebank/db/sqlc"
	"github.com/joshuandeleva/simplebank/token"
	"github.com/joshuandeleva/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	// test cases

	testCases := []struct {
		name          string
		accountID     int64
		setupAuthFunc func(t *testing.T, request *http.Request , tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)                           // build stubs
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder) // check response
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code) // check if status code is 200
				requireBodyMatcher(t, recorder.Body, account)  // check if the body is the same as the account we sent
			},
		},
		{
			name: "Unauthorized",
			accountID: account.ID,
			setupAuthFunc: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			accountID: account.ID,
			setupAuthFunc: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Any()).
				Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code) // check if status code is 200
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // check if status code is 200
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
						GetAccount(gomock.Any() , gomock.Any()).
						Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code) // check if status code is 200
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// create a new mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl) //create a mock store
			tc.buildStubs(store)               // build stubs

			// start a test server and send request
			server := newTestServer(t,store)         // create a new server
			recorder := httptest.NewRecorder() // create a new recorder
			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuthFunc(t, request, server.tokenMaker) // add authorization header to request
			server.router.ServeHTTP(recorder, request) // send request to server
			tc.checkResponse(t, recorder)              // check response
		})
	}

}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatcher(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)

	require.NoError(t, err)
	require.Equal(t, account, getAccount) // check if the account we get is the same as the account we sent

}
