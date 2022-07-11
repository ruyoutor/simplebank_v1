package db

import (
	"context"
	"testing"

	"github.com/ruyoutor/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccountForTransfer() (Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandonOwn(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	return testQueries.CreateAccount(context.Background(), arg)

}

func createRandomTransfer(t *testing.T) Transfer {

	accountFrom, err := createRandomAccountForTransfer()

	require.NoError(t, err)

	accountTo, err := createRandomAccountForTransfer()

	require.NoError(t, err)

	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        util.RandonInit(0, accountFrom.Balance),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotEmpty(t, transfer.ID)
	require.Equal(t, accountFrom.ID, transfer.FromAccountID)
	require.Equal(t, accountTo.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {

	createRandomTransfer(t)

}

func TestGetTransfer(t *testing.T) {

	transfer := createRandomTransfer(t)

	transferGet, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferGet)
	require.Equal(t, transfer.ID, transferGet.ID)
	require.Equal(t, transfer.Amount, transferGet.Amount)
	require.Equal(t, transfer.FromAccountID, transferGet.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferGet.ToAccountID)

}

func TestListTransfers(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, v := range transfers {
		require.NotEmpty(t, v)
	}

}
