package db

import (
	"context"
	"testing"
	"time"

	"github.com/ruyoutor/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccountForEntry() (Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandonOwn(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	return testQueries.CreateAccount(context.Background(), arg)

}

func createRandomEntry(t *testing.T) Entry {

	account, err := createRandomAccountForEntry()

	require.NoError(t, err)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandonInit(0, account.Balance),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.Equal(t, arg.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {

	createRandomEntry(t)

}

func TestGetEntry(t *testing.T) {

	entry := createRandomEntry(t)

	require.NotEmpty(t, entry)

	entryGet, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryGet)
	require.Equal(t, entry.ID, entryGet.ID)
	require.Equal(t, entry.Amount, entryGet.Amount)
	require.Equal(t, entry.AccountID, entryGet.AccountID)
	require.WithinDuration(t, entry.CreatedAt, entryGet.CreatedAt, time.Second)

}

func TestListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, v := range entries {
		require.NotEmpty(t, v)
	}

}
