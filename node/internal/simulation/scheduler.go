package simulation

import (
	"container/heap"
	"fmt"
	"sort" // Required by heap interface implementation detail, even if not used directly
	"time"
)

// --- Priority Queue Implementation ---

// EventQueue implements heap.Interface and holds Events, ordered by timestamp.
type EventQueue []Event

func (eq EventQueue) Len() int { return len(eq) }

func (eq EventQueue) Less(i, j int) bool {
	// Min-heap based on Timestamp
	return eq[i].Timestamp().Before(eq[j].Timestamp())
}

func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

// Push and Pop use pointer receivers because they modify the slice's length.
func (eq *EventQueue) Push(x interface{}) {
	item := x.(Event) // Note: type assertion assumes x is always Event
	*eq = append(*eq, item)
}

func (eq *EventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	if n == 0 {
		return nil // Handle empty queue case
	}
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*eq = old[0 : n-1]
	return item
}

// Ensure EventQueue satisfies sort.Interface (needed by heap functions implicitly)
// This line isn't strictly necessary for functionality but serves as a compile-time check.
var _ heap.Interface = (*EventQueue)(nil)
var _ sort.Interface = (*EventQueue)(nil) // Also satisfies sort.Interface due to Len, Less, Swap

// --- Scheduler ---

// Scheduler manages the simulation clock and event queue.
type Scheduler struct {
	CurrentTime time.Time
	EventQueue  *EventQueue      // Use pointer to the heap implementation
	Nodes       map[string]*Node // Map of known nodes by ID
}

// NewScheduler creates a new simulation scheduler.
func NewScheduler() *Scheduler {
	pq := &EventQueue{}
	heap.Init(pq) // Initialize the heap
	return &Scheduler{
		CurrentTime: time.Time{}, // Zero time
		EventQueue:  pq,
		Nodes:       make(map[string]*Node),
	}
}

// RegisterNode adds a node to the scheduler's list of known nodes.
func (s *Scheduler) RegisterNode(node *Node) {
	if node != nil && s.Nodes != nil {
		s.Nodes[node.ID] = node
	}
}

// Schedule adds an event to the priority queue.
func (s *Scheduler) Schedule(event Event) {
	if event != nil && s.EventQueue != nil {
		heap.Push(s.EventQueue, event)
	}
}

// RunUntil executes events from the priority queue until the specified simulation time.
func (s *Scheduler) RunUntil(stopTime time.Time) {
	for s.EventQueue.Len() > 0 {
		// Peek at the next event's time without popping yet
		// Note: Accessing (*s.EventQueue)[0] assumes heap property is maintained
		// and peeks at the minimum element (earliest time).
		nextEventTimestamp := (*s.EventQueue)[0].Timestamp()

		// If the next event is after the stop time, we are done
		if nextEventTimestamp.After(stopTime) {
			break
		}

		// Pop the next event (minimum timestamp)
		event := heap.Pop(s.EventQueue).(Event)

		// Update simulation time to the event's time
		// Ensure time only moves forward
		if event.Timestamp().After(s.CurrentTime) {
			s.CurrentTime = event.Timestamp()
		} else if event.Timestamp().Before(s.CurrentTime) {
			// This shouldn't happen with a correct PQ and time handling
			fmt.Printf("Warning: Event timestamp %v is before current time %v\n", event.Timestamp(), s.CurrentTime)
			// Optionally skip or process anyway depending on desired strictness
		}
		// If timestamps are equal, process in pop order (heap stability not guaranteed)

		// Dispatch the event
		targetID := event.GetTargetID()
		if targetNode, ok := s.Nodes[targetID]; ok {
			targetNode.Deliver(event) // Call the Deliver method on the node
			// TODO: Consider if node processing (ProcessInbox) should happen here
			// or be scheduled as another event immediately following.
			// For now, Deliver just adds to inbox.
		} else {
			// Handle case where target node doesn't exist
			fmt.Printf("Warning: Target node %s not found for event type %s at time %v\n", targetID, event.Type(), s.CurrentTime)
		}
	}

	// After processing events up to their timestamps, advance time to stopTime
	// if it's later than the last event processed.
	if stopTime.After(s.CurrentTime) {
		s.CurrentTime = stopTime
	}
}
