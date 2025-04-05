package core

import (
	"fmt"
	"sync"
)

// StateDB defines the interface for accessing and modifying the blockchain state.
// This could be backed by an in-memory map, a key-value store (like LevelDB), etc.
type StateDB interface {
	GetBalance(address string) (uint64, error)
	SetBalance(address string, balance uint64) error
	GetNonce(address string) (uint64, error)
	SetNonce(address string, nonce uint64) error
	// TODO: Add methods for contract storage, code, etc. later
}

// InMemoryStateDB provides a simple in-memory implementation of StateDB using maps.
// Note: This is not persistent and primarily for testing/early development.
type InMemoryStateDB struct {
	mu       sync.RWMutex // Mutex to protect concurrent access
	balances map[string]uint64
	nonces   map[string]uint64
	// TODO: Add maps for contract storage, code, etc.
}

// NewInMemoryStateDB creates a new in-memory state database.
func NewInMemoryStateDB() *InMemoryStateDB {
	return &InMemoryStateDB{
		balances: make(map[string]uint64),
		nonces:   make(map[string]uint64),
	}
}

// GetBalance retrieves the balance for a given address. Returns 0 if address not found.
func (db *InMemoryStateDB) GetBalance(address string) (uint64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	balance, _ := db.balances[address] // Returns 0 if not found, which is acceptable
	return balance, nil
}

// SetBalance sets the balance for a given address.
func (db *InMemoryStateDB) SetBalance(address string, balance uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	// TODO: Add validation for address format?
	db.balances[address] = balance
	return nil
}

// GetNonce retrieves the nonce for a given address. Returns 0 if address not found.
func (db *InMemoryStateDB) GetNonce(address string) (uint64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	nonce, _ := db.nonces[address] // Returns 0 if not found
	return nonce, nil
}

// SetNonce sets the nonce for a given address.
func (db *InMemoryStateDB) SetNonce(address string, nonce uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	// TODO: Add validation for address format?
	db.nonces[address] = nonce
	return nil
}

// StateManager orchestrates state changes by applying transactions.
type StateManager struct {
	db StateDB
}

// NewStateManager creates a new state manager.
func NewStateManager(db StateDB) *StateManager {
	if db == nil {
		// Or handle this more gracefully depending on requirements
		panic("StateDB cannot be nil for StateManager")
	}
	return &StateManager{db: db}
}

// ApplyTransaction validates a transaction against the current state and updates the state accordingly.
// Placeholder implementation - only handles basic transfer logic for now.
func (sm *StateManager) ApplyTransaction(tx *Transaction) error {
	if tx == nil {
		return fmt.Errorf("cannot apply nil transaction")
	}

	// Basic validation (stateless) - assumes signature is present for now
	if err := tx.ValidateBasic(); err != nil {
		return fmt.Errorf("basic transaction validation failed: %w", err)
	}

	// --- State Transition Logic (Minimal for Transfer) ---

	// TODO: Add signature verification check here once implemented properly
	// validSig, err := tx.VerifySignature()
	// if err != nil || !validSig { ... }

	// Get current state for sender and recipient
	senderNonce, err := sm.db.GetNonce(tx.SenderID)
	if err != nil {
		return fmt.Errorf("failed to get sender nonce for %s: %w", tx.SenderID, err)
	}
	senderBalance, err := sm.db.GetBalance(tx.SenderID)
	if err != nil {
		return fmt.Errorf("failed to get sender balance for %s: %w", tx.SenderID, err)
	}
	recipientBalance, err := sm.db.GetBalance(tx.RecipientID)
	if err != nil {
		// It's okay if recipient doesn't exist yet, balance is 0
		recipientBalance = 0
	}

	// TODO: Validate Nonce
	// if tx.Nonce != senderNonce {
	//  return fmt.Errorf("invalid nonce: expected %d, got %d", senderNonce, tx.Nonce)
	// }

	// TODO: Validate Balance (ensure senderBalance >= tx.Amount)
	// if senderBalance < tx.Amount {
	//  return fmt.Errorf("insufficient funds: sender %s has %d, needs %d", tx.SenderID, senderBalance, tx.Amount)
	// }

	// Perform state updates (minimal transfer logic)
	newSenderBalance := senderBalance - tx.Amount
	newRecipientBalance := recipientBalance + tx.Amount
	newSenderNonce := senderNonce + 1

	if err := sm.db.SetBalance(tx.SenderID, newSenderBalance); err != nil {
		return fmt.Errorf("failed to set sender balance: %w", err)
		// TODO: Consider state rollback mechanisms on partial failure
	}
	if err := sm.db.SetBalance(tx.RecipientID, newRecipientBalance); err != nil {
		return fmt.Errorf("failed to set recipient balance: %w", err)
		// TODO: Consider state rollback mechanisms on partial failure
	}
	if err := sm.db.SetNonce(tx.SenderID, newSenderNonce); err != nil {
		return fmt.Errorf("failed to set sender nonce: %w", err)
		// TODO: Consider state rollback mechanisms on partial failure
	}

	return nil // Success
}
