package core

import (
	"reflect"
	"testing"
)

func TestTransactionSerialization(t *testing.T) {
	t.Run("EncodeDecodeTransferTx", func(t *testing.T) {
		// 1. Create an original transaction
		originalTx := NewBaseTransaction(TxTypeTransfer, 1, "senderA", "recipientB", 100)
		// Add a mock signature for serialization test
		originalTx.Signature = Signature([]byte("test-sig"))

		// 2. Encode the transaction
		encodedData, err := originalTx.Encode()
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		if len(encodedData) == 0 {
			t.Fatalf("Encoded data is empty")
		}

		// 3. Decode the transaction
		decodedTx, err := DecodeTransaction(encodedData)
		if err != nil {
			t.Fatalf("DecodeTransaction failed: %v", err)
		}
		if decodedTx == nil {
			t.Fatalf("DecodeTransaction returned nil transaction")
		}

		// 4. Verify the decoded transaction matches the original
		// Using reflect.DeepEqual for struct comparison
		if !reflect.DeepEqual(originalTx, decodedTx) {
			t.Errorf("Decoded transaction does not match original.")
			t.Errorf("Original: %+v", originalTx)
			t.Errorf("Decoded:  %+v", decodedTx)
		}
	})

	// Test case: Decoding invalid data
	t.Run("DecodeInvalidData", func(t *testing.T) {
		invalidData := []byte("this is not valid gob data")
		_, err := DecodeTransaction(invalidData)
		if err == nil {
			t.Errorf("Expected error when decoding invalid data, but got nil")
		}
	})

	// Test case: Decoding empty data
	t.Run("DecodeEmptyData", func(t *testing.T) {
		emptyData := []byte{}
		_, err := DecodeTransaction(emptyData)
		if err == nil {
			t.Errorf("Expected error when decoding empty data, but got nil")
		}
		// Note: gob decoder might return EOF specifically here.
	})
}

func TestSignatureVerification(t *testing.T) {
	t.Run("SignAndVerify", func(t *testing.T) {
		// 1. Create a transaction
		tx := NewBaseTransaction(TxTypeTransfer, 1, "senderA", "recipientB", 100)

		// 2. Sign the transaction (using placeholder)
		err := tx.Sign()
		if err != nil {
			t.Fatalf("Sign failed unexpectedly: %v", err)
		}
		if tx.Signature == nil {
			t.Fatalf("Signature is nil after signing")
		}

		// 3. Verify the signature (using placeholder)
		valid, err := tx.VerifySignature()
		if err != nil {
			t.Fatalf("VerifySignature failed unexpectedly: %v", err)
		}

		// 4. Check validity (placeholder VerifySignature returns true)
		if !valid {
			t.Errorf("Expected signature to be valid, but VerifySignature returned false")
		}
		// TODO: Add test case for invalid signature once real crypto is implemented
	})

	// Test case: Verify signature on unsigned transaction
	t.Run("VerifyUnsignedTx", func(t *testing.T) {
		tx := NewBaseTransaction(TxTypeTransfer, 2, "senderC", "recipientD", 50)
		// tx.Signature is nil

		_, err := tx.VerifySignature()
		if err == nil {
			t.Errorf("Expected error when verifying signature on unsigned tx, but got nil")
		}
	})
}

// TODO: Add tests for other transaction types once defined

func TestTransactionValidation_Basic(t *testing.T) {
	t.Run("ValidSignedTx", func(t *testing.T) {
		tx := NewBaseTransaction(TxTypeTransfer, 1, "senderA", "recipientB", 100)
		_ = tx.Sign() // Use placeholder Sign

		err := tx.ValidateBasic()
		if err != nil {
			t.Errorf("ValidateBasic failed for a valid signed tx: %v", err)
		}
	})

	t.Run("InvalidUnsignedTx", func(t *testing.T) {
		tx := NewBaseTransaction(TxTypeTransfer, 1, "senderA", "recipientB", 100)
		// tx.Signature is nil

		err := tx.ValidateBasic()
		if err == nil {
			t.Errorf("Expected ValidateBasic to fail for an unsigned tx, but it passed")
		}
	})

	// TODO: Add more test cases for other basic validation rules later
}

// TODO: Add TestSignatureVerification
// TODO: Add TestTransactionValidation_Basic
