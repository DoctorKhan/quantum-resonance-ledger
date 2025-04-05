package paramspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParameterCreation tests the creation of a parameter with constraints
func TestParameterCreation(t *testing.T) {
	// Create a parameter with min=0, max=10
	param, err := NewParameter("test_param", 0.0, 10.0)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the parameter has the correct name
	assert.Equal(t, "test_param", param.Name())

	// Assert that the parameter has the correct min and max values
	assert.Equal(t, 0.0, param.Min())
	assert.Equal(t, 10.0, param.Max())
}

// TestParameterValidation tests the validation of parameter values
func TestParameterValidation(t *testing.T) {
	// Create a parameter with min=0, max=10
	param, err := NewParameter("test_param", 0.0, 10.0)
	assert.NoError(t, err)

	// Test valid values
	assert.True(t, param.IsValid(0.0))
	assert.True(t, param.IsValid(5.0))
	assert.True(t, param.IsValid(10.0))

	// Test invalid values
	assert.False(t, param.IsValid(-1.0))
	assert.False(t, param.IsValid(11.0))
}

// TestInvalidParameterCreation tests that creating a parameter with invalid constraints fails
func TestInvalidParameterCreation(t *testing.T) {
	// Try to create a parameter with min > max
	_, err := NewParameter("invalid_param", 10.0, 0.0)

	// Assert that an error occurred
	assert.Error(t, err)
}
