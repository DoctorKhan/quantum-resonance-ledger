package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/khan/qrib/pkg/simulation"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Quantum-Inspired Blockchain Simulation")
	fmt.Println("======================================")

	// Create a network
	network := simulation.NewNetwork()
	fmt.Println("Created a new network")

	// Create nodes at different positions
	node1 := simulation.NewNode("node1", 0.0, 0.0, 0.0)
	node2 := simulation.NewNode("node2", 3.0, 0.0, 0.0)
	node3 := simulation.NewNode("node3", 0.0, 4.0, 0.0)
	node4 := simulation.NewNode("node4", 3.0, 4.0, 0.0)

	// Add nodes to the network
	network.AddNode(node1)
	network.AddNode(node2)
	network.AddNode(node3)
	network.AddNode(node4)
	fmt.Printf("Added %d nodes to the network\n", len(network.Nodes))

	// Set base latencies between nodes
	node1.Latency[node2.ID] = 1.0
	node1.Latency[node3.ID] = 1.5
	node1.Latency[node4.ID] = 2.0
	node2.Latency[node1.ID] = 1.0
	node2.Latency[node3.ID] = 2.0
	node2.Latency[node4.ID] = 1.5
	node3.Latency[node1.ID] = 1.5
	node3.Latency[node2.ID] = 2.0
	node3.Latency[node4.ID] = 1.0
	node4.Latency[node1.ID] = 2.0
	node4.Latency[node2.ID] = 1.5
	node4.Latency[node3.ID] = 1.0
	fmt.Println("Set base latencies between nodes")

	// Calculate distances and latencies
	fmt.Println("\nDistances between nodes:")
	for i := 0; i < len(network.Nodes); i++ {
		for j := i + 1; j < len(network.Nodes); j++ {
			distance := network.Distance(i, j)
			latency := network.Latency(i, j)
			fmt.Printf("  Distance from %s to %s: %.2f, Latency: %.2f\n",
				network.Nodes[i].ID, network.Nodes[j].ID, distance, latency)
		}
	}

	// Create a parameter manager
	paramManager := simulation.NewParameterManager()
	fmt.Println("\nCreated a parameter manager")

	// Create parameters
	blockSizeParam := simulation.NewParameter("block_size", 1.0, 10.0)
	txThroughputParam := simulation.NewParameter("tx_throughput", 100.0, 1000.0)
	latencyParam := simulation.NewParameter("latency", 0.1, 5.0)

	// Add parameters to the manager
	paramManager.AddParameter(blockSizeParam)
	paramManager.AddParameter(txThroughputParam)
	paramManager.AddParameter(latencyParam)
	fmt.Printf("Added %d parameters to the manager\n", len(paramManager.Parameters))

	// Create distributions
	blockSizeDist := simulation.NewUniformDistribution(blockSizeParam)
	txThroughputDist := simulation.NewNormalDistribution(txThroughputParam, 500.0, 100.0)
	latencyDist := simulation.NewNormalDistribution(latencyParam, 2.0, 0.6)

	// Set distributions
	paramManager.SetDistribution(blockSizeParam, blockSizeDist)
	paramManager.SetDistribution(txThroughputParam, txThroughputDist)
	paramManager.SetDistribution(latencyParam, latencyDist)
	fmt.Println("Set distributions for all parameters")

	// Create uncertainty relations
	relation1 := simulation.NewUncertaintyRelation(blockSizeParam, latencyParam, 1.0)
	relation2 := simulation.NewUncertaintyRelation(txThroughputParam, latencyParam, 50.0)

	// Add relations to the manager
	paramManager.AddUncertaintyRelation(relation1)
	paramManager.AddUncertaintyRelation(relation2)
	fmt.Printf("Added %d uncertainty relations to the manager\n", len(paramManager.UncertaintyRelations))

	// Validate uncertainty relations
	valid, violations := paramManager.ValidateUncertaintyRelations()
	if valid {
		fmt.Println("All uncertainty relations are satisfied")
	} else {
		fmt.Printf("Some uncertainty relations are violated: %d violations\n", len(violations))
		for i, relation := range violations {
			fmt.Printf("  Violation %d: Relation between %s and %s (constant: %.2f)\n",
				i+1, relation.Parameter1.Name, relation.Parameter2.Name, relation.Constant)
		}
	}

	// Sample parameters
	paramManager.SampleParameters()
	fmt.Println("\nSampled parameter values:")
	for _, param := range paramManager.Parameters {
		fmt.Printf("  %s: %.2f\n", param.Name, param.CurrentValue)
	}

	// Create a transaction manager
	txManager := simulation.NewTransactionManager()
	fmt.Println("\nCreated a transaction manager")

	// Create transactions
	tx1 := simulation.NewTransaction("node1", "node2", 10.0, 1.0, 0.5)
	tx2 := simulation.NewTransaction("node2", "node3", 5.0, 0.5, 0.25)
	tx3 := simulation.NewTransaction("node3", "node4", 7.5, 0.75, 0.375)

	// Add transactions to the manager
	txManager.AddTransaction(tx1)
	txManager.AddTransaction(tx2)
	txManager.AddTransaction(tx3)
	fmt.Printf("Added %d transactions to the manager\n", len(txManager.Transactions))

	// Create an event queue
	eventQueue := simulation.NewEventQueue()
	fmt.Println("\nCreated an event queue")

	// Create a simulation
	sim := simulation.NewSimulation(network, eventQueue)
	fmt.Println("Created a simulation")

	// Create events for transactions
	event1 := simulation.NewEvent(simulation.EventTypeTransactionCreated, 10.0, node1, node2, tx1)
	event2 := simulation.NewEvent(simulation.EventTypeTransactionCreated, 15.0, node2, node3, tx2)
	event3 := simulation.NewEvent(simulation.EventTypeTransactionCreated, 20.0, node3, node4, tx3)

	// Add events to the queue
	eventQueue.AddEvent(event1)
	eventQueue.AddEvent(event2)
	eventQueue.AddEvent(event3)
	fmt.Printf("Added %d events to the queue\n", eventQueue.EventCount())

	// Process events
	fmt.Println("\nProcessing events:")
	for eventQueue.HasEvents() {
		event := sim.ProcessNextEvent()
		fmt.Printf("  Processed event: %s at time %.2f from %s to %s\n",
			event.Type, event.TimeScheduled, event.SourceNode.ID, event.TargetNode.ID)

		// Process the transaction
		if event.Type == simulation.EventTypeTransactionCreated {
			tx := event.Payload.(*simulation.Transaction)
			txManager.ProcessTransaction(tx)

			// Create a transaction received event
			receivedEvent := simulation.NewEvent(
				simulation.EventTypeTransactionReceived,
				event.TimeScheduled+network.Latency(
					getNodeIndex(network, event.SourceNode),
					getNodeIndex(network, event.TargetNode),
				),
				event.SourceNode,
				event.TargetNode,
				tx,
			)
			eventQueue.AddEvent(receivedEvent)
		}
	}

	fmt.Println("\nSimulation completed")
}

// Helper function to get the index of a node in the network
func getNodeIndex(network *simulation.Network, node *simulation.Node) int {
	for i, n := range network.Nodes {
		if n == node {
			return i
		}
	}
	return -1
}
