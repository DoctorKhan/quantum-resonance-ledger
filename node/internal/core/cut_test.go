package core

import (
	"bytes" // Needed for comparing slices
	"testing"
)

func TestCUT_Creation(t *testing.T) {
	t.Run("GenerateKeysBasic", func(t *testing.T) {
		sk, commitment, err := GenerateKeys()

		if err != nil {
			t.Fatalf("GenerateKeys failed unexpectedly: %v", err)
		}
		if sk == nil {
			t.Errorf("GenerateKeys returned nil secret key")
		}
		if commitment == nil {
			t.Errorf("GenerateKeys returned nil commitment")
		}
		// Basic length check for placeholder (SHA256 output)
		// TODO: Update this check when real crypto is implemented
		expectedCommitmentLen := 32
		if len(commitment) != expectedCommitmentLen {
			t.Errorf("Expected commitment length %d, got %d", expectedCommitmentLen, len(commitment))
		}
	})

	// TODO: Add tests for specific properties once real crypto is implemented
	// e.g., test commitment uniqueness, key properties, etc.
}

func TestCUT_CommitmentVerification(t *testing.T) {
	t.Run("VerifyCommitmentMatchesKey", func(t *testing.T) {
		// 1. Generate keys
		sk, generatedCommitment, err := GenerateKeys()
		if err != nil {
			t.Fatalf("GenerateKeys failed: %v", err)
		}
		if sk == nil || generatedCommitment == nil {
			t.Fatalf("GenerateKeys returned nil key or commitment")
		}

		// 2. Compute commitment separately
		computedCommitment, err := Commit(sk)
		if err != nil {
			t.Fatalf("Commit failed: %v", err)
		}
		if computedCommitment == nil {
			t.Fatalf("Commit returned nil commitment")
		}

		// 3. Verify they match (using bytes.Equal for slices)
		if !bytes.Equal(generatedCommitment, computedCommitment) {
			t.Errorf("Generated commitment (%x) does not match computed commitment (%x)",
				generatedCommitment, computedCommitment)
		}
	})

	// Test case: Commit with nil key
	t.Run("CommitNilKey", func(t *testing.T) {
		_, err := Commit(nil)
		if err == nil {
			t.Errorf("Expected error when committing nil key, but got nil")
		}
	})
}

func TestCUT_SpendProofGeneration(t *testing.T) {
	t.Run("GenerateProofBasic", func(t *testing.T) {
		// 1. Generate keys
		sk, _, err := GenerateKeys()
		if err != nil {
			t.Fatalf("GenerateKeys failed: %v", err)
		}

		// 2. Generate proof
		proof, err := GenerateSpendProof(sk)
		if err != nil {
			t.Fatalf("GenerateSpendProof failed: %v", err)
		}

		// 3. Basic check (proof is not nil)
		if proof == nil {
			t.Errorf("GenerateSpendProof returned nil proof")
		}
		// TODO: Add more specific checks when real ZKP is implemented
	})

	// Test case: Generate proof with nil key
	t.Run("GenerateProofNilKey", func(t *testing.T) {
		_, err := GenerateSpendProof(nil)
		if err == nil {
			t.Errorf("Expected error when generating proof with nil key, but got nil")
		}
	})
}

func TestCUT_SpendProofVerification(t *testing.T) {
	t.Run("VerifyValidProof", func(t *testing.T) {
		// 1. Generate keys and proof
		sk, commitment, err := GenerateKeys()
		if err != nil {
			t.Fatalf("GenerateKeys failed: %v", err)
		}
		proof, err := GenerateSpendProof(sk)
		if err != nil {
			t.Fatalf("GenerateSpendProof failed: %v", err)
		}

		// 2. Verify the proof
		valid, err := VerifySpendProof(commitment, proof)
		if err != nil {
			t.Fatalf("VerifySpendProof failed unexpectedly: %v", err)
		}

		// 3. Check validity (placeholder VerifySpendProof returns true)
		if !valid {
			t.Errorf("Expected proof to be valid, but VerifySpendProof returned false")
		}
		// TODO: Add test case for invalid proof once real ZKP is implemented
	})

	// Test case: Verify with nil commitment
	t.Run("VerifyNilCommitment", func(t *testing.T) {
		proof := SpendProof([]byte("mock-proof"))
		_, err := VerifySpendProof(nil, proof)
		if err == nil {
			t.Errorf("Expected error when verifying with nil commitment, but got nil")
		}
	})

	// Test case: Verify with nil proof
	t.Run("VerifyNilProof", func(t *testing.T) {
		commitment := Commitment([]byte("mock-commitment"))
		_, err := VerifySpendProof(commitment, nil)
		if err == nil {
			t.Errorf("Expected error when verifying with nil proof, but got nil")
		}
	})
}

// TODO: Add TestCUT_NoCloning
// TODO: Add TestCUT_Representation
