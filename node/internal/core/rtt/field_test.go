package rtt

import (
	"math" // For floating point comparison tolerance
	"testing"
)

// No helper needed now, using actual PropensityField from types.go

// --- Test Cases ---

// TestPerturbPropensityField_IncreasesDensity verifies that perturbing increases density in the target range.
func TestPerturbPropensityField_IncreasesDensity(t *testing.T) {
	// TDD Step 2: First test for PerturbPropensityField.
	// This test WILL FAIL until PerturbPropensityField is implemented in field.go.

	// Arrange: Create a test field and define perturbation parameters
	initialField := &PropensityField{
		AssetID: "ASSET_A",
		Density: make(map[string]float64), // Using the map defined in types.go
	}
	targetRange := PriceRange{Min: 100.00, Max: 110.00}
	targetRangeKey := targetRange.Key() // Use the Key() method
	initialDensity := 0.5
	initialField.Density[targetRangeKey] = initialDensity

	perturbMagnitude := 0.2
	isBuy := true

	// Act
	// Call the function under test (which is currently a placeholder)
	err := PerturbPropensityField(initialField, initialField.AssetID, targetRange, perturbMagnitude, isBuy)

	// Assert - Initial basic assertion: function should not error for valid inputs (yet)
	if err != nil {
		t.Fatalf("PerturbPropensityField returned an unexpected error: %v", err)
	}

	// Assert - Check density change (this part will fail until implementation)
	// Density should increase. The exact calculation might be complex (e.g., proportional),
	// but for a simple test, let's assume direct addition for now.
	expectedDensity := initialDensity + perturbMagnitude
	actualDensity := initialField.Density[targetRangeKey] // Direct access for testing

	// Use a tolerance for floating-point comparisons
	tolerance := 1e-9
	if math.Abs(actualDensity-expectedDensity) > tolerance {
		t.Errorf("Density mismatch after perturbation: expected=%.6f, got=%.6f", expectedDensity, actualDensity)
	}

	// --- End of Test ---
}

// TestPerturbPropensityField_Accumulates verifies that multiple perturbations sum correctly.
func TestPerturbPropensityField_Accumulates(t *testing.T) {
	// Arrange
	initialField := &PropensityField{
		AssetID: "ASSET_B",
		Density: make(map[string]float64),
	}
	targetRange := PriceRange{Min: 50.00, Max: 55.00}
	targetRangeKey := targetRange.Key()
	initialDensity := 0.1
	initialField.Density[targetRangeKey] = initialDensity

	perturb1 := 0.15
	perturb2 := 0.25
	isBuy := false

	// Act
	err1 := PerturbPropensityField(initialField, initialField.AssetID, targetRange, perturb1, isBuy)
	if err1 != nil {
		t.Fatalf("First PerturbPropensityField failed: %v", err1)
	}
	err2 := PerturbPropensityField(initialField, initialField.AssetID, targetRange, perturb2, isBuy)
	if err2 != nil {
		t.Fatalf("Second PerturbPropensityField failed: %v", err2)
	}

	// Assert
	expectedDensity := initialDensity + perturb1 + perturb2
	actualDensity := initialField.Density[targetRangeKey]
	tolerance := 1e-9

	if math.Abs(actualDensity-expectedDensity) > tolerance {
		t.Errorf("Density mismatch after accumulation: expected=%.6f, got=%.6f", expectedDensity, actualDensity)
	}
}

// TestPerturbPropensityField_EdgeCases tests various edge cases and invalid inputs.
func TestPerturbPropensityField_EdgeCases(t *testing.T) {

	t.Run("ZeroMagnitude", func(t *testing.T) {
		// Arrange
		initialField := &PropensityField{
			AssetID: "ASSET_C",
			Density: make(map[string]float64),
		}
		targetRange := PriceRange{Min: 200.00, Max: 210.00}
		targetRangeKey := targetRange.Key()
		initialDensity := 1.0
		initialField.Density[targetRangeKey] = initialDensity

		// Act
		err := PerturbPropensityField(initialField, initialField.AssetID, targetRange, 0.0, true)

		// Assert
		if err != nil {
			t.Errorf("Perturbation with zero magnitude returned unexpected error: %v", err)
		}
		actualDensity := initialField.Density[targetRangeKey]
		tolerance := 1e-9
		if math.Abs(actualDensity-initialDensity) > tolerance {
			t.Errorf("Density changed unexpectedly for zero magnitude: expected=%.6f, got=%.6f", initialDensity, actualDensity)
		}
	})

	t.Run("NegativeMagnitude", func(t *testing.T) {
		// Arrange
		initialField := &PropensityField{
			AssetID: "ASSET_D",
			Density: make(map[string]float64),
		}
		targetRange := PriceRange{Min: 300.00, Max: 310.00}

		// Act
		err := PerturbPropensityField(initialField, initialField.AssetID, targetRange, -0.1, true)

		// Assert
		if err == nil {
			t.Error("Perturbation with negative magnitude did not return an error")
		}
		// Optional: Check the error message contains expected text
		// if !strings.Contains(err.Error(), "magnitude cannot be negative") {
		// 	t.Errorf("Error message mismatch for negative magnitude: got %q", err.Error())
		// }
	})

	t.Run("NilField", func(t *testing.T) {
		// Arrange
		var nilField *PropensityField = nil
		targetRange := PriceRange{Min: 400.00, Max: 410.00}

		// Act
		err := PerturbPropensityField(nilField, "ASSET_E", targetRange, 0.1, true)

		// Assert
		if err == nil {
			t.Error("Perturbation with nil field did not return an error")
		}
		// Optional: Check error message
		// if !strings.Contains(err.Error(), "cannot perturb nil field") {
		//  t.Errorf("Error message mismatch for nil field: got %q", err.Error())
		// }
	})

	// TODO: Add test for invalid PriceRange (e.g., Min > Max) once validation is added.
}
