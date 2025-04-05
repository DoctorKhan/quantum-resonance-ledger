package core

import (
	"testing"
	"time"
)

// Basic test setup - might need refinement as scheduler evolves
func TestEventOrdering(t *testing.T) {
	scheduler := NewScheduler()
	node := NewNode("node-1", Position{})
	scheduler.RegisterNode(node)

	if node == nil || scheduler == nil {
		t.Fatalf("Failed to create scheduler or node")
	}

	// Schedule events out of order
	time1 := time.Unix(0, 100)                                   // Earlier time
	time2 := time.Unix(0, 200)                                   // Later time
	eventB := NewTransactionCreatedEvent(time2, node.ID, "tx-B") // Event B occurs later
	eventA := NewTransactionCreatedEvent(time1, node.ID, "tx-A") // Event A occurs earlier

	scheduler.Schedule(eventA) // Schedule earlier event second
	scheduler.Schedule(eventB) // Schedule later event first

	// Run simulation past both event times
	stopTime := time.Unix(0, 300)
	scheduler.RunUntil(stopTime) // Uses placeholder RunUntil

	// Verify events were processed and arrived in timestamp order in the inbox
	if len(node.Inbox) != 2 {
		t.Fatalf("Expected 2 events in inbox, got %d", len(node.Inbox))
	}

	// Check order (Event A should be first)
	if node.Inbox[0].Timestamp() != time1 {
		t.Errorf("Expected first event timestamp %v, got %v", time1, node.Inbox[0].Timestamp())
	}
	if txEventA, ok := node.Inbox[0].(*TransactionCreatedEvent); !ok || txEventA.TransactionID != "tx-A" {
		t.Errorf("Expected first event to be tx-A, got %v", node.Inbox[0])
	}

	if node.Inbox[1].Timestamp() != time2 {
		t.Errorf("Expected second event timestamp %v, got %v", time2, node.Inbox[1].Timestamp())
	}
	if txEventB, ok := node.Inbox[1].(*TransactionCreatedEvent); !ok || txEventB.TransactionID != "tx-B" {
		t.Errorf("Expected second event to be tx-B, got %v", node.Inbox[1])
	}
}

func TestEventTiming(t *testing.T) {
	scheduler := NewScheduler()
	node := NewNode("node-1", Position{})
	scheduler.RegisterNode(node)

	if node == nil || scheduler == nil {
		t.Fatalf("Failed to create scheduler or node")
	}

	eventTime := time.Unix(0, 500)
	event := NewTransactionCreatedEvent(eventTime, node.ID, "tx-timing")
	scheduler.Schedule(event)

	// Run just before the event time
	runTime1 := time.Unix(0, 499)
	scheduler.RunUntil(runTime1)
	if len(node.Inbox) != 0 {
		t.Errorf("Event processed prematurely at time %v, expected no events before %v", scheduler.CurrentTime, eventTime)
	}
	if scheduler.CurrentTime != runTime1 {
		t.Errorf("Expected scheduler time %v after RunUntil, got %v", runTime1, scheduler.CurrentTime)
	}

	// Run exactly to the event time
	runTime2 := eventTime
	scheduler.RunUntil(runTime2)
	if len(node.Inbox) != 1 {
		t.Errorf("Event not processed exactly at its timestamp. Inbox size: %d", len(node.Inbox))
	}
	if scheduler.CurrentTime != runTime2 {
		t.Errorf("Expected scheduler time %v after RunUntil, got %v", runTime2, scheduler.CurrentTime)
	}

	// Run past the event time (ensure no double processing)
	node.Inbox = make([]Event, 0) // Clear inbox manually for this check
	runTime3 := eventTime.Add(100 * time.Nanosecond)
	scheduler.RunUntil(runTime3)
	if len(node.Inbox) != 0 {
		t.Errorf("Event processed again after its timestamp. Inbox size: %d", len(node.Inbox))
	}
	if scheduler.CurrentTime != runTime3 {
		t.Errorf("Expected scheduler time %v after RunUntil, got %v", runTime3, scheduler.CurrentTime)
	}
}
