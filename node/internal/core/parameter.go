package core

import (
	"fmt"
	"math"
)

// Placeholder for probability distribution interface/structs
type Distribution interface {
	Type() string
	Sample() float64 // Draw a random sample from the distribution
	Mean() float64   // Return the configured mean
	StdDev() float64 // Return the configured standard deviation
	// Methods like PDF(float64), Update(...) will be added later
}

// Parameter represents a simulation parameter governed by a probability distribution.
// Placeholder implementation.
type Parameter struct {
	Name         string
	Distribution Distribution // Interface for the distribution
	CurrentValue float64      // The current sampled or updated value of the parameter
}

// NewParameter creates a new parameter.
// Placeholder implementation.
func NewParameter(name string, dist Distribution) *Parameter {
	// Correct implementation:
	return &Parameter{
		Name:         name,
		Distribution: dist,
	}
}

// Note: CurrentValue is typically initialized by calling Sample() or an update rule.

// Update applies the parameter update rule (simplified Langevin dynamics).
// It modifies the Parameter's CurrentValue based on the gradient, learning rate, and time step.
// It also ensures the value stays within the distribution's bounds if applicable.
// NOTE: This is a simplified version, ignoring noise and Laplacian terms for now.
func (p *Parameter) Update(hamiltonianGradient, eta, dt float64) error {
	if p.Distribution == nil {
		return fmt.Errorf("parameter '%s' has no distribution, cannot update", p.Name)
	}

	// Calculate the change based on the simplified Langevin dynamics (deterministic part)
	// newValue = oldValue - learningRate * timeStep * gradient
	newValue := p.CurrentValue - eta*dt*hamiltonianGradient

	// Check for bounds if the distribution is TruncatedGaussian
	// We need type assertion to access Min/Max fields
	if tg, ok := p.Distribution.(*TruncatedGaussian); ok {
		if newValue < tg.Min {
			newValue = tg.Min
		} else if newValue > tg.Max {
			newValue = tg.Max
		}
	}
	// TODO: Add similar bounds checks if other distribution types have them.

	// Handle potential NaN/Inf values resulting from calculation
	if math.IsNaN(newValue) || math.IsInf(newValue, 0) {
		return fmt.Errorf("parameter update resulted in invalid value (NaN or Inf)")
	}

	p.CurrentValue = newValue

	return nil
}

// UncertaintyRelation defines a constraint between the standard deviations of two parameters.
// Represents Δθᵢ ⋅ Δθⱼ ≥ Cᵢⱼ
type UncertaintyRelation struct {
	Param1   *Parameter // Pointer to the first parameter
	Param2   *Parameter // Pointer to the second parameter
	Constant float64    // The constant Cᵢⱼ
}

// NewUncertaintyRelation creates a new uncertainty relation constraint.
func NewUncertaintyRelation(p1, p2 *Parameter, constant float64) (*UncertaintyRelation, error) {
	if p1 == nil || p2 == nil {
		return nil, fmt.Errorf("parameters cannot be nil")
	}
	if p1 == p2 {
		return nil, fmt.Errorf("uncertainty relation must be between two different parameters")
	}
	if constant < 0 {
		return nil, fmt.Errorf("uncertainty relation constant cannot be negative (got %f)", constant)
	}
	return &UncertaintyRelation{
		Param1:   p1,
		Param2:   p2,
		Constant: constant,
	}, nil
}

