package simulation

import (
	"testing"
	"time"
	// "fmt" // Removed as it's unused now
)

func TestEventCreation(t *testing.T) {

	// Test case 1: Create TransactionCreatedEvent
	t.Run("CreateTransactionCreatedEvent", func(t *testing.T) {
		now := time.Now()
		targetID := "node-target-1"
		txID := "tx-123"
		event := NewTransactionCreatedEvent(now, targetID, txID)

		if event == nil {
			t.Fatalf("NewTransactionCreatedEvent returned nil")
		}

		// Check interface methods
		var baseEvent Event = event // Check if it implements the interface
		if baseEvent.Timestamp() != now {
			t.Errorf("Expected timestamp %v, got %v", now, baseEvent.Timestamp())
		}
		if baseEvent.Type() != "TransactionCreated" {
			t.Errorf("Expected type 'TransactionCreated', got '%s'", baseEvent.Type())
		}
		if baseEvent.GetTargetID() != targetID { // Check new interface method
			t.Errorf("Expected GetTargetID() '%s', got '%s'", targetID, baseEvent.GetTargetID())
		}

		// Check specific fields
		if event.TransactionID != txID {
			t.Errorf("Expected TransactionID '%s', got '%s'", txID, event.TransactionID)
		}
		if event.OccursAt != now { // Check BaseEvent field directly too
			t.Errorf("Expected BaseEvent.OccursAt %v, got %v", now, event.OccursAt)
		}
		if event.TargetID != targetID { // Check new TargetID field
			t.Errorf("Expected BaseEvent.TargetID '%s', got '%s'", targetID, event.TargetID)
		}
	})
}

// Add tests for other event types here later
// e.g., t.Run("CreateBlockCreatedEvent", func(t *testing.T) { ... })

func TestEventDispatching(t *testing.T) {
	scheduler := NewScheduler() // Uses placeholder scheduler
	targetNode := NewNode("node-1", Position{})
	scheduler.RegisterNode(targetNode)

	if targetNode == nil || scheduler == nil {
		t.Fatalf("Failed to create scheduler or node")
	}

	// Create an event targeted at the node
	eventTime := time.Now().Add(1 * time.Second) // Schedule for the future
	txID := "tx-dispatch-test"
	event := NewTransactionCreatedEvent(eventTime, targetNode.ID, txID)

	scheduler.Schedule(event) // Uses placeholder Schedule

	// Run simulation past the event time
	stopTime := eventTime.Add(1 * time.Millisecond)
	scheduler.RunUntil(stopTime) // Uses placeholder RunUntil

	// Check if the event arrived in the node's inbox
	if len(targetNode.Inbox) == 0 {
		t.Fatalf("Node inbox is empty, event was not dispatched")
	}
	if len(targetNode.Inbox) > 1 {
		t.Errorf("Expected 1 event in inbox, got %d", len(targetNode.Inbox))
	}

	receivedEvent := targetNode.Inbox[0]
	if receivedEvent == nil {
		t.Fatalf("Received event is nil")
	}

	// Verify the received event is the one we sent
	if receivedEvent.Timestamp() != eventTime {
		t.Errorf("Received event timestamp mismatch: expected %v, got %v", eventTime, receivedEvent.Timestamp())
	}
	if receivedEvent.Type() != event.Type() {
		t.Errorf("Received event type mismatch: expected %s, got %s", event.Type(), receivedEvent.Type())
	}
	// Try to type assert to check specific fields (optional but good)
	if txEvent, ok := receivedEvent.(*TransactionCreatedEvent); ok {
		if txEvent.TransactionID != txID {
			t.Errorf("Received event TransactionID mismatch: expected %s, got %s", txID, txEvent.TransactionID)
		}
		if txEvent.TargetID != targetNode.ID {
			t.Errorf("Received event TargetID mismatch: expected %s, got %s", targetNode.ID, txEvent.TargetID)
		}
	} else {
		t.Errorf("Received event is not of expected type *TransactionCreatedEvent")
	}

	// Test edge case: sending to non-existent node (Scheduler should handle this)
	// TODO: Add test case for dispatching to unknown targetID
}

// Add tests for Event Dispatching, Handling, Ordering, Timing later
