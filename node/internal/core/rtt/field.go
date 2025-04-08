package rtt

import "fmt" // Added for error formatting

// PerturbPropensityField modifies a given propensity field based on new trading intent.
// The exact implementation depends heavily on the chosen representation for PropensityField.
// This is a placeholder signature.
func PerturbPropensityField(field *PropensityField, assetID string, priceRange PriceRange, magnitude float64, isBuyField bool) error {
	// Basic implementation for TDD Step 2 to pass the first test.
	// Assumes direct addition to the density map based on PriceRange key.
	// Ignores assetID and isBuyField for now. Adds basic validation.

	if field == nil {
		return fmt.Errorf("cannot perturb nil field")
	}
	// Basic validation for magnitude (can be expanded)
	if magnitude < 0 {
		return fmt.Errorf("perturbation magnitude cannot be negative: %f", magnitude)
	}

	// Ensure the density map is initialized
	if field.Density == nil {
		field.Density = make(map[string]float64)
	}

	// Get the key for the price range
	key := priceRange.Key()

	// Add the magnitude to the existing density for this key
	field.Density[key] += magnitude

	// TODO: Add proper error handling (e.g., invalid price range).
	// TODO: Implement more sophisticated density update logic if needed.
	// TODO: Use assetID and isBuyField.
	return nil
}

// TODO: Add functions for calculating overlap, smoothing, etc. as needed.
