package simulation

// Placeholder for probability distribution interface/structs
type Distribution interface {
	Type() string
	Sample() float64 // Draw a random sample from the distribution
	// Methods like PDF(float64), Update(...) will be added later
}

// Parameter represents a simulation parameter governed by a probability distribution.
// Placeholder implementation.
type Parameter struct {
	Name         string
	Distribution Distribution // Interface for the distribution
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
