package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEventCreation ensures Event structs are created correctly
func TestEventCreation(t *testing.T) {
	// Create nodes
	sourceNode := NewNode("source", 0.0, 0.0, 0.0)
	targetNode := NewNode("target", 1.0, 1.0, 1.0)

	// Create an event
	eventTime := 10.0
	eventType := EventTypeTransactionCreated
	payload := "test payload"

	event := NewEvent(eventType, eventTime, sourceNode, targetNode, payload)

	// Assert that the event has the correct properties
	assert.Equal(t, eventType, event.Type)
	assert.Equal(t, eventTime, event.TimeScheduled)
	assert.Equal(t, sourceNode, event.SourceNode)
	assert.Equal(t, targetNode, event.TargetNode)
	assert.Equal(t, payload, event.Payload)
}

// TestEventQueueCreation ensures EventQueue structs are created correctly
func TestEventQueueCreation(t *testing.T) {
	// Create a new event queue
	queue := NewEventQueue()

	// Assert that the queue is initialized but empty
	assert.NotNil(t, queue.events)
	assert.Equal(t, 0, len(queue.events))
	assert.Equal(t, 0.0, queue.CurrentTime)
}

// TestEventQueueAddEvent tests adding events to the queue
func TestEventQueueAddEvent(t *testing.T) {
	// Create a new event queue
	queue := NewEventQueue()

	// Create nodes
	sourceNode := NewNode("source", 0.0, 0.0, 0.0)
	targetNode := NewNode("target", 1.0, 1.0, 1.0)

	// Create events
	event1 := NewEvent(EventTypeTransactionCreated, 10.0, sourceNode, targetNode, "payload1")
	event2 := NewEvent(EventTypeTransactionReceived, 5.0, targetNode, sourceNode, "payload2")

	// Add events to the queue
	queue.AddEvent(event1)
	queue.AddEvent(event2)

	// Assert that the queue has two events
	assert.Equal(t, 2, len(queue.events))
}

// TestEventQueueProcessNextEvent tests processing the next event in the queue
func TestEventQueueProcessNextEvent(t *testing.T) {
	// Create a new event queue
	queue := NewEventQueue()

	// Create nodes
	sourceNode := NewNode("source", 0.0, 0.0, 0.0)
	targetNode := NewNode("target", 1.0, 1.0, 1.0)

	// Create events with different times
	event1 := NewEvent(EventTypeTransactionCreated, 10.0, sourceNode, targetNode, "payload1")
	event2 := NewEvent(EventTypeTransactionReceived, 5.0, targetNode, sourceNode, "payload2")

	// Add events to the queue
	queue.AddEvent(event1)
	queue.AddEvent(event2)

	// Process the next event (should be event2 since it has earlier time)
	processedEvent := queue.ProcessNextEvent()

	// Assert that the processed event is event2
	assert.Equal(t, event2, processedEvent)

	// Assert that the queue now has only one event
	assert.Equal(t, 1, len(queue.events))

	// Assert that the current time has been updated
	assert.Equal(t, 5.0, queue.CurrentTime)

	// Process the next event (should be event1)
	processedEvent = queue.ProcessNextEvent()

	// Assert that the processed event is event1
	assert.Equal(t, event1, processedEvent)

	// Assert that the queue is now empty
	assert.Equal(t, 0, len(queue.events))

	// Assert that the current time has been updated
	assert.Equal(t, 10.0, queue.CurrentTime)
}

// TestEventDispatch ensures events get queued and processed in the correct time order
func TestEventDispatch(t *testing.T) {
	// Create a new event queue
	queue := NewEventQueue()

	// Create nodes
	sourceNode := NewNode("source", 0.0, 0.0, 0.0)
	targetNode := NewNode("target", 1.0, 1.0, 1.0)

	// Create events with different times, adding them out of order
	event1 := NewEvent(EventTypeTransactionCreated, 15.0, sourceNode, targetNode, "payload1")
	event2 := NewEvent(EventTypeTransactionReceived, 5.0, targetNode, sourceNode, "payload2")
	event3 := NewEvent(EventTypeBlockCreated, 10.0, sourceNode, nil, "payload3")

	// Add events to the queue
	queue.AddEvent(event1)
	queue.AddEvent(event2)
	queue.AddEvent(event3)

	// Process all events
	processedEvents := make([]*Event, 0)
	for len(queue.events) > 0 {
		processedEvents = append(processedEvents, queue.ProcessNextEvent())
	}

	// Assert that the events were processed in the correct order
	assert.Equal(t, 3, len(processedEvents))
	assert.Equal(t, event2, processedEvents[0]) // 5.0
	assert.Equal(t, event3, processedEvents[1]) // 10.0
	assert.Equal(t, event1, processedEvents[2]) // 15.0
}

// TestNodeEventHandling checks that each node updates state properly upon receiving events
func TestNodeEventHandling(t *testing.T) {
	// Create a network
	network := NewNetwork()

	// Create nodes
	node1 := NewNode("node1", 0.0, 0.0, 0.0)
	node2 := NewNode("node2", 1.0, 1.0, 1.0)

	// Add nodes to the network
	network.AddNode(node1)
	network.AddNode(node2)

	// Create an event queue
	queue := NewEventQueue()

	// Create a simulation
	sim := NewSimulation(network, queue)

	// Create a transaction event
	event := NewEvent(EventTypeTransactionCreated, 10.0, node1, node2, "test transaction")

	// Add the event to the queue
	queue.AddEvent(event)

	// Process the event
	sim.ProcessNextEvent()

	// Assert that node1's clock has been updated
	assert.Equal(t, 10.0, node1.Clock)

	// Create another event
	event2 := NewEvent(EventTypeTransactionReceived, 15.0, node1, node2, "test transaction")

	// Add the event to the queue
	queue.AddEvent(event2)

	// Process the event
	sim.ProcessNextEvent()

	// Assert that node2's clock has been updated
	assert.Equal(t, 15.0, node2.Clock)
}
