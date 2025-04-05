package paramspace

import (
	"fmt"
)

// ParameterSpace represents the entire parameter space Θ with multiple parameters
// and their uncertainty relations
type ParameterSpace struct {
	parameters map[string]*Parameter
	relations  []*UncertaintyRelation
	// Map from parameter name to its distribution
	distributions map[string]Distribution
}

// NewParameterSpace creates a new empty parameter space
func NewParameterSpace() *ParameterSpace {
	return &ParameterSpace{
		parameters:    make(map[string]*Parameter),
		relations:     make([]*UncertaintyRelation, 0),
		distributions: make(map[string]Distribution),
	}
}

// ParameterCount returns the number of parameters in the space
func (s *ParameterSpace) ParameterCount() int {
	return len(s.parameters)
}

// RelationCount returns the number of uncertainty relations in the space
func (s *ParameterSpace) RelationCount() int {
	return len(s.relations)
}

// AddParameter adds a parameter to the space
func (s *ParameterSpace) AddParameter(param *Parameter) error {
	if param == nil {
		return fmt.Errorf("parameter cannot be nil")
	}

	// Check if a parameter with the same name already exists
	if _, exists := s.parameters[param.Name()]; exists {
		return fmt.Errorf("parameter with name '%s' already exists", param.Name())
	}

	s.parameters[param.Name()] = param
	return nil
}

// GetParameterByName retrieves a parameter by its name
func (s *ParameterSpace) GetParameterByName(name string) (*Parameter, bool) {
	param, found := s.parameters[name]
	return param, found
}

// AddRelation adds an uncertainty relation to the space
func (s *ParameterSpace) AddRelation(relation *UncertaintyRelation) error {
	if relation == nil {
		return fmt.Errorf("relation cannot be nil")
	}

	// Check if both parameters in the relation exist in the space
	param1 := relation.Parameter1()
	param2 := relation.Parameter2()

	if _, exists := s.parameters[param1.Name()]; !exists {
		return fmt.Errorf("parameter '%s' does not exist in the space", param1.Name())
	}

	if _, exists := s.parameters[param2.Name()]; !exists {
		return fmt.Errorf("parameter '%s' does not exist in the space", param2.Name())
	}

	// Check if a relation between these parameters already exists
	for _, r := range s.relations {
		if (r.Parameter1() == param1 && r.Parameter2() == param2) ||
			(r.Parameter1() == param2 && r.Parameter2() == param1) {
			return fmt.Errorf("a relation between parameters '%s' and '%s' already exists",
				param1.Name(), param2.Name())
		}
	}

	s.relations = append(s.relations, relation)
	return nil
}

// GetRelationsForParameter retrieves all relations involving a parameter
func (s *ParameterSpace) GetRelationsForParameter(param *Parameter) []*UncertaintyRelation {
	var result []*UncertaintyRelation

	for _, relation := range s.relations {
		if relation.Parameter1() == param || relation.Parameter2() == param {
			result = append(result, relation)
		}
	}

	return result
}

// GetRelationBetweenParameters retrieves the relation between two specific parameters
func (s *ParameterSpace) GetRelationBetweenParameters(param1, param2 *Parameter) (*UncertaintyRelation, bool) {
	for _, relation := range s.relations {
		if relation.Parameter1() == param1 && relation.Parameter2() == param2 {
			return relation, true
		}
	}

	return nil, false
}

// SetDistribution sets the distribution for a parameter
func (s *ParameterSpace) SetDistribution(param *Parameter, dist Distribution) error {
	if param == nil {
		return fmt.Errorf("parameter cannot be nil")
	}

	if dist == nil {
		return fmt.Errorf("distribution cannot be nil")
	}

	// Check if the parameter exists in the space
	if _, exists := s.parameters[param.Name()]; !exists {
		return fmt.Errorf("parameter '%s' does not exist in the space", param.Name())
	}

	// Check if the distribution is for the correct parameter
	if dist.Parameter() != param {
		return fmt.Errorf("distribution is for parameter '%s', not '%s'",
			dist.Parameter().Name(), param.Name())
	}

	s.distributions[param.Name()] = dist
	return nil
}

// GetDistribution retrieves the distribution for a parameter
func (s *ParameterSpace) GetDistribution(param *Parameter) (Distribution, bool) {
	if param == nil {
		return nil, false
	}

	dist, found := s.distributions[param.Name()]
	return dist, found
}

// ValidateUncertaintyRelations checks if all uncertainty relations are satisfied
// by the current distributions
func (s *ParameterSpace) ValidateUncertaintyRelations() (bool, []*UncertaintyRelation) {
	var violations []*UncertaintyRelation

	for _, relation := range s.relations {
		param1 := relation.Parameter1()
		param2 := relation.Parameter2()

		dist1, found1 := s.GetDistribution(param1)
		dist2, found2 := s.GetDistribution(param2)

		// If either distribution is missing, we can't validate
		if !found1 || !found2 {
			continue
		}

		// Check if the relation is satisfied
		if !relation.IsSatisfied(dist1, dist2) {
			violations = append(violations, relation)
		}
	}

	return len(violations) == 0, violations
}
