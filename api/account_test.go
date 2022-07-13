package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ruyoutor/simplebank/db/mock"
	db "github.com/ruyoutor/simplebank/db/sqlc"
	"github.com/ruyoutor/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	//build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), account.ID).
		Times(1).
		Return(account, nil)

	//start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var jsonAccount db.Account
	err = json.Unmarshal(data, &jsonAccount)
	require.NoError(t, err)

	require.Equal(t, account, jsonAccount)

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandonInit(1, 1000),
		Owner:    util.RandonOwn(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomMoney(),
	}
}
