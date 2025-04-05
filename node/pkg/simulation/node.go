package simulation

import (
	"math"
)

// Position represents a 3D position in space
type Position struct {
	X, Y, Z float64
}

// Node represents a node in the network
type Node struct {
	ID       string
	Position Position
	Clock    float64
	Latency  map[string]float64 // Map from node ID to base latency
	Velocity float64            // Optional: for relativistic effects
}

// NewNode creates a new node with the given ID and position
func NewNode(id string, x, y, z float64) *Node {
	return &Node{
		ID: id,
		Position: Position{
			X: x,
			Y: y,
			Z: z,
		},
		Clock:   0.0,
		Latency: make(map[string]float64),
	}
}

// Network represents a network of nodes
type Network struct {
	Nodes []*Node
}

// NewNetwork creates a new empty network
func NewNetwork() *Network {
	return &Network{
		Nodes: make([]*Node, 0),
	}
}

// AddNode adds a node to the network
func (n *Network) AddNode(node *Node) {
	n.Nodes = append(n.Nodes, node)
}

// Distance calculates the Euclidean distance between two nodes
func (n *Network) Distance(i, j int) float64 {
	if i < 0 || i >= len(n.Nodes) || j < 0 || j >= len(n.Nodes) {
		return 0.0
	}

	node1 := n.Nodes[i]
	node2 := n.Nodes[j]

	dx := node1.Position.X - node2.Position.X
	dy := node1.Position.Y - node2.Position.Y
	dz := node1.Position.Z - node2.Position.Z

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Latency calculates the latency between two nodes
// This is a simple model where latency = base latency + distance
func (n *Network) Latency(i, j int) float64 {
	if i < 0 || i >= len(n.Nodes) || j < 0 || j >= len(n.Nodes) {
		return 0.0
	}

	node1 := n.Nodes[i]
	node2 := n.Nodes[j]

	// Get base latency from the map
	baseLatency, exists := node1.Latency[node2.ID]
	if !exists {
		baseLatency = 0.0
	}

	// Calculate total latency as base latency + distance
	return baseLatency + n.Distance(i, j)
}
