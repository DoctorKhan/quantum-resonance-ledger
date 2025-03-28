package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNodeInitialization ensures Node structs are created correctly
func TestNodeInitialization(t *testing.T) {
	// Create a node with ID "node1" and position (1.0, 2.0, 3.0)
	node := NewNode("node1", 1.0, 2.0, 3.0)

	// Assert that the node has the correct ID
	assert.Equal(t, "node1", node.ID)

	// Assert that the node has the correct position
	assert.Equal(t, 1.0, node.Position.X)
	assert.Equal(t, 2.0, node.Position.Y)
	assert.Equal(t, 3.0, node.Position.Z)

	// Assert that the node's clock starts at 0
	assert.Equal(t, 0.0, node.Clock)

	// Assert that the node's latency map is initialized but empty
	assert.NotNil(t, node.Latency)
	assert.Equal(t, 0, len(node.Latency))
}

// TestNetworkInitialization ensures Network structs are created correctly
func TestNetworkInitialization(t *testing.T) {
	// Create a new network
	network := NewNetwork()

	// Assert that the network's nodes slice is initialized but empty
	assert.NotNil(t, network.Nodes)
	assert.Equal(t, 0, len(network.Nodes))
}

// TestAddNode tests adding a node to the network
func TestAddNode(t *testing.T) {
	// Create a new network
	network := NewNetwork()

	// Create a node
	node := NewNode("node1", 1.0, 2.0, 3.0)

	// Add the node to the network
	network.AddNode(node)

	// Assert that the network now has one node
	assert.Equal(t, 1, len(network.Nodes))

	// Assert that the node in the network is the one we added
	assert.Equal(t, node, network.Nodes[0])
}

// TestDistanceCalculations checks correctness of distance calculations
func TestDistanceCalculations(t *testing.T) {
	// Create a new network
	network := NewNetwork()

	// Create nodes at different positions
	node1 := NewNode("node1", 0.0, 0.0, 0.0)
	node2 := NewNode("node2", 3.0, 0.0, 0.0)
	node3 := NewNode("node3", 0.0, 4.0, 0.0)

	// Add nodes to the network
	network.AddNode(node1)
	network.AddNode(node2)
	network.AddNode(node3)

	// Calculate distances
	dist12 := network.Distance(0, 1) // Distance between node1 and node2
	dist13 := network.Distance(0, 2) // Distance between node1 and node3
	dist23 := network.Distance(1, 2) // Distance between node2 and node3

	// Assert that the distances are correct
	// node1 to node2: sqrt((3-0)^2 + (0-0)^2 + (0-0)^2) = 3
	assert.InDelta(t, 3.0, dist12, 0.0001)

	// node1 to node3: sqrt((0-0)^2 + (4-0)^2 + (0-0)^2) = 4
	assert.InDelta(t, 4.0, dist13, 0.0001)

	// node2 to node3: sqrt((3-0)^2 + (0-4)^2 + (0-0)^2) = 5
	assert.InDelta(t, 5.0, dist23, 0.0001)
}

// TestLatencyCalculations ensures the formula for latency is correct
func TestLatencyCalculations(t *testing.T) {
	// Create a new network
	network := NewNetwork()

	// Create nodes at different positions
	node1 := NewNode("node1", 0.0, 0.0, 0.0)
	node2 := NewNode("node2", 3.0, 0.0, 0.0)

	// Add nodes to the network
	network.AddNode(node1)
	network.AddNode(node2)

	// Set base latency between nodes
	node1.Latency[node2.ID] = 1.0
	node2.Latency[node1.ID] = 1.0

	// Calculate latency
	latency12 := network.Latency(0, 1) // Latency from node1 to node2

	// Assert that the latency is correct
	// Base latency is 1.0, and distance is 3.0, so total latency should be 1.0 + 3.0 = 4.0
	// (assuming latency is base + distance)
	assert.InDelta(t, 4.0, latency12, 0.0001)
}
