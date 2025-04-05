package core

import (
	"fmt"
	"time"
)

// BlockHeader represents the header of a block.
// Contains metadata about the block.
type BlockHeader struct {
	ParentHash Hash      // Hash of the parent block
	Number     uint64    // Block number
	Timestamp  time.Time // Timestamp of block creation
	StateRoot  Hash      // Root hash of the state trie after applying transactions
	TxRoot     Hash      // Root hash of the transaction trie
	// TODO: Add other fields like Difficulty, GasUsed, etc.
}

// Block represents a block in the blockchain.
type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction // List of transactions included in the block
	// TODO: Add Uncles/Ommer headers if applicable
}

// Hash calculates the block's hash (typically of the header).
// Placeholder implementation.
func (b *Block) Hash() (Hash, error) {
	// TODO: Implement proper hashing (e.g., RLP encode header, then Keccak256)
	// For now, return a placeholder based on number
	h := Hash{} // Placeholder
	copy(h[:], []byte(fmt.Sprintf("block-%d", b.Header.Number)))
	return h, nil
}

// NewBlock creates a new block.
// Placeholder implementation.
func NewBlock(header *BlockHeader, txs []*Transaction) *Block {
	return &Block{
		Header:       header,
		Transactions: txs,
	}
}

// --- Helper Type ---

// Hash represents a 32-byte hash.
// Placeholder - consider using a fixed-size array type from a crypto library.
type Hash [32]byte

// TODO: Add helper functions for Hash type (e.g., String(), Bytes(), SetBytes())
