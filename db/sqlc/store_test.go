package sqlc

import (
	"context"
	"fmt"

	"github.com/razorpay/bankProject/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := queries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	testStore := NewStore(testDb)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx operation %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), struct{}{}, txName)
			result, err := testStore.TransferTx(ctx, TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = queries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = queries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccountId
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccountId
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestDeadlockTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	testStore := NewStore(testDb)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx operation %d", i+1)
		var account3 int64
		var account4 int64
		if i%2 == 0 {
			account3 = account1.ID
			account4 = account2.ID
		} else {
			account4 = account1.ID
			account3 = account2.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), struct{}{}, txName)
			_, err := testStore.TransferTx(ctx, TransferTxParam{
				FromAccountID: account3,
				ToAccountID:   account4,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check results

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
