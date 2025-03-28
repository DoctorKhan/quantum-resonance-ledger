package simulation

import (
	"fmt"
	"math"
	"math/rand"
)

// Parameter represents a parameter in the simulation
type Parameter struct {
	Name         string
	Min          float64
	Max          float64
	CurrentValue float64
}

// NewParameter creates a new parameter with the given name and bounds
func NewParameter(name string, min, max float64) *Parameter {
	// Default current value is the midpoint
	currentValue := (min + max) / 2.0

	return &Parameter{
		Name:         name,
		Min:          min,
		Max:          max,
		CurrentValue: currentValue,
	}
}

// Distribution is an interface for probability distributions
type Distribution interface {
	// Sample returns a random sample from the distribution
	Sample() float64

	// PDF returns the probability density function value at x
	PDF(x float64) float64

	// GetParameter returns the parameter associated with this distribution
	GetParameter() *Parameter

	// GetUncertainty returns the standard deviation of the distribution
	GetUncertainty() float64
}

// UniformDistribution represents a uniform distribution
type UniformDistribution struct {
	Parameter *Parameter
}

// NewUniformDistribution creates a new uniform distribution for the given parameter
func NewUniformDistribution(param *Parameter) *UniformDistribution {
	return &UniformDistribution{
		Parameter: param,
	}
}

// Sample returns a random sample from the uniform distribution
func (d *UniformDistribution) Sample() float64 {
	// Generate a random value between Min and Max
	return d.Parameter.Min + rand.Float64()*(d.Parameter.Max-d.Parameter.Min)
}

// PDF returns the probability density function value at x
func (d *UniformDistribution) PDF(x float64) float64 {
	// If x is outside the range, PDF is 0
	if x < d.Parameter.Min || x > d.Parameter.Max {
		return 0.0
	}

	// For uniform distribution, PDF is 1/(Max-Min)
	return 1.0 / (d.Parameter.Max - d.Parameter.Min)
}

// GetParameter returns the parameter associated with this distribution
func (d *UniformDistribution) GetParameter() *Parameter {
	return d.Parameter
}

// GetUncertainty returns the standard deviation of the uniform distribution
func (d *UniformDistribution) GetUncertainty() float64 {
	// For uniform distribution on [a,b], the standard deviation is (b-a)/sqrt(12)
	return (d.Parameter.Max - d.Parameter.Min) / math.Sqrt(12.0)
}

// NormalDistribution represents a normal (Gaussian) distribution
type NormalDistribution struct {
	Parameter *Parameter
	Mean      float64
	StdDev    float64
}

// NewNormalDistribution creates a new normal distribution for the given parameter
func NewNormalDistribution(param *Parameter, mean, stdDev float64) *NormalDistribution {
	return &NormalDistribution{
		Parameter: param,
		Mean:      mean,
		StdDev:    stdDev,
	}
}

// Sample returns a random sample from the normal distribution
func (d *NormalDistribution) Sample() float64 {
	// Generate a random value from a normal distribution
	sample := rand.NormFloat64()*d.StdDev + d.Mean

	// Truncate to the parameter bounds
	if sample < d.Parameter.Min {
		sample = d.Parameter.Min
	} else if sample > d.Parameter.Max {
		sample = d.Parameter.Max
	}

	return sample
}

// PDF returns the probability density function value at x
func (d *NormalDistribution) PDF(x float64) float64 {
	// If x is outside the range, PDF is 0
	if x < d.Parameter.Min || x > d.Parameter.Max {
		return 0.0
	}

	// For normal distribution, PDF is (1/(σ√(2π))) * e^(-((x-μ)²/(2σ²)))
	exponent := -math.Pow(x-d.Mean, 2) / (2 * math.Pow(d.StdDev, 2))
	coefficient := 1.0 / (d.StdDev * math.Sqrt(2*math.Pi))

	return coefficient * math.Exp(exponent)
}

// GetParameter returns the parameter associated with this distribution
func (d *NormalDistribution) GetParameter() *Parameter {
	return d.Parameter
}

// GetUncertainty returns the standard deviation of the normal distribution
func (d *NormalDistribution) GetUncertainty() float64 {
	// For normal distribution, the standard deviation is directly the StdDev field
	return d.StdDev
}

// UncertaintyRelation represents an uncertainty relation between two parameters
type UncertaintyRelation struct {
	Parameter1 *Parameter
	Parameter2 *Parameter
	Constant   float64 // The constant C in the relation Δθ_i * Δθ_j ≥ C
}

// NewUncertaintyRelation creates a new uncertainty relation
func NewUncertaintyRelation(param1, param2 *Parameter, constant float64) *UncertaintyRelation {
	return &UncertaintyRelation{
		Parameter1: param1,
		Parameter2: param2,
		Constant:   constant,
	}
}

