package core

import (
	"fmt"
	"math"
	"sync"
)

// WSIOracle defines the interface for getting external stablecoin prices.
type WSIOracle interface {
	GetPrice(stablecoinID string) (float64, error)
	// TODO: Add methods for checking oracle health/freshness?
}

// WSIManager manages the state and parameters related to the Wavefunction Stability Index.
type WSIManager struct {
	mu sync.RWMutex
	// Map stablecoin ID (e.g., "qUSDC", "qDAI") to its target weight Parameter
	weights map[string]*Parameter
	// Map stablecoin ID to its oracle interface
	oracles map[string]WSIOracle
	// The target peg value (e.g., 1.0 for USD)
	targetPeg float64
}

// NewWSIManager creates a new WSI manager.
func NewWSIManager(targetPeg float64) *WSIManager {
	return &WSIManager{
		weights:   make(map[string]*Parameter),
		oracles:   make(map[string]WSIOracle),
		targetPeg: targetPeg,
	}
}

// AddConstituent adds a stablecoin to be tracked by the WSI.
func (wm *WSIManager) AddConstituent(id string, weightParam *Parameter, oracle WSIOracle) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.weights[id]; exists {
		return fmt.Errorf("constituent '%s' already exists in WSI manager", id)
	}
	if weightParam == nil {
		return fmt.Errorf("weight parameter cannot be nil for constituent '%s'", id)
	}
	if oracle == nil {
		return fmt.Errorf("oracle cannot be nil for constituent '%s'", id)
	}

	wm.weights[id] = weightParam
	wm.oracles[id] = oracle
	// TODO: Initial validation? Ensure initial weights sum reasonably?
	return nil
}

// GetValue calculates the current value of the WSI based on constituent weights and oracle prices.
func (wm *WSIManager) GetValue() (float64, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	currentValue := 0.0
	totalWeight := 0.0 // Keep track to normalize/check later if needed

	for id, weightParam := range wm.weights {
		oracle, ok := wm.oracles[id]
		if !ok {
			// Should not happen if AddConstituent is used correctly
			return 0, fmt.Errorf("internal error: oracle missing for constituent '%s'", id)
		}

		price, err := oracle.GetPrice(id)
		if err != nil {
			// How to handle oracle errors? Exclude? Return error? Use stale price?
			// For now, return an error.
			return 0, fmt.Errorf("failed to get price for '%s' from oracle: %w", id, err)
		}
		if math.IsNaN(price) || math.IsInf(price, 0) || price < 0 {
			return 0, fmt.Errorf("invalid price received from oracle for '%s': %f", id, price)
		}

		weight := weightParam.CurrentValue // Use the parameter's current value
		currentValue += weight * price
		totalWeight += weight
	}

	// Optional: Check if total weight is close to 1.0 (depends on how weights are updated)
	// tolerance := 1e-6
	// if math.Abs(totalWeight-1.0) > tolerance {
	//  fmt.Printf("Warning: WSI weights sum to %f, not 1.0\n", totalWeight)
	// }

	return currentValue, nil
}

// CalculatePegPenalty calculates the Hamiltonian cost term for WSI peg deviation.
// Penalty = (CurrentValue - TargetPeg)^2
func (wm *WSIManager) CalculatePegPenalty() (float64, error) {
	currentValue, err := wm.GetValue()
	if err != nil {
		// If value calculation fails (e.g., oracle error), return max penalty or error?
		// Returning error for now.
		return math.Inf(1), fmt.Errorf("cannot calculate peg penalty due to value error: %w", err)
	}

	deviation := currentValue - wm.targetPeg
	penalty := deviation * deviation
	return penalty, nil
}

// TODO: Add methods related to updating weights based on Hamiltonian gradients.
// func (wm *WSIManager) ApplyWeightUpdates(gradients map[string]float64, eta, dt float64) error
