package paramspace

import (
	"fmt"
	"math"
)

// DistributionType represents the type of probability distribution
type DistributionType string

const (
	DistributionTypeUniform DistributionType = "uniform"
	DistributionTypeNormal  DistributionType = "normal"
)

// Distribution represents a probability distribution for a parameter
type Distribution interface {
	// Type returns the type of distribution
	Type() DistributionType

	// PDF returns the probability density function value at x
	PDF(x float64) float64

	// CDF returns the cumulative distribution function value at x
	CDF(x float64) float64

	// Parameter returns the associated parameter
	Parameter() *Parameter
}

// UniformDistribution represents a uniform distribution over a parameter's range
type UniformDistribution struct {
	param *Parameter
	pdf   float64 // Constant PDF value within the range
}

// NewUniformDistribution creates a new uniform distribution for the given parameter
func NewUniformDistribution(param *Parameter) (Distribution, error) {
	if param == nil {
		return nil, fmt.Errorf("parameter cannot be nil")
	}

	// Calculate the constant PDF value (1 / range)
	pdf := 1.0 / (param.Max() - param.Min())

	return &UniformDistribution{
		param: param,
		pdf:   pdf,
	}, nil
}

// Type returns the type of distribution
func (d *UniformDistribution) Type() DistributionType {
	return DistributionTypeUniform
}

// PDF returns the probability density function value at x
func (d *UniformDistribution) PDF(x float64) float64 {
	if !d.param.IsValid(x) {
		return 0.0
	}
	return d.pdf
}

// CDF returns the cumulative distribution function value at x
func (d *UniformDistribution) CDF(x float64) float64 {
	if x < d.param.Min() {
		return 0.0
	}
	if x > d.param.Max() {
		return 1.0
	}
	return (x - d.param.Min()) / (d.param.Max() - d.param.Min())
}

// Parameter returns the associated parameter
func (d *UniformDistribution) Parameter() *Parameter {
	return d.param
}

// NormalDistribution represents a normal (Gaussian) distribution
type NormalDistribution struct {
	param  *Parameter
	mean   float64
	stdDev float64
}

// NewNormalDistribution creates a new normal distribution with the given mean and standard deviation
func NewNormalDistribution(param *Parameter, mean, stdDev float64) (Distribution, error) {
	if param == nil {
		return nil, fmt.Errorf("parameter cannot be nil")
	}

	if stdDev <= 0 {
		return nil, fmt.Errorf("standard deviation must be positive")
	}

	return &NormalDistribution{
		param:  param,
		mean:   mean,
		stdDev: stdDev,
	}, nil
}

// Type returns the type of distribution
func (d *NormalDistribution) Type() DistributionType {
	return DistributionTypeNormal
}

// PDF returns the probability density function value at x
func (d *NormalDistribution) PDF(x float64) float64 {
	// Standard normal PDF formula: (1 / (σ * sqrt(2π))) * e^(-(x-μ)²/(2σ²))
	exponent := -math.Pow(x-d.mean, 2) / (2 * math.Pow(d.stdDev, 2))
	coefficient := 1 / (d.stdDev * math.Sqrt(2*math.Pi))
	return coefficient * math.Exp(exponent)
}

// CDF returns the cumulative distribution function value at x
func (d *NormalDistribution) CDF(x float64) float64 {
	// Approximation of the normal CDF
	// Using the error function: CDF(x) = 0.5 * (1 + erf((x-μ)/(σ*sqrt(2))))
	return 0.5 * (1 + math.Erf((x-d.mean)/(d.stdDev*math.Sqrt(2))))
}

// Parameter returns the associated parameter
func (d *NormalDistribution) Parameter() *Parameter {
	return d.param
}
