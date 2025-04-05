package core

import (
	"testing"
)

func TestStateDB_Initialization(t *testing.T) {
	t.Run("NewInMemoryStateDB", func(t *testing.T) {
		db := NewInMemoryStateDB()

		if db == nil {
			t.Fatalf("NewInMemoryStateDB returned nil")
		}
		if db.balances == nil {
			t.Errorf("InMemoryStateDB balances map is nil after initialization")
		}
		if db.nonces == nil {
			t.Errorf("InMemoryStateDB nonces map is nil after initialization")
		}
		// Check initial size (should be 0)
		if len(db.balances) != 0 {
			t.Errorf("Expected initial balances map size 0, got %d", len(db.balances))
		}
		if len(db.nonces) != 0 {
			t.Errorf("Expected initial nonces map size 0, got %d", len(db.nonces))
		}
	})
}

func TestStateDB_ReadWrite(t *testing.T) {
	db := NewInMemoryStateDB()
	addr1 := "address1"
	addr2 := "address2"

	// Test Balances
	t.Run("BalanceReadWrite", func(t *testing.T) {
		// 1. Check initial balance (should be 0)
		bal1Initial, err := db.GetBalance(addr1)
		if err != nil {
			t.Fatalf("GetBalance failed for initial state: %v", err)
		}
		if bal1Initial != 0 {
			t.Errorf("Expected initial balance 0 for %s, got %d", addr1, bal1Initial)
		}

		// 2. Set balance
		setBal := uint64(12345)
		err = db.SetBalance(addr1, setBal)
		if err != nil {
			t.Fatalf("SetBalance failed: %v", err)
		}

		// 3. Read balance back
		bal1Read, err := db.GetBalance(addr1)
		if err != nil {
			t.Fatalf("GetBalance failed after set: %v", err)
		}
		if bal1Read != setBal {
			t.Errorf("Expected balance %d for %s after set, got %d", setBal, addr1, bal1Read)
		}

		// 4. Check balance of another address (should still be 0)
		bal2, err := db.GetBalance(addr2)
		if err != nil {
			t.Fatalf("GetBalance failed for addr2: %v", err)
		}
		if bal2 != 0 {
			t.Errorf("Expected balance 0 for %s, got %d", addr2, bal2)
		}

		// 5. Overwrite balance
		overwriteBal := uint64(500)
		err = db.SetBalance(addr1, overwriteBal)
		if err != nil {
			t.Fatalf("SetBalance (overwrite) failed: %v", err)
		}
		bal1Overwrite, err := db.GetBalance(addr1)
		if err != nil {
			t.Fatalf("GetBalance failed after overwrite: %v", err)
		}
		if bal1Overwrite != overwriteBal {
			t.Errorf("Expected balance %d for %s after overwrite, got %d", overwriteBal, addr1, bal1Overwrite)
		}
	})

	// Test Nonces
	t.Run("NonceReadWrite", func(t *testing.T) {
		// 1. Check initial nonce (should be 0)
		nonce1Initial, err := db.GetNonce(addr1)
		if err != nil {
			t.Fatalf("GetNonce failed for initial state: %v", err)
		}
		if nonce1Initial != 0 {
			t.Errorf("Expected initial nonce 0 for %s, got %d", addr1, nonce1Initial)
		}

		// 2. Set nonce
		setNonce := uint64(5)
		err = db.SetNonce(addr1, setNonce)
		if err != nil {
			t.Fatalf("SetNonce failed: %v", err)
		}

		// 3. Read nonce back
		nonce1Read, err := db.GetNonce(addr1)
		if err != nil {
			t.Fatalf("GetNonce failed after set: %v", err)
		}
		if nonce1Read != setNonce {
			t.Errorf("Expected nonce %d for %s after set, got %d", setNonce, addr1, nonce1Read)
		}

		// 4. Check nonce of another address (should still be 0)
		nonce2, err := db.GetNonce(addr2)
		if err != nil {
			t.Fatalf("GetNonce failed for addr2: %v", err)
		}
		if nonce2 != 0 {
			t.Errorf("Expected nonce 0 for %s, got %d", addr2, nonce2)
		}

		// 5. Overwrite nonce
		overwriteNonce := uint64(10)
		err = db.SetNonce(addr1, overwriteNonce)
		if err != nil {
			t.Fatalf("SetNonce (overwrite) failed: %v", err)
		}
		nonce1Overwrite, err := db.GetNonce(addr1)
		if err != nil {
			t.Fatalf("GetNonce failed after overwrite: %v", err)
		}
		if nonce1Overwrite != overwriteNonce {
			t.Errorf("Expected nonce %d for %s after overwrite, got %d", overwriteNonce, addr1, nonce1Overwrite)
		}
	})
}

func TestStateTransition_Basic(t *testing.T) {
	db := NewInMemoryStateDB()
	sm := NewStateManager(db)

	addr1 := "senderA"
	addr2 := "recipientB"
	initialBal := uint64(1000)
	initialNonce := uint64(0)
	txAmount := uint64(100)

	// Setup initial state
	err := db.SetBalance(addr1, initialBal)
	if err != nil {
		t.Fatalf("Setup: SetBalance failed for sender: %v", err)
	}
	err = db.SetNonce(addr1, initialNonce)
	if err != nil {
		t.Fatalf("Setup: SetNonce failed for sender: %v", err)
	}

	t.Run("ApplyValidTransfer", func(t *testing.T) {
		// Create a basic transfer transaction
		tx := NewBaseTransaction(TxTypeTransfer, initialNonce, addr1, addr2, txAmount)
		_ = tx.Sign() // Use placeholder sign

		// Apply the transaction (this will fail until ApplyTransaction is implemented)
		err := sm.ApplyTransaction(tx)

		// --- Test Assertions (will fail initially) ---
		if err != nil {
			t.Errorf("ApplyTransaction failed unexpectedly for valid transfer: %v", err)
		}

		// Check sender's final state
		finalBalSender, _ := db.GetBalance(addr1)
		expectedBalSender := initialBal - txAmount
		if finalBalSender != expectedBalSender {
			t.Errorf("Sender balance incorrect: expected %d, got %d", expectedBalSender, finalBalSender)
		}
		finalNonceSender, _ := db.GetNonce(addr1)
		expectedNonceSender := initialNonce + 1
		if finalNonceSender != expectedNonceSender {
			t.Errorf("Sender nonce incorrect: expected %d, got %d", expectedNonceSender, finalNonceSender)
		}

		// Check recipient's final state
		finalBalRecipient, _ := db.GetBalance(addr2)
		expectedBalRecipient := txAmount // Assuming recipient started at 0
		if finalBalRecipient != expectedBalRecipient {
			t.Errorf("Recipient balance incorrect: expected %d, got %d", expectedBalRecipient, finalBalRecipient)
		}

		// TODO: Add tests for invalid nonce, insufficient balance etc. once ApplyTransaction implements checks

		// --- Reset state for potential future sub-tests ---
		_ = db.SetBalance(addr1, initialBal)
		_ = db.SetNonce(addr1, initialNonce)
		_ = db.SetBalance(addr2, 0)
	})

	// Test applying a nil transaction
	t.Run("ApplyNilTransaction", func(t *testing.T) {
		err := sm.ApplyTransaction(nil)
		if err == nil {
			t.Errorf("Expected error when applying nil transaction, but got nil")
		}
	})
}

// TODO: Add TestStateDB_ReadWrite
// TODO: Add TestStateTransition_Basic
