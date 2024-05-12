package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Mgeorg1/simpleBank/db/mock"
	db "github.com/Mgeorg1/simpleBank/db/sqlc"
	"github.com/Mgeorg1/simpleBank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	accountFrom := db.Account{
		ID:       1,
		Owner:    util.RandomOwner(),
		Balance:  util.RandomInt(10, 1000),
		Currency: util.USD,
	}
	accountTo := db.Account{
		ID:       2,
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.USD,
	}

	amount := int64(10)
	testCases := []struct {
		name          string
		accountIDFrom int64
		accountIDTo   int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "OK",
			accountIDFrom: accountFrom.ID,
			accountIDTo:   accountTo.ID,
			body: gin.H{
				"from_account_id": accountFrom.ID,
				"to_account_id":   accountTo.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).Times(1).Return(accountFrom, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountTo.ID)).Times(1).Return(accountTo, nil)

				arg := db.CreateTransferParams{
					FromAccountID: accountFrom.ID,
					ToAccountID:   accountTo.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "Currency mismatch",
			accountIDFrom: accountFrom.ID,
			accountIDTo:   accountTo.ID,
			body: gin.H{
				"from_account_id": accountFrom.ID,
				"to_account_id":   accountTo.ID,
				"amount":          amount,
				"currency":        util.EUR,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountFrom.ID)).Times(1).Return(accountFrom, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//TODO: Make more test cases
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
