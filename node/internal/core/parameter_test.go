package core

import (
	"testing"
)

// Mock Distribution for testing Parameter creation
type mockDistribution struct {
	distType string
	mean     float64 // Add fields to store mock values
	stdDev   float64
}

func (m *mockDistribution) Type() string { return m.distType }

// Sample returns a dummy value for the mock.
func (m *mockDistribution) Sample() float64 { return 0.0 }

// Mean returns the mock mean value.
func (m *mockDistribution) Mean() float64 { return m.mean }

// StdDev returns the mock stddev value.
func (m *mockDistribution) StdDev() float64 { return m.stdDev }

// Add other methods if needed later

func TestParameterCreation(t *testing.T) {
	// Test case 1: Create parameter with a mock distribution
	t.Run("CreateParameter", func(t *testing.T) { // Keep existing tests
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

func TestUncertaintyRelationHandling(t *testing.T) {
	manager := NewParameterManager()

	// Create mock parameters with specific StdDevs
	dist1 := &mockDistribution{distType: "Mock1", stdDev: 2.0}
	param1 := NewParameter("ParamA", dist1)
	err := manager.AddParameter(param1)
	if err != nil {
		t.Fatalf("Failed to add param1: %v", err)
	}

	dist2 := &mockDistribution{distType: "Mock2", stdDev: 3.0}
	param2 := NewParameter("ParamB", dist2)
	err = manager.AddParameter(param2)
	if err != nil {
		t.Fatalf("Failed to add param2: %v", err)
	}

	dist3 := &mockDistribution{distType: "Mock3", stdDev: 0.5}
	param3 := NewParameter("ParamC", dist3)
	err = manager.AddParameter(param3)
	if err != nil {
		t.Fatalf("Failed to add param3: %v", err)
	}

	// Test case 1: Valid relation (2.0 * 3.0 >= 5.0)
	t.Run("ValidRelation", func(t *testing.T) {
		relation, err := NewUncertaintyRelation(param1, param2, 5.0)
		if err != nil {
			t.Fatalf("Failed to create relation: %v", err)
		}
		err = manager.AddUncertaintyRelation(relation)
		if err != nil {
			t.Fatalf("Failed to add relation: %v", err)
		}

		allValid, violations, err := manager.ValidateAllUncertaintyRelations()
		if err != nil {
			t.Errorf("Validation failed unexpectedly: %v", err)
		}
		if !allValid {
			t.Errorf("Expected all relations to be valid, but got violations: %v", violations)
		}
		if len(violations) != 0 {
			t.Errorf("Expected 0 violations, got %d", len(violations))
		}

		// Clean up for next test - remove the relation
		manager.UncertaintyRelations = manager.UncertaintyRelations[:len(manager.UncertaintyRelations)-1]
	})

	// Test case 2: Invalid relation (2.0 * 0.5 < 1.5)
	t.Run("InvalidRelation", func(t *testing.T) {
		relation, err := NewUncertaintyRelation(param1, param3, 1.5) // 2.0 * 0.5 = 1.0, which is < 1.5
		if err != nil {
			t.Fatalf("Failed to create relation: %v", err)
		}
		err = manager.AddUncertaintyRelation(relation)
		if err != nil {
			t.Fatalf("Failed to add relation: %v", err)
		}

		allValid, violations, err := manager.ValidateAllUncertaintyRelations()
		if err != nil {
			t.Errorf("Validation failed unexpectedly: %v", err)
		}
		if allValid {
			t.Errorf("Expected relation to be invalid, but validation passed")
		}
		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d", len(violations))
		}
		if len(violations) > 0 && violations[0] != relation {
			t.Errorf("Violation list does not contain the expected relation")
		}
		// Clean up for next test
		manager.UncertaintyRelations = manager.UncertaintyRelations[:len(manager.UncertaintyRelations)-1]
	})

	// Test case 3: Multiple relations (one valid, one invalid)
	t.Run("MultipleRelationsMixed", func(t *testing.T) {
		validRel, _ := NewUncertaintyRelation(param1, param2, 5.0)   // 2*3=6 >= 5 (Valid)
		invalidRel, _ := NewUncertaintyRelation(param1, param3, 1.5) // 2*0.5=1 < 1.5 (Invalid)
		_ = manager.AddUncertaintyRelation(validRel)
		_ = manager.AddUncertaintyRelation(invalidRel)

		allValid, violations, err := manager.ValidateAllUncertaintyRelations()
		if err != nil {
			t.Errorf("Validation failed unexpectedly: %v", err)
		}
		if allValid {
			t.Errorf("Expected validation to fail due to one invalid relation")
		}
		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d", len(violations))
		}
		if len(violations) > 0 && violations[0] != invalidRel {
			t.Errorf("Violation list does not contain the expected invalid relation")
		}
		// Clean up
		manager.UncertaintyRelations = manager.UncertaintyRelations[:0] // Clear all relations
	})

	// Test case 4: Error during validation (e.g., nil distribution)
	t.Run("ValidationError", func(t *testing.T) {
		// Create a parameter without assigning a distribution
		paramErr := &Parameter{Name: "ParamErr"}
		_ = manager.AddParameter(paramErr)

		relation, err := NewUncertaintyRelation(param1, paramErr, 1.0)
		if err != nil {
			t.Fatalf("Failed to create relation: %v", err)
		}
		err = manager.AddUncertaintyRelation(relation)
		if err != nil {
			t.Fatalf("Failed to add relation: %v", err)
		}

		allValid, violations, err := manager.ValidateAllUncertaintyRelations()
		if err == nil {
			t.Errorf("Expected an error during validation due to nil distribution, but got nil")
		}
		if allValid {
			t.Errorf("Expected validation to be false when an error occurs")
		}
		if len(violations) != 1 { // Should still report the relation causing the error as a violation
			t.Errorf("Expected 1 violation (due to error), got %d", len(violations))
		}

		// Clean up
		delete(manager.Parameters, "ParamErr")
		manager.UncertaintyRelations = manager.UncertaintyRelations[:0]
	})

	// Test case 5: Add relation with non-managed parameter
	t.Run("AddRelationWithUnmanagedParam", func(t *testing.T) {
		unmanagedParam := NewParameter("Unmanaged", &mockDistribution{stdDev: 1.0})
		relation, err := NewUncertaintyRelation(param1, unmanagedParam, 1.0)
		if err != nil {
			t.Fatalf("Failed to create relation: %v", err)
		}
		err = manager.AddUncertaintyRelation(relation)
		if err == nil {
			t.Errorf("Expected error when adding relation with unmanaged parameter, but got nil")
		}
	})
}

func TestParameterUpdateRule(t *testing.T) {
	// Use a TruncatedGaussian for testing bounds later
	dist, err := NewTruncatedGaussian(5.0, 1.0, 0.0, 10.0)
	if err != nil {
		t.Fatalf("Failed to create distribution: %v", err)
	}
	param := NewParameter("TestParam", dist)
	initialValue := 5.0
	param.CurrentValue = initialValue // Initialize the value

	// Test case 1: Basic update (positive gradient -> decrease value)
	t.Run("BasicUpdatePositiveGradient", func(t *testing.T) {
		gradient := 2.0
		eta := 0.1
		dt := 1.0
		expectedChange := -eta * dt * gradient         // -0.1 * 1.0 * 2.0 = -0.2
		expectedValue := initialValue + expectedChange // 5.0 - 0.2 = 4.8

		err := param.Update(gradient, eta, dt)
		if err != nil {
			t.Errorf("Update failed unexpectedly: %v", err)
		}

		// This assertion will fail until Update logic is implemented
		if param.CurrentValue != expectedValue {
			t.Errorf("Expected CurrentValue %.4f, got %.4f", expectedValue, param.CurrentValue)
		}

		// Reset value for next test
		param.CurrentValue = initialValue
	})

	// Test case 2: Basic update (negative gradient -> increase value)
	t.Run("BasicUpdateNegativeGradient", func(t *testing.T) {
		gradient := -1.5
		eta := 0.2
		dt := 0.5
		expectedChange := -eta * dt * gradient         // -0.2 * 0.5 * (-1.5) = 0.15
		expectedValue := initialValue + expectedChange // 5.0 + 0.15 = 5.15

		err := param.Update(gradient, eta, dt)
		if err != nil {
			t.Errorf("Update failed unexpectedly: %v", err)
		}

		// This assertion will fail until Update logic is implemented
		if param.CurrentValue != expectedValue {
			t.Errorf("Expected CurrentValue %.4f, got %.4f", expectedValue, param.CurrentValue)
		}

		// Reset value
		param.CurrentValue = initialValue
	})

	// Test case 3: Update with nil distribution (should error)
	t.Run("UpdateNilDistribution", func(t *testing.T) {
		paramNoDist := NewParameter("NoDistParam", nil)
		paramNoDist.CurrentValue = 1.0

		err := paramNoDist.Update(1.0, 0.1, 1.0)
		if err == nil {
			t.Errorf("Expected an error when updating parameter with nil distribution, but got nil")
		}
	})

	// TODO: Add tests for bounds checking once implemented in Update
}
