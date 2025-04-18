package core

import (
	"testing"
)

// Helper function to create mock parameters for QSD tests
func createMockQSDParams() (*Parameter, *Parameter, *Parameter) {
	// Use mock distributions for parameters
	distCR := &mockDistribution{mean: 1.5, stdDev: 0.1} // Example Collateral Ratio
	paramCR := NewParameter("qsd_collateral_ratio", distCR)
	paramCR.CurrentValue = 1.5 // Set initial value

	distSF := &mockDistribution{mean: 0.02, stdDev: 0.005} // Example Stability Fee (2%)
	paramSF := NewParameter("qsd_stability_fee", distSF)
	paramSF.CurrentValue = 0.02

	distLP := &mockDistribution{mean: 0.1, stdDev: 0.02} // Example Liquidation Penalty (10%)
	paramLP := NewParameter("qsd_liquidation_penalty", distLP)
	paramLP.CurrentValue = 0.1

	return paramCR, paramSF, paramLP
}

func TestQSD_Minting(t *testing.T) {
	paramCR, paramSF, paramLP := createMockQSDParams()
	manager := NewQSDManager(paramCR, paramSF, paramLP)
	owner := "userA"
	collateralType := "qETH"
	collateralID := Hash{1, 2, 3}                   // Mock CUT ID
	collateralAmount := uint64(2000000000000000000) // 2 qETH
	qsdToMint := uint64(1000)                       // Mint 1000 QSD

	t.Run("Mint", func(t *testing.T) {
		err := manager.Mint(owner, collateralType, collateralID, collateralAmount, qsdToMint)
		if err != nil {
			t.Errorf("Expected Mint function to return no error, but got %v", err)
		}

		// Basic check if vault was created
		manager.mu.RLock()
		vaultExists := false
		if ownerVaults, ok := manager.vaults[owner]; ok {
			if _, ok := ownerVaults[collateralType]; ok {
				vaultExists = true
			}
		}
		manager.mu.RUnlock()
		if !vaultExists {
			t.Errorf("Mint logic did not create the vault entry")
		}

		// Call Mint with placeholder owner and expect an error
		err = manager.Mint("placeholder", collateralType, collateralID, collateralAmount, qsdToMint)
		if err == nil {
			t.Errorf("Expected Mint function to return an error for placeholder owner, but got nil")
		} else {
			t.Logf("Got expected error from Mint for placeholder owner: %v", err)
		}
		// TODO: Add checks for collateral ratio, state updates etc. once implemented
	})

	t.Run("MintInvalidInputs", func(t *testing.T) {
		// Zero amount
		err := manager.Mint(owner, collateralType, collateralID, collateralAmount, 0)
		if err == nil {
			t.Errorf("Expected Mint function to return an error for zero amount, but got nil")
		} else {
			t.Logf("Got expected error from Mint for zero amount: %v", err)
		}
	})

	t.Run("MintInsufficientCollateral", func(t *testing.T) {
		// Insufficient collateral
		qsdToMint := uint64(10000) // Mint 10000 QSD
		err := manager.Mint(owner, collateralType, collateralID, collateralAmount, qsdToMint)
		if err == nil {
			t.Errorf("Expected Mint function to return an error for insufficient collateral, but got nil")
		} else {
			t.Logf("Got expected error from Mint for insufficient collateral: %v", err)
		}
	})

	// TODO: Add tests for insufficient collateral once implemented
}

func TestQSD_Burning(t *testing.T) {
	paramCR, paramSF, paramLP := createMockQSDParams()
	manager := NewQSDManager(paramCR, paramSF, paramLP)
	owner := "userA"
	collateralType := "qETH"
	collateralID := Hash{1, 2, 3}
	collateralAmount := uint64(2000000000000000000)
	qsdMinted := uint64(1000)
	qsdToRepay := uint64(500)

	// Setup: Mint first (using placeholder which creates vault entry)
	_ = manager.Mint(owner, collateralType, collateralID, collateralAmount, qsdMinted)

	t.Run("Burn", func(t *testing.T) {
		err := manager.Burn(owner, collateralType, qsdToRepay)
		if err != nil {
			t.Errorf("Expected Burn function to return no error, but got %v", err)
		}

		// Call Burn with placeholder owner and expect an error
		err = manager.Burn("placeholder", collateralType, qsdToRepay)
		if err == nil {
			t.Errorf("Expected Burn function to return an error for placeholder owner, but got nil")
		} else {
			t.Logf("Got expected error from Burn for placeholder owner: %v", err)
		}
		// TODO: Add checks for debt reduction, collateral unlock etc. once implemented
	})

	t.Run("BurnNonExistentVault", func(t *testing.T) {
		err := manager.Burn("nonExistentUser", collateralType, qsdToRepay)
		if err == nil {
			t.Errorf("Expected error when burning from non-existent user vault, but got nil")
		}
		err = manager.Burn(owner, "qBTC", qsdToRepay) // Correct owner, wrong collateral type
		if err == nil {
			t.Errorf("Expected error when burning from non-existent collateral type vault, but got nil")
		}

		// Zero amount
		err = manager.Burn(owner, collateralType, 0)
		if err == nil {
			t.Errorf("Expected Burn function to return an error for zero amount, but got nil")
		} else {
			t.Logf("Got expected error from Burn for zero amount: %v", err)
		}
	})

	// TODO: Add tests for repaying full amount, partial amount, stability fees etc.
}

// TODO: Add TestQSD_LiquidationTrigger
// TODO: Add TestQSD_LiquidationProcess
// TODO: Add TestQSD_CollateralRatioCheck
// TODO: Add TestQSD_StabilityFee
