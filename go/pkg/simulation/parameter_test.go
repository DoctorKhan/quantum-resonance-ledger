package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDistributionSampling tests sampling from different distributions
func TestDistributionSampling(t *testing.T) {
	// Test uniform distribution
	param := NewParameter("test_param", 0.0, 10.0)
	uniformDist := NewUniformDistribution(param)

	// Sample from the uniform distribution
	for i := 0; i < 100; i++ {
		sample := uniformDist.Sample()
		// Assert that the sample is within the parameter bounds
		assert.GreaterOrEqual(t, sample, param.Min)
		assert.LessOrEqual(t, sample, param.Max)
	}

	// Test normal distribution
	normalDist := NewNormalDistribution(param, 5.0, 2.0)

	// Sample from the normal distribution
	for i := 0; i < 100; i++ {
		sample := normalDist.Sample()
		// Assert that the sample is within the parameter bounds
		assert.GreaterOrEqual(t, sample, param.Min)
		assert.LessOrEqual(t, sample, param.Max)
	}
}

// TestParameterInitialization ensures parameters are created correctly
func TestParameterInitialization(t *testing.T) {
	// Create a parameter
	param := NewParameter("test_param", 0.0, 10.0)

	// Assert that the parameter has the correct properties
	assert.Equal(t, "test_param", param.Name)
	assert.Equal(t, 0.0, param.Min)
	assert.Equal(t, 10.0, param.Max)
	assert.Equal(t, 5.0, param.CurrentValue) // Default value should be the midpoint
}

// TestUncertaintyRelation tests the uncertainty relation between parameters
func TestUncertaintyRelation(t *testing.T) {
	// Create parameters
	param1 := NewParameter("param1", 0.0, 10.0)
	param2 := NewParameter("param2", 0.0, 10.0)

	// Create distributions
	dist1 := NewUniformDistribution(param1)
	dist2 := NewUniformDistribution(param2)

	// Create an uncertainty relation
	relation := NewUncertaintyRelation(param1, param2, 1.0)

	// Assert that the relation has the correct properties
	assert.Equal(t, param1, relation.Parameter1)
	assert.Equal(t, param2, relation.Parameter2)
	assert.Equal(t, 1.0, relation.Constant)

	// Check if the relation is satisfied
	satisfied := relation.IsSatisfied(dist1, dist2)

	// For uniform distributions on [0,10], the standard deviation is (10-0)/sqrt(12) ≈ 2.89
	// So the product of uncertainties is 2.89 * 2.89 ≈ 8.35, which is > 1.0
	assert.True(t, satisfied)

	// Create a distribution with smaller uncertainty
	dist3 := NewNormalDistribution(param2, 5.0, 0.2)

	// Check if the relation is still satisfied
	satisfied = relation.IsSatisfied(dist1, dist3)

	// For uniform distribution on [0,10], the standard deviation is ≈ 2.89
	// For normal distribution with stddev=0.2, the uncertainty is 0.2
	// So the product of uncertainties is 2.89 * 0.2 ≈ 0.58, which is < 1.0
	assert.False(t, satisfied)
}

// TestParameterManager tests the parameter manager
func TestParameterManager(t *testing.T) {
	// Create a parameter manager
	manager := NewParameterManager()

	// Create parameters
	param1 := NewParameter("param1", 0.0, 10.0)
	param2 := NewParameter("param2", 0.0, 10.0)

	// Add parameters to the manager
	manager.AddParameter(param1)
	manager.AddParameter(param2)

	// Assert that the manager has the correct number of parameters
	assert.Equal(t, 2, len(manager.Parameters))

	// Create distributions
	dist1 := NewUniformDistribution(param1)
	dist2 := NewUniformDistribution(param2)

	// Set distributions
	manager.SetDistribution(param1, dist1)
	manager.SetDistribution(param2, dist2)

	// Assert that the distributions were set correctly
	assert.Equal(t, dist1, manager.Distributions[param1.Name])
	assert.Equal(t, dist2, manager.Distributions[param2.Name])

	// Create an uncertainty relation
	relation := NewUncertaintyRelation(param1, param2, 1.0)

	// Add the relation to the manager
	manager.AddUncertaintyRelation(relation)

	// Assert that the manager has the correct number of relations
	assert.Equal(t, 1, len(manager.UncertaintyRelations))

	// Validate the uncertainty relations
	valid, violations := manager.ValidateUncertaintyRelations()

	// Assert that all relations are satisfied
	assert.True(t, valid)
	assert.Equal(t, 0, len(violations))

	// Create a distribution with smaller uncertainty
	dist3 := NewNormalDistribution(param2, 5.0, 0.2)

	// Set the new distribution
	manager.SetDistribution(param2, dist3)

	// Validate the uncertainty relations again
	valid, violations = manager.ValidateUncertaintyRelations()

	// Assert that some relations are violated
	assert.False(t, valid)
	assert.Equal(t, 1, len(violations))
	assert.Equal(t, relation, violations[0])
}
