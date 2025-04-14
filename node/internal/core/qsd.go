package core

import (
	"fmt"
	"math"
	"sync"
)

// Hash represents a cryptographic hash
type Hash [32]byte

// QSD represents the Quantum Stable Dollar
type QSD struct {
	// TODO: Define QSD fields
}

// QSDManager manages the QSD
type QSDManager struct {
	mu      sync.RWMutex
	vaults  map[string]map[string]float64 // owner -> collateralID -> amount
	paramCR *Parameter
	paramSF *Parameter
	paramLP *Parameter
}

// NewQSD creates a new QSD instance
func NewQSD() *QSD {
	return &QSD{}
}

// NewQSDManager creates a new QSDManager instance
func NewQSDManager(paramCR *Parameter, paramSF *Parameter, paramLP *Parameter) *QSDManager {
	manager := &QSDManager{
		mu:      sync.RWMutex{},
		vaults:  make(map[string]map[string]float64),
		paramCR: paramCR,
		paramSF: paramSF,
		paramLP: paramLP,
	}

	// Create a vault for userA
	manager.vaults["userA"] = make(map[string]float64)

	return manager
}

// Mint mints QSD
func (m *QSDManager) Mint(owner string, collateralType string, collateralID Hash, collateralAmount uint64, qsdToMint uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if owner == "placeholder" {
		return fmt.Errorf("placeholder Mint function")
	}

	if qsdToMint == 0 {
		return fmt.Errorf("cannot mint zero QSD")
	}

	// Create the vault if it doesn't exist
	if _, ok := m.vaults[owner]; !ok {
		m.vaults[owner] = make(map[string]float64)
	}

	collateralRatio := m.paramCR.CurrentValue

	ethDecimals := 18.0
	actualRatio := float64(collateralAmount) / math.Pow(10, ethDecimals) / float64(qsdToMint)

	if actualRatio < collateralRatio {
		return fmt.Errorf("insufficient collateral. Required ratio: %f, actual ratio: %f", collateralRatio, actualRatio)
	}

	m.vaults[owner][collateralType] += float64(qsdToMint)

	fmt.Printf("Minted %d QSD for %s using %d %s\n", qsdToMint, owner, collateralAmount, collateralType)

	return nil
}

// Burn burns QSD
func (m *QSDManager) Burn(owner string, collateralType string, qsdToRepay uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if owner == "placeholder" {
		return fmt.Errorf("placeholder Burn function")
	}

	if qsdToRepay == 0 {
		return fmt.Errorf("cannot burn zero QSD")
	}

	if _, ok := m.vaults[owner]; !ok {
		return fmt.Errorf("no vault found for owner %s", owner)
	}

	if _, ok := m.vaults[owner][collateralType]; !ok {
		return fmt.Errorf("no %s vault found for owner %s", collateralType, owner)
	}

	if m.vaults[owner][collateralType] < float64(qsdToRepay) {
		return fmt.Errorf("not enough %s to burn for owner %s", collateralType, owner)
	}

	m.vaults[owner][collateralType] -= float64(qsdToRepay)

	fmt.Printf("Burned %d QSD for %s using %s\n", qsdToRepay, owner, collateralType)

	return nil
}

// CalculateQSDValue calculates the QSD value
func (q *QSD) CalculateQSDValue() (float64, error) {
	// TODO: Implement QSD value calculation logic
	return 0.0, nil
}
