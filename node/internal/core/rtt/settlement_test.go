package rtt

import (
	"math"
	"testing"
)

// TestCalculateOverlap verifies the calculation of overlap between buy/sell fields.
func TestCalculateOverlap(t *testing.T) {
	// t.Skip("TDD Step 3: TestCalculateOverlap now implemented.") // Unskipped
	// TDD Step 3: First test for settlement logic. This will fail initially.

	// Arrange: Create buy and sell fields with some overlapping density
	range1 := PriceRange{Min: 100, Max: 110}
	range1Key := range1.Key()
	range2 := PriceRange{Min: 110, Max: 120} // Non-overlapping range
	range2Key := range2.Key()

	buyField := &PropensityField{
		AssetID: "ASSET_A",
		Density: map[string]float64{
			range1Key: 0.6, // Density in the target range
			range2Key: 0.2,
		},
	}
	sellField := &PropensityField{
		AssetID: "ASSET_A",
		Density: map[string]float64{
			range1Key: 0.8, // Density in the target range
			// No density in range2
		},
	}

	// Act: Calculate overlap for the specific range
	// Note: The function signature might need refinement (e.g., does it calculate
	// overlap for a specific range or overall? Plan assumes specific range for now).
	overlap, err := CalculateOverlap(buyField, sellField, range1)

	// Assert: Check for errors and expected overlap value
	if err != nil {
		t.Fatalf("CalculateOverlap returned an unexpected error: %v", err)
	}

	// Define expected overlap based on a simple calculation (e.g., product)
	// This definition will drive the implementation.
	// Let's assume overlap = buyDensity * sellDensity for the range for this test.
	expectedOverlap := 0.6 * 0.8 // 0.48

	tolerance := 1e-9
	if math.Abs(overlap-expectedOverlap) > tolerance {
		t.Errorf("Overlap mismatch: expected=%.6f, got=%.6f", expectedOverlap, overlap)
	}

	// TODO: Add tests for zero overlap, nil fields, different calculation methods.
}

// TestAttemptLocalSettlement verifies the probabilistic settlement logic.
func TestAttemptLocalSettlement(t *testing.T) {
	// t.Skip("TDD Step 3: TestAttemptLocalSettlement now implemented.") // Unskipped
	// TDD Step 3: Testing AttemptLocalSettlement

	t.Run("SettlementOccursWhenOverlapMeetsThreshold", func(t *testing.T) {
		// Arrange
		assetID := "ASSET_SETTLE_1"
		targetRange := PriceRange{Min: 100, Max: 110}
		targetRangeKey := targetRange.Key()

		buyField := &PropensityField{
			AssetID: assetID,
			Density: map[string]float64{targetRangeKey: 0.8}, // 0.8 * 0.6 = 0.48 overlap
		}
		sellField := &PropensityField{
			AssetID: assetID,
			Density: map[string]float64{targetRangeKey: 0.6},
		}

		// Create a minimal state containing the fields
		// Note: LocalNodeStateRTT needs BuyFields/SellFields maps initialized
		state := InitializeRTTState() // Use the initializer from types.go
		state.BuyFields[assetID] = buyField
		state.SellFields[assetID] = sellField
		// TODO: Add mock CUTs and NodeID to state once needed

		threshold := 0.40 // Overlap (0.48) > threshold (0.40)

		// Act
		// Note: AttemptLocalSettlement needs refinement to specify which range to check.
		// Assuming for now it checks *all* ranges or a default one implicitly.
		// We will need to adapt the function signature or logic based on this test.
		// Let's assume it finds the fields for assetID and calculates overlap implicitly for now.
		settlementRecord, err := AttemptLocalSettlement(state, assetID, threshold)

		// Assert
		if err != nil {
			t.Fatalf("AttemptLocalSettlement returned an unexpected error: %v", err)
		}
		if settlementRecord == nil {
			t.Fatal("Expected a settlement record, but got nil")
		}
		if settlementRecord.AssetID != assetID {
			t.Errorf("Settlement record has wrong AssetID: expected=%s, got=%s", assetID, settlementRecord.AssetID)
		}
		// TODO: Add more detailed checks on settlementRecord fields (Amount, PriceRange, CUTs) later.
	})

	t.Run("NoSettlementWhenOverlapBelowThreshold", func(t *testing.T) {
		// Arrange
		assetID := "ASSET_SETTLE_2"
		targetRange := PriceRange{Min: 200, Max: 210}
		targetRangeKey := targetRange.Key()

		buyField := &PropensityField{
			AssetID: assetID,
			Density: map[string]float64{targetRangeKey: 0.5}, // 0.5 * 0.7 = 0.35 overlap
		}
		sellField := &PropensityField{
			AssetID: assetID,
			Density: map[string]float64{targetRangeKey: 0.7},
		}

		state := InitializeRTTState()
		state.BuyFields[assetID] = buyField
		state.SellFields[assetID] = sellField

		threshold := 0.40 // Overlap (0.35) < threshold (0.40)

		// Act
		settlementRecord, err := AttemptLocalSettlement(state, assetID, threshold)

		// Assert
		if err != nil {
			t.Fatalf("AttemptLocalSettlement returned an unexpected error: %v", err)
		}
		if settlementRecord != nil {
			t.Errorf("Expected no settlement record, but got one: %+v", settlementRecord)
		}
	})

	t.Run("SettlementProbabilityMatchesOverlap", func(t *testing.T) {
		// Arrange
		assetID := "ASSET_SETTLE_PROB"
		targetRange := PriceRange{Min: 300, Max: 310}
		targetRangeKey := targetRange.Key()

		expectedProbability := 0.7 // Target overlap probability
		threshold := 0.1           // Must be below expectedProbability

		// Set densities to achieve the target overlap (e.g., 1.0 * 0.7 = 0.7)
		buyField := &PropensityField{
			AssetID: assetID,
			Density: map[string]float64{targetRangeKey: 1.0},
		}
		sellField := &PropensityField{
			AssetID: assetID,
			Density: map[string]float64{targetRangeKey: expectedProbability},
		}

		state := InitializeRTTState()
		state.BuyFields[assetID] = buyField
		state.SellFields[assetID] = sellField

		// Act: Run many trials
		numTrials := 10000 // Increase for higher accuracy, decrease for faster tests
		successCount := 0
		for i := 0; i < numTrials; i++ {
			// Use a copy of the state if the function modifies it, though current impl doesn't
			// stateCopy := *state // Shallow copy might be enough if fields aren't modified deeply
			settlementRecord, err := AttemptLocalSettlement(state, assetID, threshold)
			if err != nil {
				t.Fatalf("Trial %d: AttemptLocalSettlement returned an unexpected error: %v", i, err)
			}
			if settlementRecord != nil {
				successCount++
			}
		}

		// Assert: Check if observed frequency is close to expected probability
		observedFrequency := float64(successCount) / float64(numTrials)
		tolerance := 0.05 // Adjust tolerance based on numTrials (higher N -> lower tolerance)

		if math.Abs(observedFrequency-expectedProbability) > tolerance {
			t.Errorf("Observed settlement frequency (%.4f) is outside tolerance (%.2f) of expected probability (%.4f) after %d trials",
				observedFrequency, tolerance, expectedProbability, numTrials)
		}
	})

	// TODO: Add tests with CUT management integrated
	// TODO: Add tests for Q imbalance updates (state modification)
}

// TODO: Add helper functions to create test states if needed.
