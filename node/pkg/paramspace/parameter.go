package paramspace

import (
	"fmt"
)

// Parameter represents a parameter in the parameter space with constraints
type Parameter struct {
	name string
	min  float64
	max  float64
}

// NewParameter creates a new parameter with the given name and constraints
func NewParameter(name string, min, max float64) (*Parameter, error) {
	// Validate that min <= max
	if min > max {
		return nil, fmt.Errorf("invalid parameter constraints: min (%f) > max (%f)", min, max)
	}

	return &Parameter{
		name: name,
		min:  min,
		max:  max,
	}, nil
}

// Name returns the name of the parameter
func (p *Parameter) Name() string {
	return p.name
}

// Min returns the minimum value of the parameter
func (p *Parameter) Min() float64 {
	return p.min
}

// Max returns the maximum value of the parameter
func (p *Parameter) Max() float64 {
	return p.max
}

// IsValid checks if a value is valid for this parameter
func (p *Parameter) IsValid(value float64) bool {
	return value >= p.min && value <= p.max
}
