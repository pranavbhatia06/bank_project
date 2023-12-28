package sqlc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateEntry(t *testing.T) {
	CreateAccountData := CreateAccountParams{
		Owner:    "Pranav",
		Balance:  1950000,
		Currency: "INR",
	}
	account, err := queries.CreateAccount(context.Background(), CreateAccountData)
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.Owner, CreateAccountData.Owner)
	assert.Equal(t, account.Balance, CreateAccountData.Balance)
	assert.Equal(t, account.Currency, CreateAccountData.Currency)
	assert.NotNil(t, account.ID)
	assert.NotNil(t, account.CreatedAt)

	entryData := CreateEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}
	entry, err := queries.CreateEntry(context.Background(), entryData)
	assert.Nil(t, err)
	assert.Equal(t, entry.Amount, account.Balance)
	assert.Equal(t, entry.AccountID, account.ID)
	assert.NotNil(t, entry.ID)
	assert.NotNil(t, entry.CreatedAt)

}

func TestListEntry(t *testing.T) {
	CreateAccountData := CreateAccountParams{
		Owner:    "Pranav",
		Balance:  1950000,
		Currency: "INR",
	}
	account, err := queries.CreateAccount(context.Background(), CreateAccountData)
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.Owner, CreateAccountData.Owner)
	assert.Equal(t, account.Balance, CreateAccountData.Balance)
	assert.Equal(t, account.Currency, CreateAccountData.Currency)
	assert.NotNil(t, account.ID)
	assert.NotNil(t, account.CreatedAt)

	entryData := CreateEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}
	_, err = queries.CreateEntry(context.Background(), entryData)
	entryData.Amount = 15000
	_, err = queries.CreateEntry(context.Background(), entryData)
	entryData.Amount = -1500
	_, err = queries.CreateEntry(context.Background(), entryData)
	listEntry, err := queries.ListAccountStatement(context.Background(), account.ID)
	assert.NotNil(t, listEntry)
	assert.Equal(t, listEntry[2].Amount, int64(1950000))
	assert.Equal(t, listEntry[1].Amount, int64(15000))
	assert.Equal(t, listEntry[0].Amount, int64(-1500))

	assert.Nil(t, err)

}
