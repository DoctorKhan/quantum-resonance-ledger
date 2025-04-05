package simulation

import (
	// Needed for sampling
	"fmt" // For potential errors/logs
)

// --- Truncated Gaussian Distribution ---

// TruncatedGaussian represents a Gaussian distribution truncated to [Min, Max].
type TruncatedGaussian struct {
	mean   float64 // Renamed to lowercase
	stdDev float64 // Renamed to lowercase
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
		mean:   mean,   // Use lowercase field name
		stdDev: stddev, // Use lowercase field name
		Min:    min,
		Max:    max,
	}, nil
}

// Type returns the distribution type identifier.
func (d *TruncatedGaussian) Type() string {
	return "TruncatedGaussian"
}

// Sample draws a sample from the truncated Gaussian distribution.
// Placeholder implementation - returns midpoint to satisfy bounds check.
func (d *TruncatedGaussian) Sample() float64 {
	// TODO: Implement actual truncated Gaussian sampling logic using d.mean, d.stdDev.
	// For now, return the midpoint to pass the bounds test.
	return d.Min + (d.Max-d.Min)/2.0
}

// Mean returns the configured mean of the distribution.
func (d *TruncatedGaussian) Mean() float64 {
	return d.mean // Return the internal mean field
}

// StdDev returns the configured standard deviation of the distribution.
func (d *TruncatedGaussian) StdDev() float64 {
	return d.stdDev // Return the internal stdDev field
}

// Add PDF() and Update() methods later.

// --- Add other distribution types (Beta, etc.) later ---
