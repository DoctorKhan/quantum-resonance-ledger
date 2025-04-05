package simulation

// EventType represents the type of event
type EventType string

const (
	EventTypeTransactionCreated  EventType = "transaction_created"
	EventTypeTransactionReceived EventType = "transaction_received"
	EventTypeBlockCreated        EventType = "block_created"
	EventTypeBlockReceived       EventType = "block_received"
	EventTypeParameterUpdate     EventType = "parameter_update"
)

// Event represents an event in the simulation
type Event struct {
	Type          EventType
	TimeScheduled float64
	SourceNode    *Node
	TargetNode    *Node
	Payload       interface{}
}

// NewEvent creates a new event
func NewEvent(eventType EventType, timeScheduled float64, sourceNode, targetNode *Node, payload interface{}) *Event {
	return &Event{
		Type:          eventType,
		TimeScheduled: timeScheduled,
		SourceNode:    sourceNode,
		TargetNode:    targetNode,
		Payload:       payload,
	}
}

// EventQueue represents a priority queue of events ordered by time
type EventQueue struct {
	events      []*Event
	CurrentTime float64
}

// NewEventQueue creates a new empty event queue
func NewEventQueue() *EventQueue {
	return &EventQueue{
		events:      make([]*Event, 0),
		CurrentTime: 0.0,
	}
}

// AddEvent adds an event to the queue
func (q *EventQueue) AddEvent(event *Event) {
	q.events = append(q.events, event)
}

// ProcessNextEvent processes the next event in the queue (the one with the earliest time)
// and returns it
func (q *EventQueue) ProcessNextEvent() *Event {
	if len(q.events) == 0 {
		return nil
	}

	// Find the event with the earliest time
	earliestIndex := 0
	for i := 1; i < len(q.events); i++ {
		if q.events[i].TimeScheduled < q.events[earliestIndex].TimeScheduled {
			earliestIndex = i
		}
	}

	// Get the earliest event
	event := q.events[earliestIndex]

	// Remove the event from the queue
	q.events = append(q.events[:earliestIndex], q.events[earliestIndex+1:]...)

	// Update the current time
	q.CurrentTime = event.TimeScheduled

	return event
}

// EventCount returns the number of events in the queue
func (q *EventQueue) EventCount() int {
	return len(q.events)
}

// HasEvents returns true if the queue has events
func (q *EventQueue) HasEvents() bool {
	return len(q.events) > 0
}

// Simulation represents a simulation of the network
type Simulation struct {
	Network    *Network
	EventQueue *EventQueue
}

// NewSimulation creates a new simulation
func NewSimulation(network *Network, eventQueue *EventQueue) *Simulation {
	return &Simulation{
		Network:    network,
		EventQueue: eventQueue,
	}
}

// ProcessNextEvent processes the next event in the queue
func (s *Simulation) ProcessNextEvent() *Event {
	event := s.EventQueue.ProcessNextEvent()
	if event == nil {
		return nil
	}

	// Handle the event based on its type
	switch event.Type {
	case EventTypeTransactionCreated:
		// Update the source node's clock
		if event.SourceNode != nil {
			event.SourceNode.Clock = event.TimeScheduled
		}

	case EventTypeTransactionReceived:
		// Update the target node's clock
		if event.TargetNode != nil {
			event.TargetNode.Clock = event.TimeScheduled
		}

	case EventTypeBlockCreated:
		// Update the source node's clock
		if event.SourceNode != nil {
			event.SourceNode.Clock = event.TimeScheduled
		}

	case EventTypeBlockReceived:
		// Update the target node's clock
		if event.TargetNode != nil {
			event.TargetNode.Clock = event.TimeScheduled
		}

	case EventTypeParameterUpdate:
		// Handle parameter update
		// This will be implemented later
	}

	return event
}
