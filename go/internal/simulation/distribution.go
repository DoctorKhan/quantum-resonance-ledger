package simulation

import (
	// Needed for sampling
	"fmt"  // For potential errors/logs
	"math" // Needed for IsNaN etc.
)

// --- Truncated Gaussian Distribution ---

// TruncatedGaussian represents a Gaussian distribution truncated to [Min, Max].
type TruncatedGaussian struct {
	Mean   float64
	StdDev float64
	Min    float64
	Max    float64
}

// NewTruncatedGaussian creates a new TruncatedGaussian distribution.
// Performs basic validation.
func NewTruncatedGaussian(mean, stddev, min, max float64) (*TruncatedGaussian, error) {
	if stddev <= 0 {
		return nil, fmt.Errorf("standard deviation must be positive (got %f)", stddev)
	}
	if min >= max {
		return nil, fmt.Errorf("min must be less than max (got min=%f, max=%f)", min, max)
	}
	// Add checks for NaN/Inf if necessary
	return &TruncatedGaussian{
		Mean:   mean,
		StdDev: stddev,
		Min:    min,
		Max:    max,
	}, nil
}

// Type returns the distribution type identifier.
func (d *TruncatedGaussian) Type() string {
	return "TruncatedGaussian"
}

// Sample draws a sample from the truncated Gaussian distribution.
// Placeholder implementation - returns a fixed value for initial test failure.
func (d *TruncatedGaussian) Sample() float64 {
	// TODO: Implement actual truncated sampling logic
	// For now, return a value guaranteed to be in range if Min/Max are valid,
	// but not necessarily following the distribution, or return something fixed.
	// Returning Min ensures it's in range but doesn't test sampling logic.
	// Let's return a value likely outside simple bounds if possible, or NaN.
	return math.NaN() // This will likely fail the bounds check in the test.
}

// Add PDF() and Update() methods later.

// --- Add other distribution types (Beta, etc.) later ---
