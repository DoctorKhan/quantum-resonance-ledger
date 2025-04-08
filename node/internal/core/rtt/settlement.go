package rtt

import (
	"fmt"
	"math/rand" // Keep for placeholder RecordID for now
	"strconv"
	"strings"
	// "time" // Potentially needed for timestamps
)

// Helper to parse PriceRange key string (inverse of Key())
// WARNING: Brittle, assumes specific "%.2f-%.2f" format from Key().
func ParsePriceRangeKey(key string) (PriceRange, error) {
	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return PriceRange{}, fmt.Errorf("invalid price range key format: %s", key)
	}
	min, errMin := strconv.ParseFloat(parts[0], 64)
	max, errMax := strconv.ParseFloat(parts[1], 64)
	if errMin != nil || errMax != nil {
		return PriceRange{}, fmt.Errorf("error parsing floats from key '%s': %v, %v", key, errMin, errMax)
	}
	return PriceRange{Min: min, Max: max}, nil
}

// CalculateOverlap computes the overlap integral or metric between buy and sell propensity fields.
// The exact calculation depends heavily on the PropensityField representation.
// Returns a metric representing the potential for a trade in a given range/overall.
func CalculateOverlap(buyField, sellField *PropensityField, priceRange PriceRange) (float64, error) {
	// Basic implementation for TDD Step 3.
	// Calculates overlap simply as the product of densities in the specific priceRange key.
	// Assumes the map[string]float64 representation in PropensityField.

	if buyField == nil || sellField == nil {
		return 0.0, fmt.Errorf("cannot calculate overlap with nil fields")
	}
	// TODO: Add validation for priceRange if needed (e.g., Min <= Max)

	key := priceRange.Key()
	buyDensity := 0.0
	sellDensity := 0.0

	if buyField.Density != nil {
		buyDensity = buyField.Density[key] // Defaults to 0 if key not present
	}
	if sellField.Density != nil {
		sellDensity = sellField.Density[key] // Defaults to 0 if key not present
	}

	// Simple overlap calculation: product of densities in the range.
	// More sophisticated methods (integration, min function, etc.) might be needed later.
	overlap := buyDensity * sellDensity

	// Ensure overlap is not negative (densities should be non-negative)
	if overlap < 0 {
		// This case implies negative density, which should ideally be prevented earlier.
		// Log a warning or return an error if necessary.
		// For now, clamp to 0.
		overlap = 0.0
		// Consider returning an error: fmt.Errorf("negative density detected during overlap calculation")
	}

	return overlap, nil
}

// AttemptLocalSettlement evaluates local fields and probabilistically initiates settlement.
// It uses the calculated overlap and a threshold (theta_trade_threshold) to decide.
func AttemptLocalSettlement(state *LocalNodeStateRTT, assetID string, tradeThreshold float64) (*SettlementRecord, error) {
	// TODO: Implement the full logic:
	// 1. Get relevant buy and sell fields from state for the assetID.
	// 2. Define the relevant PriceRange(s) to check for overlap (might be the whole field initially).
	// 3. Call CalculateOverlap for those fields/ranges.
	// 4. Compare overlap metric to tradeThreshold.
	// 5. If overlap >= tradeThreshold:
	//    a. Determine settlement amount (proportional to overlap?).
	//    b. Determine settlement price range.
	//    c. Check/select available CUTs from state.
	//    d. Generate a SettlementRecord.
	//    e. Mark used CUTs in the state.
	//    f. Update local balances/state (potentially reflected in Q).
	//    g. Return the SettlementRecord.
	// 6. If overlap < tradeThreshold, return nil (no settlement).

	if state == nil {
		return nil, fmt.Errorf("cannot attempt settlement with nil state")
	}

	// --- Actual Logic (Minimal for TDD Step 3) ---

	// 1. Get fields from state
	buyField, okBuy := state.BuyFields[assetID]
	sellField, okSell := state.SellFields[assetID]
	if !okBuy || !okSell || buyField == nil || sellField == nil {
		return nil, nil // No fields for this asset, no settlement possible
	}
	if buyField.Density == nil {
		return nil, nil // No buy density defined
	}

	// 2. Iterate through price ranges in the buy field to find potential overlap
	// TODO: This is inefficient. A better approach might involve iterating common keys
	// or using a more structured field representation. Sticking to minimal change for TDD.
	for key := range buyField.Density {
		priceRange, err := ParsePriceRangeKey(key)
		if err != nil {
			// Log error? Skip invalid key? For now, skip.
			// fmt.Printf("Warning: Skipping invalid price range key '%s': %v\n", key, err)
			continue
		}

		// 3. Calculate overlap for this specific range
		overlap, err := CalculateOverlap(buyField, sellField, priceRange)
		if err != nil {
			// Log error? Return error? For now, continue checking other ranges.
			// fmt.Printf("Warning: Error calculating overlap for range %s: %v\n", key, err)
			continue
		}

		// 4. Check if overlap meets minimum threshold
		if overlap >= tradeThreshold {
			// 5. Probabilistic check: Settle with probability = overlap
			// Ensure overlap is treated as probability (e.g., clamp between 0 and 1 if necessary,
			// although CalculateOverlap should ideally ensure non-negativity).
			settlementProbability := overlap
			// Clamp probability just in case (optional, depends on CalculateOverlap guarantees)
			// if settlementProbability > 1.0 { settlementProbability = 1.0 }
			// if settlementProbability < 0.0 { settlementProbability = 0.0 } // Already handled in CalculateOverlap

			if rand.Float64() < settlementProbability {
				// Settlement occurs!
				// Create placeholder settlement record.
				// TODO: Refine Amount, PriceRange representation, CUTs, NodeID, Timestamp
				record := &SettlementRecord{
					RecordID:   fmt.Sprintf("settle_%d", rand.Int()), // Placeholder ID
					AssetID:    assetID,
					Amount:     overlap * 10.0, // Placeholder amount logic (proportional to overlap?)
					PriceRange: key,            // Use the key string for now
					Timestamp:  0,              // TODO: Use time.Now().UnixNano()
					NodeID:     "TODO",         // TODO: Get local node ID from state/config
					// InvolvedCUTs: []CUT_ID{"cut1", "cut2"}, // Placeholder CUTs
				}
				// TODO: Add logic to update state (mark CUTs, update Q)
				return record, nil // Return the first settlement found
			}
			// else: Overlap met threshold, but probabilistic check failed. Continue loop.
		}
		// else: Overlap below threshold. Continue loop.
	}

	// 6. If loop completes, no range met the threshold
	return nil, nil // No settlement occurred
}