// Validate checks if the uncertainty relation is currently satisfied.
// It requires the Distribution interface to have a StdDev() method.
func (ur *UncertaintyRelation) Validate() (bool, error) {
	if ur.Param1.Distribution == nil || ur.Param2.Distribution == nil {
		return false, fmt.Errorf("parameters must have distributions assigned to validate uncertainty relation")
	}

	stdDev1 := ur.Param1.Distribution.StdDev()
	stdDev2 := ur.Param2.Distribution.StdDev()

	// Handle potential NaN or Inf from StdDev() if necessary
	if math.IsNaN(stdDev1) || math.IsNaN(stdDev2) || math.IsInf(stdDev1, 0) || math.IsInf(stdDev2, 0) {
		return false, fmt.Errorf("invalid standard deviation encountered (NaN or Inf)")
	}
	// Ensure standard deviations are non-negative
	if stdDev1 < 0 || stdDev2 < 0 {
		// This shouldn't happen if StdDev() is implemented correctly, but check defensively
		return false, fmt.Errorf("standard deviation cannot be negative (got %.4f, %.4f)", stdDev1, stdDev2)
	}

	// The core check: Δθ₁ ⋅ Δθ₂ ≥ C₁₂
	return (stdDev1 * stdDev2) >= ur.Constant, nil
}

// ParameterManager manages a collection of parameters and their uncertainty relations.
type ParameterManager struct {
	Parameters           map[string]*Parameter  // Map parameter name to Parameter struct
	UncertaintyRelations []*UncertaintyRelation // List of relations
}

// NewParameterManager creates a new ParameterManager.
func NewParameterManager() *ParameterManager {
	return &ParameterManager{
		Parameters:           make(map[string]*Parameter),
		UncertaintyRelations: make([]*UncertaintyRelation, 0),
	}
}

// AddParameter adds a parameter to the manager.
func (pm *ParameterManager) AddParameter(param *Parameter) error {
	if param == nil {
		return fmt.Errorf("cannot add a nil parameter")
	}
	if _, exists := pm.Parameters[param.Name]; exists {
		return fmt.Errorf("parameter with name '%s' already exists", param.Name)
	}
	pm.Parameters[param.Name] = param
	return nil
}

// GetParameter retrieves a parameter by name.
func (pm *ParameterManager) GetParameter(name string) (*Parameter, error) {
	param, exists := pm.Parameters[name]
	if !exists {
		return nil, fmt.Errorf("parameter '%s' not found", name)
	}
	return param, nil
}

// AddUncertaintyRelation adds an uncertainty relation constraint.
// It ensures the involved parameters are already managed.
func (pm *ParameterManager) AddUncertaintyRelation(ur *UncertaintyRelation) error {
	if ur == nil {
		return fmt.Errorf("cannot add a nil uncertainty relation")
	}
	// Check if parameters exist in the manager
	_, p1Exists := pm.Parameters[ur.Param1.Name]
	_, p2Exists := pm.Parameters[ur.Param2.Name]
	if !p1Exists || !p2Exists {
		return fmt.Errorf("parameters '%s' and/or '%s' not found in manager", ur.Param1.Name, ur.Param2.Name)
	}
	// Ensure the relation uses the parameters managed by this manager instance
	if pm.Parameters[ur.Param1.Name] != ur.Param1 || pm.Parameters[ur.Param2.Name] != ur.Param2 {
		return fmt.Errorf("uncertainty relation parameters do not match managed parameters")
	}

	pm.UncertaintyRelations = append(pm.UncertaintyRelations, ur)
	return nil
}

// ValidateAllUncertaintyRelations checks all registered uncertainty relations.
// Returns true if all are satisfied, false otherwise, along with a list of violated relations.
func (pm *ParameterManager) ValidateAllUncertaintyRelations() (bool, []*UncertaintyRelation, error) {
	violations := make([]*UncertaintyRelation, 0)
	allValid := true
	var firstError error

	for _, ur := range pm.UncertaintyRelations {
		valid, err := ur.Validate()
		if err != nil {
			// Collect the first error encountered, but continue checking others
			if firstError == nil {
				firstError = fmt.Errorf("error validating relation between '%s' and '%s': %w", ur.Param1.Name, ur.Param2.Name, err)
			}
			allValid = false                    // Consider relation invalid if validation errors out
			violations = append(violations, ur) // Add to violations if error occurs
			continue                            // Move to the next relation
		}
		if !valid {
			allValid = false
			violations = append(violations, ur)
		}
	}
	return allValid, violations, firstError
}

// TODO: Add methods for Parameter Update Rule (e.g., ApplyUpdateRule)
