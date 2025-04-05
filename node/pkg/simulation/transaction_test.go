package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTransactionCreation tests the creation of a transaction
func TestTransactionCreation(t *testing.T) {
	// Create a transaction
	from := "sender"
	to := "receiver"
	amount := 10.0
	fee := 1.0
	receiverFee := 0.5

	tx := NewTransaction(from, to, amount, fee, receiverFee)

	// Assert that the transaction has the correct properties
	assert.Equal(t, from, tx.From)
	assert.Equal(t, to, tx.To)
	assert.Equal(t, amount, tx.Amount)
	assert.Equal(t, fee, tx.Fee)
	assert.Equal(t, receiverFee, tx.ReceiverFee)
}

// TestTransactionValidation tests the validation of transactions
func TestTransactionValidation(t *testing.T) {
	// Create a valid transaction
	tx := NewTransaction("sender", "receiver", 10.0, 1.0, 0.5)

	// Validate the transaction
	valid, err := tx.Validate()

	// Assert that the transaction is valid
	assert.True(t, valid)
	assert.NoError(t, err)

	// Create an invalid transaction with negative amount
	invalidTx1 := NewTransaction("sender", "receiver", -10.0, 1.0, 0.5)

	// Validate the transaction
	valid, err = invalidTx1.Validate()

	// Assert that the transaction is invalid
	assert.False(t, valid)
	assert.Error(t, err)

	// Create an invalid transaction with negative fee
	invalidTx2 := NewTransaction("sender", "receiver", 10.0, -1.0, 0.5)

	// Validate the transaction
	valid, err = invalidTx2.Validate()

	// Assert that the transaction is invalid
	assert.False(t, valid)
	assert.Error(t, err)

	// Create an invalid transaction with negative receiver fee
	invalidTx3 := NewTransaction("sender", "receiver", 10.0, 1.0, -0.5)

	// Validate the transaction
	valid, err = invalidTx3.Validate()

	// Assert that the transaction is invalid
	assert.False(t, valid)
	assert.Error(t, err)

	// Create an invalid transaction with empty sender
	invalidTx4 := NewTransaction("", "receiver", 10.0, 1.0, 0.5)

	// Validate the transaction
	valid, err = invalidTx4.Validate()

	// Assert that the transaction is invalid
	assert.False(t, valid)
	assert.Error(t, err)

	// Create an invalid transaction with empty receiver
	invalidTx5 := NewTransaction("sender", "", 10.0, 1.0, 0.5)

	// Validate the transaction
	valid, err = invalidTx5.Validate()

	// Assert that the transaction is invalid
	assert.False(t, valid)
	assert.Error(t, err)
}

// TestTransactionProcessing tests the processing of transactions
func TestTransactionProcessing(t *testing.T) {
	// Create a transaction manager
	manager := NewTransactionManager()

	// Create a transaction
	tx := NewTransaction("sender", "receiver", 10.0, 1.0, 0.5)

	// Add the transaction to the manager
	err := manager.AddTransaction(tx)
	assert.NoError(t, err)

	// Assert that the manager has one transaction
	assert.Equal(t, 1, len(manager.Transactions))

	// Create an invalid transaction
	invalidTx := NewTransaction("sender", "receiver", -10.0, 1.0, 0.5)

	// Try to add the invalid transaction to the manager
	err = manager.AddTransaction(invalidTx)
	assert.Error(t, err)

	// Assert that the manager still has only one transaction
	assert.Equal(t, 1, len(manager.Transactions))
}
