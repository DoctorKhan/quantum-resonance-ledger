package main

import (
	"fmt"
	"os"

	"github.com/khan/qrib/pkg/paramspace"
)

func main() {
	fmt.Println("Quantum-Inspired Blockchain Framework - Parameter Space Example")
	fmt.Println("===========================================================")

	// Create a new parameter space
	space := paramspace.NewParameterSpace()
	fmt.Println("Created a new parameter space")

	// Create parameters
	blockSizeParam, err := paramspace.NewParameter("block_size", 1.0, 10.0)
	if err != nil {
		fmt.Printf("Error creating parameter: %v\n", err)
		os.Exit(1)
	}

	txThroughputParam, err := paramspace.NewParameter("tx_throughput", 100.0, 1000.0)
	if err != nil {
		fmt.Printf("Error creating parameter: %v\n", err)
		os.Exit(1)
	}

	latencyParam, err := paramspace.NewParameter("latency", 0.1, 5.0)
	if err != nil {
		fmt.Printf("Error creating parameter: %v\n", err)
		os.Exit(1)
	}

	// Add parameters to the space
	err = space.AddParameter(blockSizeParam)
	if err != nil {
		fmt.Printf("Error adding parameter: %v\n", err)
		os.Exit(1)
	}

	err = space.AddParameter(txThroughputParam)
	if err != nil {
		fmt.Printf("Error adding parameter: %v\n", err)
		os.Exit(1)
	}

	err = space.AddParameter(latencyParam)
	if err != nil {
		fmt.Printf("Error adding parameter: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added parameters: block_size[%.1f, %.1f], tx_throughput[%.1f, %.1f], latency[%.1f, %.1f]\n",
		blockSizeParam.Min(), blockSizeParam.Max(),
		txThroughputParam.Min(), txThroughputParam.Max(),
		latencyParam.Min(), latencyParam.Max())

	// Create uncertainty relations
	// Larger block sizes increase latency
	blockSizeLatencyRelation, err := paramspace.NewUncertaintyRelation(blockSizeParam, latencyParam, 1.0)
	if err != nil {
		fmt.Printf("Error creating relation: %v\n", err)
		os.Exit(1)
	}

	// Higher throughput decreases latency (trade-off)
	throughputLatencyRelation, err := paramspace.NewUncertaintyRelation(txThroughputParam, latencyParam, 50.0)
	if err != nil {
		fmt.Printf("Error creating relation: %v\n", err)
		os.Exit(1)
	}

	// Add relations to the space
	err = space.AddRelation(blockSizeLatencyRelation)
	if err != nil {
		fmt.Printf("Error adding relation: %v\n", err)
		os.Exit(1)
	}

	err = space.AddRelation(throughputLatencyRelation)
	if err != nil {
		fmt.Printf("Error adding relation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Added uncertainty relations between parameters")

	// Create distributions for parameters
	blockSizeDist, err := paramspace.NewUniformDistribution(blockSizeParam)
	if err != nil {
		fmt.Printf("Error creating distribution: %v\n", err)
		os.Exit(1)
	}

	txThroughputDist, err := paramspace.NewNormalDistribution(txThroughputParam, 500.0, 100.0)
	if err != nil {
		fmt.Printf("Error creating distribution: %v\n", err)
		os.Exit(1)
	}

	// For uniform distribution on [1,10], the uncertainty is (10-1)/sqrt(12) ≈ 2.6
	// For normal distribution with stddev=100, the uncertainty is 100
	// So for the latency, we need uncertainty > 1.0/2.6 ≈ 0.38 and > 50.0/100 = 0.5
	// So we need uncertainty > 0.5, which means stddev > 0.5
	latencyDist, err := paramspace.NewNormalDistribution(latencyParam, 2.0, 0.6)
	if err != nil {
		fmt.Printf("Error creating distribution: %v\n", err)
		os.Exit(1)
	}

	// Set distributions in the space
	err = space.SetDistribution(blockSizeParam, blockSizeDist)
	if err != nil {
		fmt.Printf("Error setting distribution: %v\n", err)
		os.Exit(1)
	}

	err = space.SetDistribution(txThroughputParam, txThroughputDist)
	if err != nil {
		fmt.Printf("Error setting distribution: %v\n", err)
		os.Exit(1)
	}

	err = space.SetDistribution(latencyParam, latencyDist)
	if err != nil {
		fmt.Printf("Error setting distribution: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Set probability distributions for all parameters")

	// Validate uncertainty relations
	valid, violations := space.ValidateUncertaintyRelations()
	if valid {
		fmt.Println("All uncertainty relations are satisfied!")
	} else {
		fmt.Printf("Some uncertainty relations are violated: %d violations\n", len(violations))
		for i, relation := range violations {
			fmt.Printf("Violation %d: Relation between %s and %s (constant: %.2f)\n",
				i+1, relation.Parameter1().Name(), relation.Parameter2().Name(), relation.Constant())
		}
	}

	// Now let's try with a distribution that violates the relations
	fmt.Println("\nChanging latency distribution to have smaller uncertainty...")

	// Create a distribution with smaller uncertainty
	latencyDist2, err := paramspace.NewNormalDistribution(latencyParam, 2.0, 0.3)
	if err != nil {
		fmt.Printf("Error creating distribution: %v\n", err)
		os.Exit(1)
	}

	// Set the new distribution
	err = space.SetDistribution(latencyParam, latencyDist2)
	if err != nil {
		fmt.Printf("Error setting distribution: %v\n", err)
		os.Exit(1)
	}

	// Validate again
	valid, violations = space.ValidateUncertaintyRelations()
	if valid {
		fmt.Println("All uncertainty relations are satisfied!")
	} else {
		fmt.Printf("Some uncertainty relations are violated: %d violations\n", len(violations))
		for i, relation := range violations {
			fmt.Printf("Violation %d: Relation between %s and %s (constant: %.2f)\n",
				i+1, relation.Parameter1().Name(), relation.Parameter2().Name(), relation.Constant())
		}
	}
}
