package core

import (
	"fmt"
	"math"
	"testing"
)

// Mock WSIOracle for testing
type mockWSIOracle struct {
	prices map[string]float64
	err    error // Optional error to simulate oracle failure
}

func (m *mockWSIOracle) GetPrice(id string) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	price, ok := m.prices[id]
	if !ok {
		return 0, fmt.Errorf("mock oracle price not found for %s", id)
	}
	return price, nil
}

// --- Test WSIManager ---

func TestWSI_ParameterHandling(t *testing.T) {
	manager := NewWSIManager(1.0) // Target $1.0

	// Mock parameter and oracle
	dist1 := &mockDistribution{stdDev: 0.1} // Need StdDev for potential future use
	param1 := NewParameter("w_qUSDC", dist1)
	oracle1 := &mockWSIOracle{prices: map[string]float64{"qUSDC": 1.01}}

	dist2 := &mockDistribution{stdDev: 0.2}
	param2 := NewParameter("w_qDAI", dist2)
	oracle2 := &mockWSIOracle{prices: map[string]float64{"qDAI": 0.99}}

	t.Run("AddConstituentValid", func(t *testing.T) {
		err := manager.AddConstituent("qUSDC", param1, oracle1)
		if err != nil {
			t.Fatalf("AddConstituent failed for qUSDC: %v", err)
		}
		err = manager.AddConstituent("qDAI", param2, oracle2)
		if err != nil {
			t.Fatalf("AddConstituent failed for qDAI: %v", err)
		}

		// Verify internal state
		manager.mu.RLock()
		if len(manager.weights) != 2 {
			t.Errorf("Expected 2 constituents, got %d", len(manager.weights))
		}
		if _, ok := manager.weights["qUSDC"]; !ok {
			t.Errorf("qUSDC weight parameter not found")
		}
		if _, ok := manager.oracles["qDAI"]; !ok {
			t.Errorf("qDAI oracle not found")
		}
		manager.mu.RUnlock()
	})

	t.Run("AddConstituentDuplicate", func(t *testing.T) {
		// Need a fresh manager for this test to be independent
		freshManager := NewWSIManager(1.0)
		_ = freshManager.AddConstituent("qUSDC", param1, oracle1)
		err := freshManager.AddConstituent("qUSDC", param1, oracle1) // Add again
		if err == nil {
			t.Errorf("Expected error when adding duplicate constituent, but got nil")
		}
	})

	t.Run("AddConstituentNilParam", func(t *testing.T) {
		freshManager := NewWSIManager(1.0)
		err := freshManager.AddConstituent("qEURC", nil, oracle1)
		if err == nil {
			t.Errorf("Expected error when adding constituent with nil parameter, but got nil")
		}
	})

	t.Run("AddConstituentNilOracle", func(t *testing.T) {
		freshManager := NewWSIManager(1.0)
		param3 := NewParameter("w_qEURC", dist1)
		err := freshManager.AddConstituent("qEURC", param3, nil)
		if err == nil {
			t.Errorf("Expected error when adding constituent with nil oracle, but got nil")
		}
	})
}

