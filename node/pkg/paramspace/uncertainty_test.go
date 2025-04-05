package paramspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUncertaintyRelation tests the creation and validation of uncertainty relations
func TestUncertaintyRelation(t *testing.T) {
	// Create two parameters
	param1, err := NewParameter("param1", 0.0, 10.0)
	assert.NoError(t, err)

	param2, err := NewParameter("param2", 0.0, 10.0)
	assert.NoError(t, err)

	// Create an uncertainty relation with constant C = 2.0
	relation, err := NewUncertaintyRelation(param1, param2, 2.0)
	assert.NoError(t, err)

	// Test that the relation has the correct parameters and constant
	assert.Equal(t, param1, relation.Parameter1())
	assert.Equal(t, param2, relation.Parameter2())
	assert.Equal(t, 2.0, relation.Constant())
}

// TestUncertaintyRelationValidation tests the validation of parameter uncertainties
func TestUncertaintyRelationValidation(t *testing.T) {
	// Create two parameters
	param1, err := NewParameter("param1", 0.0, 10.0)
	assert.NoError(t, err)

	param2, err := NewParameter("param2", 0.0, 10.0)
	assert.NoError(t, err)

	// Create distributions with different uncertainties
	dist1, err := NewUniformDistribution(param1)
	assert.NoError(t, err)

	// For uniform distribution on [0,10], the variance is (10-0)²/12 = 8.33
	// So the standard deviation (uncertainty) is sqrt(8.33) ≈ 2.89

	dist2, err := NewNormalDistribution(param2, 5.0, 0.5)
	assert.NoError(t, err)
	// For normal distribution, the standard deviation is directly 0.5

	// Create an uncertainty relation with constant C = 1.0
	relation, err := NewUncertaintyRelation(param1, param2, 1.0)
	assert.NoError(t, err)

	// Test that the relation is satisfied (2.89 * 0.5 = 1.445 > 1.0)
	assert.True(t, relation.IsSatisfied(dist1, dist2))

	// Create a distribution with smaller uncertainty
	dist3, err := NewNormalDistribution(param2, 5.0, 0.3)
	assert.NoError(t, err)

	// Test that the relation is not satisfied (2.89 * 0.3 = 0.867 < 1.0)
	assert.False(t, relation.IsSatisfied(dist1, dist3))
}

// TestInvalidUncertaintyRelationCreation tests that creating uncertainty relations with invalid parameters fails
func TestInvalidUncertaintyRelationCreation(t *testing.T) {
	// Create a parameter
	param, err := NewParameter("param", 0.0, 10.0)
	assert.NoError(t, err)

	// Try to create an uncertainty relation with the same parameter twice
	_, err = NewUncertaintyRelation(param, param, 1.0)
	assert.Error(t, err)

	// Try to create an uncertainty relation with a negative constant
	_, err = NewUncertaintyRelation(param, param, -1.0)
	assert.Error(t, err)
}
