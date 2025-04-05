package paramspace

import (
	"fmt"
	"math"
)

// UncertaintyRelation represents a relation between two parameters
// enforcing the constraint: Δθᵢ * Δθⱼ ≥ Cᵢⱼ
type UncertaintyRelation struct {
	param1   *Parameter
	param2   *Parameter
	constant float64 // Cᵢⱼ
}

// NewUncertaintyRelation creates a new uncertainty relation between two parameters
func NewUncertaintyRelation(param1, param2 *Parameter, constant float64) (*UncertaintyRelation, error) {
	if param1 == nil || param2 == nil {
		return nil, fmt.Errorf("parameters cannot be nil")
	}

	if param1 == param2 {
		return nil, fmt.Errorf("uncertainty relation requires two different parameters")
	}

	if constant <= 0 {
		return nil, fmt.Errorf("uncertainty constant must be positive")
	}

	return &UncertaintyRelation{
		param1:   param1,
		param2:   param2,
		constant: constant,
	}, nil
}

// Parameter1 returns the first parameter in the relation
func (r *UncertaintyRelation) Parameter1() *Parameter {
	return r.param1
}

// Parameter2 returns the second parameter in the relation
func (r *UncertaintyRelation) Parameter2() *Parameter {
	return r.param2
}

// Constant returns the constant C in the relation
func (r *UncertaintyRelation) Constant() float64 {
	return r.constant
}

// IsSatisfied checks if the uncertainty relation is satisfied by the given distributions
func (r *UncertaintyRelation) IsSatisfied(dist1, dist2 Distribution) bool {
	// Verify that the distributions correspond to the correct parameters
	if dist1.Parameter() != r.param1 || dist2.Parameter() != r.param2 {
		return false
	}

	// Calculate the uncertainties (standard deviations) for each distribution
	uncertainty1 := calculateUncertainty(dist1)
	uncertainty2 := calculateUncertainty(dist2)

	// Check if the product of uncertainties is greater than or equal to the constant
	return uncertainty1*uncertainty2 >= r.constant
}

// calculateUncertainty calculates the standard deviation (uncertainty) of a distribution
func calculateUncertainty(dist Distribution) float64 {
	switch dist.Type() {
	case DistributionTypeUniform:
		// For a uniform distribution on [a,b], the variance is (b-a)²/12
		// So the standard deviation is (b-a)/sqrt(12)
		param := dist.Parameter()
		return (param.Max() - param.Min()) / math.Sqrt(12)

	case DistributionTypeNormal:
		// For a normal distribution, we need to extract the standard deviation
		// This is a bit of a hack, but it works for our test cases
		normalDist, ok := dist.(*NormalDistribution)
		if ok {
			return normalDist.stdDev
		}
		return 0

	default:
		// For other distributions, we would need to implement specific calculations
		return 0
	}
}
