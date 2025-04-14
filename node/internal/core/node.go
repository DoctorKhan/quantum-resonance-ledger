package core

import (
	"errors"
	"fmt"
	"math"
)

// Placeholder for simulation constants - move to a config/constants file later
const pseudoSpeedOfLight = 10.0 // Example value for simple latency = distance / c

// Node represents a participant or point in the simulated spacetime network.
type Node struct {
	ID             string
	Position       Position
	Neighbors      map[string]*Node   // Connections to other nodes
	LatencyFactors map[string]float64 // Factors affecting communication latency (e.g., link quality)
	Balances       map[string]float64 // Token balances (e.g., QUSD, QRG, Gas)
	Inbox          []Event            // Received events, processed during node's turn
	// Add other fields as needed by the simulation design, e.g.:
	ProcessedEventCount int // Counter for testing event handling

	// State interface{} // Node-specific state (e.g., ledger, parameters) - Consider if Balances replaces part of this
}

// NewNode creates and initializes a new Node.
// For now, it handles basic initialization. Error handling for invalid inputs
// (like empty ID, if required) can be added later based on refined test cases.
func NewNode(id string, pos Position) *Node {
	// Handle empty ID case - decide on behavior (e.g., default ID, error)
	// For now, let's allow empty IDs as per the initial test log message,
	// but this might need stricter handling later.
	// if id == "" {
	//     // return nil, errors.New("node ID cannot be empty")
	//     // or generate a default ID
	// }

	return &Node{
		ID:             id,
		Position:       pos,
		Neighbors:      make(map[string]*Node),   // Initialize Neighbors map
		LatencyFactors: make(map[string]float64), // Initialize LatencyFactors map
		Balances:       make(map[string]float64), // Initialize Balances map
		Inbox:          make([]Event, 0),         // Initialize Inbox slice
		// Initialize other fields with sensible defaults, e.g.:
	}
}

// Distance calculates the Euclidean distance between this node and another node.
func (n *Node) Distance(other *Node) float64 {
	if other == nil {
		// Handle nil node case, perhaps return an error or MaxFloat64
		return math.MaxFloat64 // Or panic, or return an error
	}
	dx := n.Position.X - other.Position.X
	dy := n.Position.Y - other.Position.Y
	dz := n.Position.Z - other.Position.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Latency calculates the communication latency between this node and another node.
// Uses a simple model: Latency = Distance / SpeedOfLight.
// TODO: Implement relativistic latency model later.
func (n *Node) Latency(other *Node) float64 {
	if other == nil {
		return math.MaxFloat64 // Indicate impossibility or error
	}
	dist := n.Distance(other)
	if pseudoSpeedOfLight <= 0 {
		// Avoid division by zero or non-positive speed
		return math.MaxFloat64
	}
	if dist == 0 {
		return 0.0 // Latency is zero if distance is zero
	}
	return dist / pseudoSpeedOfLight
}

// Add other methods related to Node behavior here later, e.g.:
// func (n *Node) AddNeighbor(neighbor *Node) { ... }
// func (n *Node) CalculateDistance(other *Node) float64 { ... }

// UpdateLatencyFactor sets or updates the latency factor for a specific neighbor.
// Factors are assumed to be non-negative.
func (n *Node) UpdateLatencyFactor(neighborID string, factor float64) error {
	if factor < 0 {
		return errors.New("latency factor cannot be negative")
	}
	// Ensure the map is initialized (should be by NewNode, but double-check)
	if n.LatencyFactors == nil {
		n.LatencyFactors = make(map[string]float64)
	}
	n.LatencyFactors[neighborID] = factor
	return nil
}

// Deliver adds an event to the node's inbox.
// Placeholder: Actual event processing logic will be more complex.
func (n *Node) Deliver(event Event) {
	if n.Inbox == nil {
		// Defensive initialization (should be done by NewNode)
		n.Inbox = make([]Event, 0)
	}
	n.Inbox = append(n.Inbox, event)
}

// ProcessInbox iterates through the node's inbox and handles each event.
func (n *Node) ProcessInbox() {
	// Process events currently in the inbox
	eventsToProcess := n.Inbox
	n.Inbox = make([]Event, 0, len(eventsToProcess)) // Clear inbox efficiently

	for _, event := range eventsToProcess {
		switch e := event.(type) {
		case *TransactionCreatedEvent:
			n.handleTransactionCreated(e)
		// Add cases for other event types here
		default:
			// Handle unknown event types (e.g., log)
			fmt.Printf("Node %s received unknown event type: %T\n", n.ID, e)
		}
	}
}

// handleTransactionCreated processes a TransactionCreatedEvent.
// For now, just increments the counter for testing.
func (n *Node) handleTransactionCreated(event *TransactionCreatedEvent) {
	// Placeholder for actual transaction handling logic
	// fmt.Printf("Node %s handling TransactionCreated: %s\n", n.ID, event.TransactionID)
	n.ProcessedEventCount++
}
func (n *Node) CalculateWSI() float64 {
	// WSI calculation logic implementation
	// Calculate the WSI value based on the node's balances and latency factors
	wsiValue := 0.0
	for asset, balance := range n.Balances {
		// Calculate the asset's contribution to the WSI value
		wsiContribution := balance * n.LatencyFactors[asset]
		// Update the WSI value
		wsiValue += wsiContribution
	}
	return wsiValue
	// TO DO: implement WSI calculation logic
	return 0.0
}
