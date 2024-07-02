package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshuandeleva/simplebank/token"
	"github.com/stretchr/testify/require"
)

func addAuthorizationHeader(t *testing.T, request *http.Request, tokenMaker token.Maker , authorizationType string, username string , duration time.Duration) {
	token , err := tokenMaker.CreateToken(username , duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token) // Bearer <token>
	request.Header.Set(authorizationHeaderKey, authorizationHeader)

}

func TestAuthMiddleWare(t *testing.T){
	testCases := []struct{
		name string
		setupAuthFunc func(t *testing.T, request *http.Request , tokenMaker token.Maker)
		checkResponseFunc func(t *testing.T, recorder *httptest.ResponseRecorder) 
	}{
		{
			name: "OK",
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, "user1", time.Minute)
			},
			checkResponseFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {},
			checkResponseFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnSupportedAuthorization",
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, "unsupported", "user1", time.Minute)
			},
			checkResponseFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorization",
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, "", "user1", time.Minute)
			},
			checkResponseFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorization",
			setupAuthFunc: func(t *testing.T, request *http.Request , tokenMaker token.Maker) {
				addAuthorizationHeader(t, request, tokenMaker, authorizationTypeBearer, "user1", -time.Minute)
			},
			checkResponseFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath, 
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			tc.setupAuthFunc(t, request, server.tokenMaker) //setup the request to be authenticated
			server.router.ServeHTTP(recorder, request) // call the auth middleware
			tc.checkResponseFunc(t, recorder) // check the response
			
		})
	}

}