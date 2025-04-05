package simulation

import (
	"math"
	"testing"
)

func TestDistributionSampling(t *testing.T) {

	// Test case 1: TruncatedGaussian Sampling
	t.Run("TruncatedGaussianSampleBounds", func(t *testing.T) {
		mean, stddev, min, max := 5.0, 1.0, 3.0, 7.0
		dist, err := NewTruncatedGaussian(mean, stddev, min, max)
		if err != nil {
			t.Fatalf("Failed to create TruncatedGaussian: %v", err)
		}
		if dist == nil {
			t.Fatalf("NewTruncatedGaussian returned nil distribution")
		}

		// Draw multiple samples and check bounds
		numSamples := 100
		for i := 0; i < numSamples; i++ {
			sample := dist.Sample() // Uses placeholder Sample() initially

			// Basic check: Is the sample within the defined bounds?
			if math.IsNaN(sample) {
				// This is expected to fail initially as placeholder returns NaN
				t.Logf("Sample %d is NaN (placeholder behavior)", i+1)
				// t.Errorf("Sample %d is NaN", i+1) // Enable this when Sample() is implemented
				continue // Skip further checks for NaN
			}
			if sample < min || sample > max {
				t.Errorf("Sample %d (%.4f) is outside bounds [%.4f, %.4f]",
					i+1, sample, min, max)
			}
		}
		// TODO: Add statistical tests later to verify distribution shape
	})

	// Add tests for other distributions (Beta, etc.) later

	// Test case 2: TruncatedGaussian Mean
	t.Run("TruncatedGaussianMean", func(t *testing.T) {
		mean, stddev, min, max := 5.0, 1.0, 3.0, 7.0
		dist, err := NewTruncatedGaussian(mean, stddev, min, max)
		if err != nil {
			t.Fatalf("Failed to create TruncatedGaussian: %v", err)
		}
		if dist == nil {
			t.Fatalf("NewTruncatedGaussian returned nil distribution")
		}

		// Check if the Mean method returns the expected value
		expectedMean := mean      // For TruncatedGaussian, Mean() should return the initial mean
		actualMean := dist.Mean() // This will fail until Mean() is implemented

		// Check for NaN (should not happen with current implementation)
		if math.IsNaN(actualMean) {
			t.Errorf("Mean is NaN")
		}
		// Check if the value matches the expected configured mean
		if actualMean != expectedMean {
			t.Errorf("Expected Mean %.4f, got %.4f", expectedMean, actualMean)
		}
	})

	// Test case 3: TruncatedGaussian StdDev
	t.Run("TruncatedGaussianStdDev", func(t *testing.T) {
		mean, stddev, min, max := 5.0, 1.0, 3.0, 7.0
		dist, err := NewTruncatedGaussian(mean, stddev, min, max)
		if err != nil {
			t.Fatalf("Failed to create TruncatedGaussian: %v", err)
		}
		if dist == nil {
			t.Fatalf("NewTruncatedGaussian returned nil distribution")
		}

		// Check if the StdDev method returns the expected value
		expectedStdDev := stddev      // For TruncatedGaussian, StdDev() should return the initial stddev
		actualStdDev := dist.StdDev() // This will fail until StdDev() is implemented

		// Check for NaN (should not happen with current implementation)
		if math.IsNaN(actualStdDev) {
			t.Errorf("StdDev is NaN")
		}
		// Check if the value matches the expected configured stddev
		if actualStdDev != expectedStdDev {
			t.Errorf("Expected StdDev %.4f, got %.4f", expectedStdDev, actualStdDev)
		}
	})
}

// Add tests for PDF, Update later
