package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Signature represents a cryptographic signature.
// Placeholder - replace with actual signature type (e.g., ECDSA signature).
type Signature []byte

// TransactionType defines the type of transaction.
type TransactionType uint8

const (
	TxTypeTransfer TransactionType = iota // Basic transfer
	TxTypeAnchor                          // Anchoring a proof/hash
	// Add other types later: Vote, BridgeIntent, Anchor, QSDMint, etc.
)

// Transaction represents a basic transaction structure.
// This will be expanded significantly for different transaction types.
type Transaction struct {
	Type        TransactionType
	Nonce       uint64
	SenderID    string // Placeholder for sender identifier (e.g., public key hash)
	RecipientID string // Placeholder for recipient identifier
	Amount      uint64 // Using uint64 for amount, assuming smallest unit (0 for anchor)
	Payload     []byte // Data payload (e.g., the hash/proof being anchored)
	Signature   Signature
	// TODO: Add GasPrice, GasLimit, Payload, etc. later
}

// NewTransaction creates a basic transfer transaction (unsigned).
// Placeholder - signing should happen separately.
func NewBaseTransaction(txType TransactionType, nonce uint64, sender, recipient string, amount uint64) *Transaction {
	return &Transaction{
		Type:        txType,
		Nonce:       nonce,
		SenderID:    sender,
		RecipientID: recipient,
		Amount:      amount, // Should be 0 for anchor
		Payload:     nil,    // Payload added separately for anchor
		Signature:   nil,    // Signature added via Sign method
	}
}

// Sign generates and attaches a signature to the transaction.
// Placeholder implementation.
func (tx *Transaction) Sign( /* privateKey ... */ ) error {
	// TODO: Implement actual signing using sender's private key
	// This should sign a hash of the transaction data (excluding the signature field itself)
	tx.Signature = []byte("mock-signature-for-" + tx.SenderID) // Insecure placeholder
	return nil
}

// VerifySignature checks if the transaction's signature is valid.
// Placeholder implementation.
func (tx *Transaction) VerifySignature() (bool, error) {
	// TODO: Implement actual signature verification using sender's public key
	// This involves hashing the transaction data (excluding signature) and verifying the signature against the hash and public key.
	if tx.Signature == nil {
		return false, fmt.Errorf("transaction has no signature")
	}
	fmt.Printf("Warning: Using placeholder signature verification for tx from %s\n", tx.SenderID)
	return true, nil // Assume valid for now
}

// Encode serializes the transaction into a byte slice using gob encoding.

// ValidateBasic performs stateless validation checks on the transaction.
// Checks format, presence of signature, etc. Does NOT check nonce or balance.
func (tx *Transaction) ValidateBasic() error {
	// TODO: Add more checks (e.g., non-zero amount for transfers? Sender/Recipient format?)
	if tx.Signature == nil {
		return fmt.Errorf("transaction is missing signature")
	}
	// Placeholder: Assume valid if signature exists for now
	return nil
}

func (tx *Transaction) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to gob encode transaction: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeTransaction deserializes a byte slice into a Transaction struct.
func DecodeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&tx)
	if err != nil {
		return nil, fmt.Errorf("failed to gob decode transaction: %w", err)
	}
	return &tx, nil
}
