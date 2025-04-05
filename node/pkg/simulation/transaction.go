package simulation

import (
	"fmt"
)

// Transaction represents a transaction in the blockchain
type Transaction struct {
	From        string  // Sender's address
	To          string  // Receiver's address
	Amount      float64 // Amount to transfer
	Fee         float64 // Transaction fee paid by sender
	ReceiverFee float64 // Transaction fee paid by receiver (for receiver-pays model)
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount, fee, receiverFee float64) *Transaction {
	return &Transaction{
		From:        from,
		To:          to,
		Amount:      amount,
		Fee:         fee,
		ReceiverFee: receiverFee,
	}
}

// Validate checks if the transaction is valid
func (tx *Transaction) Validate() (bool, error) {
	// Check if sender and receiver are specified
	if tx.From == "" {
		return false, fmt.Errorf("sender address is empty")
	}

	if tx.To == "" {
		return false, fmt.Errorf("receiver address is empty")
	}

	// Check if amount is positive
	if tx.Amount <= 0 {
		return false, fmt.Errorf("amount must be positive")
	}

	// Check if fees are non-negative
	if tx.Fee < 0 {
		return false, fmt.Errorf("fee cannot be negative")
	}

	if tx.ReceiverFee < 0 {
		return false, fmt.Errorf("receiver fee cannot be negative")
	}

	return true, nil
}

// TransactionManager manages transactions in the blockchain
type TransactionManager struct {
	Transactions []*Transaction
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		Transactions: make([]*Transaction, 0),
	}
}

// AddTransaction adds a transaction to the manager after validating it
func (m *TransactionManager) AddTransaction(tx *Transaction) error {
	// Validate the transaction
	valid, err := tx.Validate()
	if !valid {
		return err
	}

	// Add the transaction to the list
	m.Transactions = append(m.Transactions, tx)

	return nil
}

// GetTransaction gets a transaction by index
func (m *TransactionManager) GetTransaction(index int) (*Transaction, error) {
	if index < 0 || index >= len(m.Transactions) {
		return nil, fmt.Errorf("transaction index out of range")
	}

	return m.Transactions[index], nil
}

// ProcessTransaction processes a transaction (in a real blockchain, this would update account balances)
// For now, this is just a placeholder
func (m *TransactionManager) ProcessTransaction(tx *Transaction) error {
	// Validate the transaction
	valid, err := tx.Validate()
	if !valid {
		return err
	}

	// In a real blockchain, we would update account balances here
	// For now, we'll just print a message
	fmt.Printf("Processing transaction: %s -> %s, Amount: %.2f, Fee: %.2f, ReceiverFee: %.2f\n",
		tx.From, tx.To, tx.Amount, tx.Fee, tx.ReceiverFee)

	return nil
}
