package paramspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUniformDistribution tests the creation and properties of a uniform distribution
func TestUniformDistribution(t *testing.T) {
	// Create a parameter with min=0, max=10
	param, err := NewParameter("test_param", 0.0, 10.0)
	assert.NoError(t, err)

	// Create a uniform distribution for the parameter
	dist, err := NewUniformDistribution(param)
	assert.NoError(t, err)

	// Test that the distribution type is correct
	assert.Equal(t, DistributionTypeUniform, dist.Type())

	// Test the PDF values at different points
	// For a uniform distribution on [0,10], the PDF should be 0.1 everywhere in the range
	assert.InDelta(t, 0.1, dist.PDF(0.0), 0.0001)
	assert.InDelta(t, 0.1, dist.PDF(5.0), 0.0001)
	assert.InDelta(t, 0.1, dist.PDF(10.0), 0.0001)

	// Test that the PDF is 0 outside the range
	assert.InDelta(t, 0.0, dist.PDF(-1.0), 0.0001)
	assert.InDelta(t, 0.0, dist.PDF(11.0), 0.0001)

	// Test that the CDF is correct
	assert.InDelta(t, 0.0, dist.CDF(-1.0), 0.0001)
	assert.InDelta(t, 0.0, dist.CDF(0.0), 0.0001)
	assert.InDelta(t, 0.5, dist.CDF(5.0), 0.0001)
	assert.InDelta(t, 1.0, dist.CDF(10.0), 0.0001)
	assert.InDelta(t, 1.0, dist.CDF(11.0), 0.0001)
}

// TestNormalDistribution tests the creation and properties of a normal distribution
func TestNormalDistribution(t *testing.T) {
	// Create a parameter with min=-inf, max=inf (effectively)
	param, err := NewParameter("test_param", -1000.0, 1000.0)
	assert.NoError(t, err)

	// Create a normal distribution with mean=0, stddev=1
	dist, err := NewNormalDistribution(param, 0.0, 1.0)
	assert.NoError(t, err)

	// Test that the distribution type is correct
	assert.Equal(t, DistributionTypeNormal, dist.Type())

	// Test the PDF values at different points
	// For a standard normal distribution, PDF(0) = 1/sqrt(2π) ≈ 0.3989
	assert.InDelta(t, 0.3989, dist.PDF(0.0), 0.0001)

	// Test that the CDF is correct
	// For a standard normal distribution, CDF(0) = 0.5
	assert.InDelta(t, 0.5, dist.CDF(0.0), 0.0001)

	// Test that the CDF approaches 0 for very negative values
	assert.InDelta(t, 0.0, dist.CDF(-5.0), 0.0001)

	// Test that the CDF approaches 1 for very positive values
	assert.InDelta(t, 1.0, dist.CDF(5.0), 0.0001)
}

// TestInvalidDistributionCreation tests that creating distributions with invalid parameters fails
func TestInvalidDistributionCreation(t *testing.T) {
	// Create a parameter with min=0, max=10
	param, err := NewParameter("test_param", 0.0, 10.0)
	assert.NoError(t, err)

	// Try to create a normal distribution with negative standard deviation
	_, err = NewNormalDistribution(param, 5.0, -1.0)
	assert.Error(t, err)
}