func TestWSI_ValueCalculation(t *testing.T) {
	manager := NewWSIManager(1.0)

	// Setup constituents
	dist1 := &mockDistribution{stdDev: 0.1}
	param1 := NewParameter("w_qUSDC", dist1)
	param1.CurrentValue = 0.6 // Set initial weight value

	dist2 := &mockDistribution{stdDev: 0.2}
	param2 := NewParameter("w_qDAI", dist2)
	param2.CurrentValue = 0.4 // Set initial weight value

	oracle1 := &mockWSIOracle{prices: map[string]float64{"qUSDC": 1.01}}
	oracle2 := &mockWSIOracle{prices: map[string]float64{"qDAI": 0.99}}

	_ = manager.AddConstituent("qUSDC", param1, oracle1)
	_ = manager.AddConstituent("qDAI", param2, oracle2)

	t.Run("CalculateValueBasic", func(t *testing.T) {
		expectedValue := (0.6 * 1.01) + (0.4 * 0.99) // 0.606 + 0.396 = 1.002
		currentValue, err := manager.GetValue()

		if err != nil {
			t.Fatalf("GetValue failed unexpectedly: %v", err)
		}

		tolerance := 1e-9
		if math.Abs(currentValue-expectedValue) > tolerance {
			t.Errorf("Expected WSI value %.9f, got %.9f", expectedValue, currentValue)
		}
	})

	t.Run("CalculateValueOracleError", func(t *testing.T) {
		// Simulate oracle error for qDAI
		oracle2.err = fmt.Errorf("oracle timeout")
		_, err := manager.GetValue()
		if err == nil {
			t.Errorf("Expected error from GetValue due to oracle failure, but got nil")
		}
		oracle2.err = nil // Reset error for other tests
	})

	t.Run("CalculateValueInvalidPrice", func(t *testing.T) {
		// Simulate invalid price for qUSDC
		originalPrice := oracle1.prices["qUSDC"]
		oracle1.prices["qUSDC"] = -1.0 // Invalid price
		_, err := manager.GetValue()
		if err == nil {
			t.Errorf("Expected error from GetValue due to invalid price, but got nil")
		}
		oracle1.prices["qUSDC"] = originalPrice // Reset price
	})

	t.Run("CalculateValueNoConstituents", func(t *testing.T) {
		emptyManager := NewWSIManager(1.0)
		value, err := emptyManager.GetValue()
		if err != nil {
			t.Errorf("GetValue failed unexpectedly for empty manager: %v", err)
		}
		if value != 0.0 {
			t.Errorf("Expected value 0.0 for empty manager, got %f", value)
		}
	})
}

func TestWSI_PegDeviationPenalty(t *testing.T) {
	manager := NewWSIManager(1.0) // Target $1.0

	// Setup constituents (same as ValueCalculation test)
	dist1 := &mockDistribution{stdDev: 0.1}
	param1 := NewParameter("w_qUSDC", dist1)
	param1.CurrentValue = 0.6
	dist2 := &mockDistribution{stdDev: 0.2}
	param2 := NewParameter("w_qDAI", dist2)
	param2.CurrentValue = 0.4
	oracle1 := &mockWSIOracle{prices: map[string]float64{"qUSDC": 1.01}}
	oracle2 := &mockWSIOracle{prices: map[string]float64{"qDAI": 0.99}}
	_ = manager.AddConstituent("qUSDC", param1, oracle1)
	_ = manager.AddConstituent("qDAI", param2, oracle2)

	t.Run("CalculatePenaltyBasic", func(t *testing.T) {
		// From previous test, value is 1.002
		currentValue := 1.002
		targetPeg := 1.0
		expectedDeviation := currentValue - targetPeg            // 0.002
		expectedPenalty := expectedDeviation * expectedDeviation // 0.000004

		penalty, err := manager.CalculatePegPenalty()
		if err != nil {
			t.Fatalf("CalculatePegPenalty failed unexpectedly: %v", err)
		}

		tolerance := 1e-9
		if math.Abs(penalty-expectedPenalty) > tolerance {
			t.Errorf("Expected penalty %.9f, got %.9f", expectedPenalty, penalty)
		}
	})

	t.Run("CalculatePenaltyOracleError", func(t *testing.T) {
		// Simulate oracle error for qDAI
		oracle2.err = fmt.Errorf("oracle timeout")
		_, err := manager.CalculatePegPenalty()
		if err == nil {
			t.Errorf("Expected error from CalculatePegPenalty due to oracle failure, but got nil")
		}
		oracle2.err = nil // Reset error
	})

	t.Run("CalculatePenaltyNoConstituents", func(t *testing.T) {
		emptyManager := NewWSIManager(1.0)
		penalty, err := emptyManager.CalculatePegPenalty()
		if err != nil {
			t.Errorf("CalculatePegPenalty failed unexpectedly for empty manager: %v", err)
		}
		// Deviation is 0.0 - 1.0 = -1.0. Penalty = (-1.0)^2 = 1.0
		expectedPenalty := 1.0
		tolerance := 1e-9
		if math.Abs(penalty-expectedPenalty) > tolerance {
			t.Errorf("Expected penalty %.9f for empty manager, got %.9f", expectedPenalty, penalty)
		}
	})
}

// TODO: Add TestWSI_WeightUpdateTrigger (requires Hamiltonian integration)
