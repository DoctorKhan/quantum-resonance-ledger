package core

import (
	"fmt"
)

// VerificationAnchor represents the data being anchored.
// For now, just a simple hash, but could be more complex (e.g., Merkle root, commitment).
type VerificationAnchor struct {
	ProofHash Hash
	// TODO: Add metadata like timestamp, source identifier?
}

// CreateAnchorTransaction creates a transaction specifically for anchoring a proof hash.
func CreateAnchorTransaction(nonce uint64, sender string, proofHash Hash) (*Transaction, error) {
	if sender == "" {
		return nil, fmt.Errorf("anchor transaction requires a sender")
	}

	// Anchor transactions typically have 0 amount and no specific recipient (or self?)
	// The proof hash is stored in the payload.
	tx := NewBaseTransaction(TxTypeAnchor, nonce, sender, "", 0) // Empty recipient, 0 amount
	tx.Payload = proofHash[:]                                    // Store hash bytes in payload

	// The transaction still needs to be signed by the sender.
	// Signing happens separately via tx.Sign()

	return tx, nil
}

// TODO: Add VerificationManager if needed for tracking/querying anchors later.
// type VerificationManager struct { ... }
