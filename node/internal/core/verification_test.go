package core

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestVerification_AnchorTxCreation(t *testing.T) {
	sender := "verifierOrg"
	nonce := uint64(0)
	proofData := []byte("data to be anchored")
	proofHash := sha256.Sum256(proofData)

	t.Run("CreateValidAnchorTx", func(t *testing.T) {
		tx, err := CreateAnchorTransaction(nonce, sender, proofHash)

		if err != nil {
			t.Fatalf("CreateAnchorTransaction failed: %v", err)
		}
		if tx == nil {
			t.Fatalf("CreateAnchorTransaction returned nil transaction")
		}

		// Verify transaction fields
		if tx.Type != TxTypeAnchor {
			t.Errorf("Expected transaction type %v, got %v", TxTypeAnchor, tx.Type)
		}
		if tx.Nonce != nonce {
			t.Errorf("Expected nonce %d, got %d", nonce, tx.Nonce)
		}
		if tx.SenderID != sender {
			t.Errorf("Expected sender %s, got %s", sender, tx.SenderID)
		}
		if tx.RecipientID != "" { // Expect empty recipient for anchor
			t.Errorf("Expected empty recipient, got %s", tx.RecipientID)
		}
		if tx.Amount != 0 { // Expect 0 amount for anchor
			t.Errorf("Expected amount 0, got %d", tx.Amount)
		}
		if !bytes.Equal(tx.Payload, proofHash[:]) {
			t.Errorf("Payload mismatch: expected %x, got %x", proofHash[:], tx.Payload)
		}
		if tx.Signature != nil { // Should not be signed yet
			t.Errorf("Expected nil signature before signing, got %x", tx.Signature)
		}
	})

	t.Run("CreateAnchorTxEmptySender", func(t *testing.T) {
		_, err := CreateAnchorTransaction(nonce, "", proofHash)
		if err == nil {
			t.Errorf("Expected error for empty sender, but got nil")
		}
	})

	// Note: Test for nil proofHash might depend on Hash type definition,
	// assuming sha256.Sum256 always returns non-nil [32]byte.
}

// TODO: Add TestVerification_AnchorTxValidation (requires consensus/state integration)
// TODO: Add TestVerification_ProofRetrieval (requires state/storage implementation)
