package paramspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParameterSpaceCreation tests the creation of a parameter space
func TestParameterSpaceCreation(t *testing.T) {
	// Create a new parameter space
	space := NewParameterSpace()

	// Test that the space is initially empty
	assert.Equal(t, 0, space.ParameterCount())
	assert.Equal(t, 0, space.RelationCount())
}

// TestParameterSpaceAddParameter tests adding parameters to the space
func TestParameterSpaceAddParameter(t *testing.T) {
	// Create a new parameter space
	space := NewParameterSpace()

	// Create parameters
	param1, err := NewParameter("param1", 0.0, 10.0)
	assert.NoError(t, err)

	param2, err := NewParameter("param2", -5.0, 5.0)
	assert.NoError(t, err)

	// Add parameters to the space
	err = space.AddParameter(param1)
	assert.NoError(t, err)

	err = space.AddParameter(param2)
	assert.NoError(t, err)

	// Test that the parameters were added
	assert.Equal(t, 2, space.ParameterCount())

	// Test retrieving parameters by name
	retrievedParam1, found := space.GetParameterByName("param1")
	assert.True(t, found)
	assert.Equal(t, param1, retrievedParam1)

	retrievedParam2, found := space.GetParameterByName("param2")
	assert.True(t, found)
	assert.Equal(t, param2, retrievedParam2)

	// Test retrieving a non-existent parameter
	_, found = space.GetParameterByName("non_existent")
	assert.False(t, found)
}

// TestParameterSpaceAddRelation tests adding uncertainty relations to the space
func TestParameterSpaceAddRelation(t *testing.T) {
	// Create a new parameter space
	space := NewParameterSpace()

	// Create parameters
	param1, err := NewParameter("param1", 0.0, 10.0)
	assert.NoError(t, err)

	param2, err := NewParameter("param2", -5.0, 5.0)
	assert.NoError(t, err)

	// Add parameters to the space
	err = space.AddParameter(param1)
	assert.NoError(t, err)

	err = space.AddParameter(param2)
	assert.NoError(t, err)

	// Create an uncertainty relation
	relation, err := NewUncertaintyRelation(param1, param2, 2.0)
	assert.NoError(t, err)

	// Add the relation to the space
	err = space.AddRelation(relation)
	assert.NoError(t, err)

	// Test that the relation was added
	assert.Equal(t, 1, space.RelationCount())

	// Test retrieving relations for a parameter
	relations := space.GetRelationsForParameter(param1)
	assert.Equal(t, 1, len(relations))
	assert.Equal(t, relation, relations[0])

	// Test retrieving relations between two parameters
	relation2, found := space.GetRelationBetweenParameters(param1, param2)
	assert.True(t, found)
	assert.Equal(t, relation, relation2)

	// Test retrieving a non-existent relation
	_, found = space.GetRelationBetweenParameters(param2, param1)
	assert.False(t, found)
}

// TestParameterSpaceValidation tests validating distributions against uncertainty relations
func TestParameterSpaceValidation(t *testing.T) {
	// Create a new parameter space
	space := NewParameterSpace()

	// Create parameters
	param1, err := NewParameter("param1", 0.0, 10.0)
	assert.NoError(t, err)

	param2, err := NewParameter("param2", -5.0, 5.0)
	assert.NoError(t, err)

	// Add parameters to the space
	err = space.AddParameter(param1)
	assert.NoError(t, err)

	err = space.AddParameter(param2)
	assert.NoError(t, err)

	// Create an uncertainty relation
	relation, err := NewUncertaintyRelation(param1, param2, 1.0)
	assert.NoError(t, err)

	// Add the relation to the space
	err = space.AddRelation(relation)
	assert.NoError(t, err)

	// Create distributions
	dist1, err := NewUniformDistribution(param1)
	assert.NoError(t, err)

	dist2, err := NewNormalDistribution(param2, 0.0, 0.5)
	assert.NoError(t, err)

	// Add distributions to the space
	err = space.SetDistribution(param1, dist1)
	assert.NoError(t, err)

	err = space.SetDistribution(param2, dist2)
	assert.NoError(t, err)

	// Test that the distributions satisfy the uncertainty relation
	valid, violations := space.ValidateUncertaintyRelations()
	assert.True(t, valid)
	assert.Equal(t, 0, len(violations))

	// Create a distribution with smaller uncertainty
	dist3, err := NewNormalDistribution(param2, 0.0, 0.2)
	assert.NoError(t, err)

	// Replace the distribution
	err = space.SetDistribution(param2, dist3)
	assert.NoError(t, err)

	// Test that the distributions no longer satisfy the uncertainty relation
	valid, violations = space.ValidateUncertaintyRelations()
	assert.False(t, valid)
	assert.Equal(t, 1, len(violations))
	assert.Equal(t, relation, violations[0])
}
