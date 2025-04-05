package core

import (
	"fmt"
	"sync"
)

// TxPool manages pending transactions that haven't been included in a block yet.
type TxPool struct {
	mu sync.RWMutex
	// Simple storage: map[senderAddress]map[nonce]*Transaction
	pending map[string]map[uint64]*Transaction
	// TODO: Add more sophisticated data structures for prioritization (e.g., heap based on gas price)
	// TODO: Add limits (max transactions per account, max total transactions)
}

// NewTxPool creates a new transaction pool.
func NewTxPool() *TxPool {
	return &TxPool{
		pending: make(map[string]map[uint64]*Transaction),
	}
}

// AddTransaction attempts to add a transaction to the pool.
// Performs validation checks.
func (pool *TxPool) AddTransaction(tx *Transaction) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if tx == nil {
		return fmt.Errorf("cannot add nil transaction to pool")
	}

	// Basic validation (stateless)
	if err := tx.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid transaction (basic validation): %w", err)
	}

	// TODO: Add stateful validation using StateManager/StateDB:
	// - Check signature (properly, once implemented)
	// - Check nonce (must be current sender nonce from state)
	// - Check balance (sender must have sufficient funds for amount + gas)

	sender := tx.SenderID
	nonce := tx.Nonce

	// Initialize sender map if it doesn't exist
	if _, exists := pool.pending[sender]; !exists {
		pool.pending[sender] = make(map[uint64]*Transaction)
	}

	// Check if a transaction with the same sender and nonce already exists
	if existingTx, exists := pool.pending[sender][nonce]; exists {
		// TODO: Implement replacement logic (e.g., higher gas price)
		return fmt.Errorf("transaction with sender %s and nonce %d already exists in pool (tx hash: %v)", sender, nonce, existingTx) // Placeholder error
	}

	// Add the transaction
	pool.pending[sender][nonce] = tx
	fmt.Printf("TxPool: Added transaction from %s with nonce %d\n", sender, nonce) // Placeholder log

	return nil
}

// RemoveTransaction removes a transaction from the pool (e.g., after inclusion in a block).
func (pool *TxPool) RemoveTransaction(tx *Transaction) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if tx == nil {
		return // Nothing to remove
	}

	sender := tx.SenderID
	nonce := tx.Nonce

	if senderMap, senderExists := pool.pending[sender]; senderExists {
		if _, txExists := senderMap[nonce]; txExists {
			delete(senderMap, nonce)
			fmt.Printf("TxPool: Removed transaction from %s with nonce %d\n", sender, nonce) // Placeholder log
			// Clean up sender map if empty
			if len(senderMap) == 0 {
				delete(pool.pending, sender)
			}
		}
	}
}

// TODO: Add methods like:
// - GetPendingTransactions (for block proposal)
// - GetTransaction(hash)
// - UpdatePool (e.g., remove transactions invalidated by a new block)
// - PromoteExecutable (move transactions from future queue to pending when nonce matches)
