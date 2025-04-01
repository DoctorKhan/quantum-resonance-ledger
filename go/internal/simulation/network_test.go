package simulation

import (
	"testing"
)

func TestNetworkCreation(t *testing.T) {
	// Test case 1: Create a small network (e.g., 3 nodes)
	t.Run("SmallNetwork", func(t *testing.T) {
		numNodes := 3
		// Assume a simple constructor for now, topology details later
		network := NewNetwork(numNodes)

		if network == nil {
			t.Fatalf("NewNetwork returned nil for %d nodes", numNodes)
		}
		if len(network.Nodes) != numNodes {
			t.Errorf("Expected %d nodes, got %d", numNodes, len(network.Nodes))
		}
		// Check if nodes have unique IDs (or other basic properties)
		ids := make(map[string]bool)
		for _, node := range network.Nodes {
			if node == nil {
				t.Errorf("Network contains a nil node")
				continue
			}
			if ids[node.ID] {
				t.Errorf("Duplicate node ID found: %s", node.ID)
			}
			ids[node.ID] = true
			// Basic check for position initialization (e.g., not zero value if random)
			// if node.Position == (Position{}) {
			// 	t.Errorf("Node %s position not initialized", node.ID)
			// }
		}
	})

	// Test case 2: Edge Case - Zero Nodes
	t.Run("ZeroNodes", func(t *testing.T) {
		numNodes := 0
		network := NewNetwork(numNodes)

		if network == nil {
			t.Fatalf("NewNetwork returned nil for %d nodes", numNodes)
		}
		if len(network.Nodes) != numNodes {
			t.Errorf("Expected %d nodes, got %d", numNodes, len(network.Nodes))
		}
	})

	// Test case 3: Edge Case - Single Node
	t.Run("SingleNode", func(t *testing.T) {
		numNodes := 1
		network := NewNetwork(numNodes)

		if network == nil {
			t.Fatalf("NewNetwork returned nil for %d nodes", numNodes)
		}
		if len(network.Nodes) != numNodes {
			t.Errorf("Expected %d nodes, got %d", numNodes, len(network.Nodes))
		}
		if network.Nodes[0] == nil {
			t.Errorf("Network contains a nil node for single node case")
		}
	})

	// Add tests for specific topologies (ring, mesh) later
	// when the constructor supports them.
}

func TestNetworkTopology(t *testing.T) {
	// Test case 1: Ring Topology
	t.Run("RingNetwork", func(t *testing.T) {
		numNodes := 5
		network := NewRingNetwork(numNodes) // Assumes NewRingNetwork exists

		if network == nil {
			t.Fatalf("NewRingNetwork returned nil for %d nodes", numNodes)
		}
		if len(network.Nodes) != numNodes {
			t.Fatalf("Expected %d nodes, got %d", numNodes, len(network.Nodes))
		}

		// Verify ring connections
		for i, node := range network.Nodes {
			if node == nil {
				t.Errorf("Node at index %d is nil", i)
				continue
			}
			if len(node.Neighbors) != 2 {
				t.Errorf("Node %s (index %d) expected 2 neighbors, got %d", node.ID, i, len(node.Neighbors))
				continue
			}

			// Calculate expected neighbor indices (wrapping around)
			prevIndex := (i - 1 + numNodes) % numNodes
			nextIndex := (i + 1) % numNodes
			expectedPrevID := network.Nodes[prevIndex].ID
			expectedNextID := network.Nodes[nextIndex].ID

			// Check if the expected neighbors are present
			if _, ok := node.Neighbors[expectedPrevID]; !ok {
				t.Errorf("Node %s missing neighbor %s (index %d)", node.ID, expectedPrevID, prevIndex)
			}
			if _, ok := node.Neighbors[expectedNextID]; !ok {
				t.Errorf("Node %s missing neighbor %s (index %d)", node.ID, expectedNextID, nextIndex)
			}
		}
	})

	// Add tests for other topologies (mesh, random) later
}
