package simulation

import (
	"testing"
)

// Mock Distribution for testing Parameter creation
type mockDistribution struct {
	distType string
}

func (m *mockDistribution) Type() string { return m.distType }

// Sample returns a dummy value for the mock.
func (m *mockDistribution) Sample() float64 { return 0.0 }

// Add other methods if needed by Parameter constructor later

func TestParameterCreation(t *testing.T) {
	// Test case 1: Create parameter with a mock distribution
	t.Run("CreateParameter", func(t *testing.T) {
		paramName := "NetworkLatencyFactor"
		mockDist := &mockDistribution{distType: "TruncatedGaussian"}
		param := NewParameter(paramName, mockDist)

		if param == nil {
			t.Fatalf("NewParameter returned nil")
		}
		if param.Name != paramName {
			t.Errorf("Expected parameter name '%s', got '%s'", paramName, param.Name)
		}
		if param.Distribution == nil {
			t.Fatalf("Parameter distribution is nil")
		}
		if param.Distribution.Type() != mockDist.Type() {
			t.Errorf("Expected distribution type '%s', got '%s'", mockDist.Type(), param.Distribution.Type())
		}
		// We don't compare param.Distribution == mockDist directly for interface values
	})

	// Test case 2: Edge case - Empty name
	t.Run("EmptyNameParameter", func(t *testing.T) {
		mockDist := &mockDistribution{distType: "Beta"}
		param := NewParameter("", mockDist) // Assume empty name is allowed for now

		if param == nil {
			t.Fatalf("NewParameter returned nil for empty name")
		}
		if param.Name != "" {
			t.Errorf("Expected empty parameter name, got '%s'", param.Name)
		}
		if param.Distribution == nil {
			t.Errorf("Parameter distribution is nil for empty name case")
		}
	})

	// Test case 3: Edge case - Nil distribution
	t.Run("NilDistribution", func(t *testing.T) {
		paramName := "ConsensusThreshold"
		// Behavior with nil distribution depends on design (error? default? panic?)
		// Let's assume for now it might create the parameter but with nil distribution.
		param := NewParameter(paramName, nil)

		if param == nil {
			t.Fatalf("NewParameter returned nil for nil distribution")
		}
		if param.Name != paramName {
			t.Errorf("Expected parameter name '%s', got '%s'", paramName, param.Name)
		}
		// Depending on design, this might be expected or an error state
		if param.Distribution != nil {
			t.Errorf("Expected nil distribution, but got %v", param.Distribution)
		} else {
			t.Logf("Parameter created with nil distribution as expected (for now).")
		}
	})
}
