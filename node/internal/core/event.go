package core

import (
	"fmt"  // For Printf in placeholder Execute
	"time" // For event timestamps
)

// Event defines the interface for simulation events.
type Event interface {
	Timestamp() time.Time            // When the event occurs
	Type() string                    // Type identifier string
	GetTargetID() string             // ID of the node the event is intended for
	Execute(world interface{}) error // Perform the event's action
	// Add other common methods if needed, e.g., SourceNodeID() string
}

// BaseEvent provides common fields for events.
type BaseEvent struct {
	OccursAt time.Time
	TargetID string // ID of the node the event is intended for
	// SourceID string // Optional: ID of the node originating the event
}

func (e BaseEvent) Timestamp() time.Time {
	return e.OccursAt
}

// GetTargetID returns the target node ID for the event.
func (e BaseEvent) GetTargetID() string {
	return e.TargetID
}

// --- Example Concrete Event Types ---

// TransactionCreatedEvent represents the creation of a new transaction.
// Placeholder - details to be filled in later (Phase 3.1)
type TransactionCreatedEvent struct {
	BaseEvent
	TransactionID string
	// Add transaction details: SenderID, RecipientID, Amount, etc.
}

// Type returns the event type identifier.
func (e TransactionCreatedEvent) Type() string {
	return "TransactionCreated"
}

// NewTransactionCreatedEvent creates a new TransactionCreatedEvent.
// Placeholder constructor - implementation might change.
func NewTransactionCreatedEvent(ts time.Time, targetID, txID string) *TransactionCreatedEvent {
	// Basic implementation for now
	return &TransactionCreatedEvent{
		BaseEvent:     BaseEvent{OccursAt: ts, TargetID: targetID},
		TransactionID: txID,
	}
}

// Execute performs the action associated with the TransactionCreatedEvent.
// Placeholder implementation.
func (e *TransactionCreatedEvent) Execute(world interface{}) error {
	// TODO: Implement actual event execution logic.
	// This might involve finding the target node in the 'world'
	// and calling a method on it, e.g., node.HandleTransactionCreated(e)
	fmt.Printf("Executing TransactionCreatedEvent for Tx %s (Target: %s) - Placeholder\n", e.TransactionID, e.TargetID)
	return nil // Return nil for now to satisfy the interface
}

// Add other event types later (BlockCreated, MessageSent, etc.)