// IsSatisfied checks if the uncertainty relation is satisfied by the given distributions
func (r *UncertaintyRelation) IsSatisfied(dist1, dist2 Distribution) bool {
	// Check if the distributions correspond to the correct parameters
	if dist1.GetParameter() != r.Parameter1 || dist2.GetParameter() != r.Parameter2 {
		return false
	}

	// Get the uncertainties (standard deviations)
	uncertainty1 := dist1.GetUncertainty()
	uncertainty2 := dist2.GetUncertainty()

	// Check if the product of uncertainties is greater than or equal to the constant
	return uncertainty1*uncertainty2 >= r.Constant
}

// ParameterManager manages parameters, distributions, and uncertainty relations
type ParameterManager struct {
	Parameters           []*Parameter
	Distributions        map[string]Distribution
	UncertaintyRelations []*UncertaintyRelation
}

// NewParameterManager creates a new parameter manager
func NewParameterManager() *ParameterManager {
	return &ParameterManager{
		Parameters:           make([]*Parameter, 0),
		Distributions:        make(map[string]Distribution),
		UncertaintyRelations: make([]*UncertaintyRelation, 0),
	}
}

// AddParameter adds a parameter to the manager
func (m *ParameterManager) AddParameter(param *Parameter) {
	m.Parameters = append(m.Parameters, param)
}

// SetDistribution sets the distribution for a parameter
func (m *ParameterManager) SetDistribution(param *Parameter, dist Distribution) {
	m.Distributions[param.Name] = dist
}

// GetDistribution gets the distribution for a parameter
func (m *ParameterManager) GetDistribution(param *Parameter) (Distribution, bool) {
	dist, ok := m.Distributions[param.Name]
	return dist, ok
}

// AddUncertaintyRelation adds an uncertainty relation to the manager
func (m *ParameterManager) AddUncertaintyRelation(relation *UncertaintyRelation) {
	m.UncertaintyRelations = append(m.UncertaintyRelations, relation)
}

// ValidateUncertaintyRelations checks if all uncertainty relations are satisfied
func (m *ParameterManager) ValidateUncertaintyRelations() (bool, []*UncertaintyRelation) {
	violations := make([]*UncertaintyRelation, 0)

	for _, relation := range m.UncertaintyRelations {
		dist1, ok1 := m.GetDistribution(relation.Parameter1)
		dist2, ok2 := m.GetDistribution(relation.Parameter2)

		if !ok1 || !ok2 {
			// If either distribution is missing, we can't validate
			continue
		}

		if !relation.IsSatisfied(dist1, dist2) {
			violations = append(violations, relation)
		}
	}

	return len(violations) == 0, violations
}

// SampleParameters samples values for all parameters from their distributions
func (m *ParameterManager) SampleParameters() {
	for _, param := range m.Parameters {
		dist, ok := m.GetDistribution(param)
		if ok {
			param.CurrentValue = dist.Sample()
		}
	}
}

// AdjustDistributions adjusts distributions to satisfy uncertainty relations
func (m *ParameterManager) AdjustDistributions() error {
	valid, violations := m.ValidateUncertaintyRelations()
	if valid {
		return nil
	}

	// For each violation, try to adjust the distributions
	for _, relation := range violations {
		dist1, ok1 := m.GetDistribution(relation.Parameter1)
		dist2, ok2 := m.GetDistribution(relation.Parameter2)

		if !ok1 || !ok2 {
			continue
		}

		// Get the current uncertainties
		uncertainty1 := dist1.GetUncertainty()
		uncertainty2 := dist2.GetUncertainty()

		// Calculate the required product
		requiredProduct := relation.Constant

		// Calculate the current product
		currentProduct := uncertainty1 * uncertainty2

		if currentProduct < requiredProduct {
			// We need to increase the product
			// For simplicity, we'll increase both uncertainties by the same factor
			factor := math.Sqrt(requiredProduct / currentProduct)

			// Create new distributions with adjusted uncertainties
			// This is a simplified approach and might not work for all distribution types
			switch d1 := dist1.(type) {
			case *NormalDistribution:
				m.SetDistribution(relation.Parameter1, NewNormalDistribution(relation.Parameter1, d1.Mean, d1.StdDev*factor))
			}

			switch d2 := dist2.(type) {
			case *NormalDistribution:
				m.SetDistribution(relation.Parameter2, NewNormalDistribution(relation.Parameter2, d2.Mean, d2.StdDev*factor))
			}
		}
	}

	// Check if we've resolved all violations
	valid, violations = m.ValidateUncertaintyRelations()
	if !valid {
		return fmt.Errorf("failed to adjust distributions to satisfy uncertainty relations")
	}

	return nil
}
