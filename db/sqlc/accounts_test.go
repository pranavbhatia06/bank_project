package sqlc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	CreateAccountData := CreateAccountParams{
		Owner:    "Lovish",
		Balance:  1500,
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
}
