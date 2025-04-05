package core

import (
	"testing"
)

func TestTxPool_AddTransaction(t *testing.T) {
	pool := NewTxPool()
	sender := "senderA"
	recipient := "recipientB"

	// Test case 1: Add a valid transaction
	t.Run("AddValidTx", func(t *testing.T) {
		tx1 := NewBaseTransaction(TxTypeTransfer, 0, sender, recipient, 100)
		_ = tx1.Sign() // Placeholder sign

		err := pool.AddTransaction(tx1)
		if err != nil {
			t.Fatalf("AddTransaction failed for valid tx: %v", err)
		}

		// Verify transaction is in the pool (internal check)
		pool.mu.RLock()
		if _, senderExists := pool.pending[sender]; !senderExists {
			t.Errorf("Sender map not created for %s", sender)
		} else if _, txExists := pool.pending[sender][tx1.Nonce]; !txExists {
			t.Errorf("Transaction with nonce %d not found for sender %s", tx1.Nonce, sender)
		}
		pool.mu.RUnlock()

		// Clean up pool for next test (optional, but good practice)
		pool.RemoveTransaction(tx1)
	})

	// Test case 2: Add a transaction with the same nonce (should fail for now)
	t.Run("AddDuplicateNonceTx", func(t *testing.T) {
		tx1 := NewBaseTransaction(TxTypeTransfer, 1, sender, recipient, 100)
		_ = tx1.Sign()
		err1 := pool.AddTransaction(tx1) // Add first tx
		if err1 != nil {
			t.Fatalf("Failed to add initial tx: %v", err1)
		}

		tx2 := NewBaseTransaction(TxTypeTransfer, 1, sender, recipient, 200) // Same sender, same nonce
		_ = tx2.Sign()
		err2 := pool.AddTransaction(tx2) // Attempt to add second tx

		if err2 == nil {
			t.Errorf("Expected error when adding transaction with duplicate nonce, but got nil")
		}
		// TODO: Check specific error type/message once replacement logic is defined

		// Verify only the first transaction is present
		pool.mu.RLock()
		if senderMap, ok := pool.pending[sender]; ok {
			if len(senderMap) != 1 || senderMap[tx1.Nonce] != tx1 {
				t.Errorf("Pool state incorrect after adding duplicate nonce tx")
			}
		} else {
			t.Errorf("Sender map missing after adding duplicate nonce tx")
		}
		pool.mu.RUnlock()

		// Clean up
		pool.RemoveTransaction(tx1)
	})

	// Test case 3: Add a nil transaction (should fail)
	t.Run("AddNilTx", func(t *testing.T) {
		err := pool.AddTransaction(nil)
		if err == nil {
			t.Errorf("Expected error when adding nil transaction, but got nil")
		}
	})

	// Test case 4: Add an unsigned transaction (should fail basic validation)
	t.Run("AddUnsignedTx", func(t *testing.T) {
		tx := NewBaseTransaction(TxTypeTransfer, 2, sender, recipient, 50)
		// tx is not signed

		err := pool.AddTransaction(tx)
		if err == nil {
			t.Errorf("Expected error when adding unsigned transaction, but got nil")
		}
		// TODO: Check specific error message related to ValidateBasic failure
	})

	// TODO: Add tests for stateful validation (nonce, balance) once implemented
	// TODO: Add tests for pool limits (max txs, max per account) once implemented
	// TODO: Add tests for transaction replacement logic once implemented
}

func TestTxPool_RemoveTransaction(t *testing.T) {
	pool := NewTxPool()
	sender := "senderA"
	recipient := "recipientB"

	// Add some transactions first
	tx1 := NewBaseTransaction(TxTypeTransfer, 0, sender, recipient, 100)
	_ = tx1.Sign()
	_ = pool.AddTransaction(tx1)

	tx2 := NewBaseTransaction(TxTypeTransfer, 1, sender, recipient, 200)
	_ = tx2.Sign()
	_ = pool.AddTransaction(tx2)

	// Test case 1: Remove an existing transaction
	t.Run("RemoveExistingTx", func(t *testing.T) {
		pool.RemoveTransaction(tx1)

		// Verify tx1 is removed, tx2 remains
		pool.mu.RLock()
		if senderMap, ok := pool.pending[sender]; ok {
			if _, tx1Exists := senderMap[tx1.Nonce]; tx1Exists {
				t.Errorf("Transaction tx1 (nonce %d) was not removed", tx1.Nonce)
			}
			if _, tx2Exists := senderMap[tx2.Nonce]; !tx2Exists {
				t.Errorf("Transaction tx2 (nonce %d) was unexpectedly removed", tx2.Nonce)
			}
		} else {
			t.Errorf("Sender map for %s missing after removing tx1", sender)
		}
		pool.mu.RUnlock()
	})

	// Test case 2: Remove the same transaction again (should do nothing)
	t.Run("RemoveTxAgain", func(t *testing.T) {
		pool.RemoveTransaction(tx1) // tx1 is already removed

		// Verify tx2 still remains
		pool.mu.RLock()
		if senderMap, ok := pool.pending[sender]; ok {
			if _, tx2Exists := senderMap[tx2.Nonce]; !tx2Exists {
				t.Errorf("Transaction tx2 (nonce %d) was removed when removing tx1 again", tx2.Nonce)
			}
			if len(senderMap) != 1 {
				t.Errorf("Expected 1 transaction remaining for sender, got %d", len(senderMap))
			}
		} else {
			t.Errorf("Sender map for %s missing after removing tx1 again", sender)
		}
		pool.mu.RUnlock()
	})

	// Test case 3: Remove a transaction that was never added
	t.Run("RemoveNonExistentTx", func(t *testing.T) {
		tx3 := NewBaseTransaction(TxTypeTransfer, 0, "otherSender", recipient, 50)
		_ = tx3.Sign()
		pool.RemoveTransaction(tx3)

		// Verify tx2 still remains and nothing else changed
		pool.mu.RLock()
		if senderMap, ok := pool.pending[sender]; ok {
			if _, tx2Exists := senderMap[tx2.Nonce]; !tx2Exists {
				t.Errorf("Transaction tx2 (nonce %d) was removed when removing non-existent tx3", tx2.Nonce)
			}
			if len(senderMap) != 1 {
				t.Errorf("Expected 1 transaction remaining for sender, got %d", len(senderMap))
			}
		} else {
			t.Errorf("Sender map for %s missing after removing non-existent tx3", sender)
		}
		if _, otherSenderExists := pool.pending["otherSender"]; otherSenderExists {
			t.Errorf("Map for otherSender should not exist")
		}
		pool.mu.RUnlock()
	})

	// Test case 4: Remove nil transaction (should do nothing)
	t.Run("RemoveNilTx", func(t *testing.T) {
		pool.RemoveTransaction(nil)
		// Verify state is unchanged (tx2 still present)
		pool.mu.RLock()
		if senderMap, ok := pool.pending[sender]; ok {
			if _, tx2Exists := senderMap[tx2.Nonce]; !tx2Exists {
				t.Errorf("Transaction tx2 (nonce %d) was removed when removing nil tx", tx2.Nonce)
			}
		} else {
			t.Errorf("Sender map for %s missing after removing nil tx", sender)
		}
		pool.mu.RUnlock()
	})
}

// TODO: Add TestTxPool_RemoveTransaction
// TODO: Add TestTxPool_GetPendingTransactions
