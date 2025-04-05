package simulation

import (
	"fmt"
	"math/rand" // For basic random position generation initially
	// Consider using a more robust random source if needed later
)

// Position represents a point in 3D space.
type Position struct {
	X, Y, Z float64
}

// Network represents the collection of nodes and their connections.
type Network struct {
	Nodes []*Node
	// Add other network-wide properties later, e.g.:
	// AdjacencyMatrix [][]float64 // Representing connections/distances
	// TopologyType string // e.g., "ring", "mesh", "random"
}

// NewNetwork creates a basic network with a specified number of nodes.
// Nodes are given sequential IDs and random initial positions.
// Topology information (connections) is not established here yet.
func NewNetwork(numNodes int) *Network {
	if numNodes < 0 {
		// Or return an error: return nil, errors.New("number of nodes cannot be negative")
		numNodes = 0 // Handle negative input gracefully for now
	}

	nodes := make([]*Node, numNodes)
	for i := 0; i < numNodes; i++ {
		id := fmt.Sprintf("node-%d", i)
		// Assign random positions for now. Topology-specific placement can be added later.
		// Ensure Position struct is accessible (defined in node.go currently)
		pos := Position{
			X: rand.Float64() * 100.0, // Example range
			Y: rand.Float64() * 100.0,
			Z: rand.Float64() * 100.0,
		}
		// Use the NewNode constructor from node.go
		nodes[i] = NewNode(id, pos)
		if nodes[i] == nil {
			// This shouldn't happen with the current NewNode, but good practice
			// Handle error, maybe return partially created network and an error
			fmt.Printf("Error: Failed to create node %s\n", id) // Basic logging
			// Consider returning an error instead of continuing
		}
	}

	return &Network{
		Nodes: nodes,
	}
}

// NewRingNetwork creates a network with the specified number of nodes
// connected in a ring topology.
func NewRingNetwork(numNodes int) *Network {
	if numNodes < 0 {
		numNodes = 0
	}

	// Create nodes using the basic constructor
	network := NewNetwork(numNodes)
	if network == nil || numNodes <= 1 {
		// Ring topology requires at least 2 nodes for meaningful connections
		// Return the network as is (0 or 1 node with no neighbors)
		return network
	}

	nodes := network.Nodes

	// Connect nodes in a ring
	for i := 0; i < numNodes; i++ {
		prevIndex := (i - 1 + numNodes) % numNodes
		nextIndex := (i + 1) % numNodes

		currentNode := nodes[i]
		prevNode := nodes[prevIndex]
		nextNode := nodes[nextIndex]

		// Neighbors map should be initialized by NewNode
		currentNode.Neighbors[prevNode.ID] = prevNode
		currentNode.Neighbors[nextNode.ID] = nextNode
	}

	return network
}

// Add methods for network operations later, e.g.:
// func (net *Network) GetNodeByID(id string) *Node { ... }
// func (net *Network) CalculateDistances() { ... }
// func (net *Network) EstablishConnections(topologyType string) { ... }
