import pytest
import numpy as np
from qrl_simulation.quantum_blockchain import QuantumBlockchain
from qrl_simulation.node import Node
from qrl_simulation.parameters import Parameter
from qrl_simulation.distributions import TruncatedGaussian

# Assuming QSD functionality will be added to QuantumBlockchain and Node

@pytest.fixture
def qsd_chain():
    """Fixture to create a basic QuantumBlockchain instance for QSD tests."""
    initial_nodes = ["Alice", "Bob"]
    chain = QuantumBlockchain(initial_nodes=initial_nodes)
    
    # Give Alice some collateral (e.g., qETH, assuming it's tracked like other balances)
    # We might need a more formal way to handle different asset types later
    if "qETH" not in chain.nodes["Alice"].balances:
         chain.nodes["Alice"].balances["qETH"] = 0.0 # Initialize if not present
    chain.nodes["Alice"].balances["qETH"] += 10.0 
    
    # Initialize QSD balance if it doesn't exist
    if "QSD" not in chain.nodes["Alice"].balances:
        chain.nodes["Alice"].balances["QSD"] = 0.0
        
    # Assume some initial parameters needed for minting (will be refined)
    # Example: A collateral ratio parameter
    if 'collateral_ratio' not in chain.params:
        # Create the Parameter object without initial_value
        collateral_ratio_param_obj = Parameter(
            name='collateral_ratio',
            distribution=TruncatedGaussian(mean=1.5, std_dev=0.1, min_val=1.1, max_val=3.0)
        )
        # Store the object if needed (e.g., if nodes store Parameter objects directly)
        # chain.nodes["Alice"].parameters['collateral_ratio'] = collateral_ratio_param_obj
        # chain.nodes["Bob"].parameters['collateral_ratio'] = collateral_ratio_param_obj

        # Set the initial value in the chain's legacy params structure
        chain.params['collateral_ratio'] = {
             'value': {},
             'min': 1.1,
             'max': 3.0
        }
        # Initialize node-specific values in the legacy params structure
        initial_collateral_ratio_value = 1.5 # Match the mean used in TruncatedGaussian
        for node_id in chain.nodes:
             if node_id not in chain.params['collateral_ratio']['value']:
                  chain.params['collateral_ratio']['value'][node_id] = initial_collateral_ratio_value


    # Assume an oracle price for collateral (e.g., 1 qETH = $2000)
    # This might need a dedicated oracle mechanism later
    chain.oracle_prices = {"qETH": 2000.0}

    return chain

def test_mint_qsd_basic(qsd_chain):
    """
    Test basic QSD minting by locking collateral.
    This test WILL FAIL until mint_qsd and vault logic are implemented.
    """
    chain = qsd_chain
    alice_node = chain.nodes["Alice"]
    initial_qsd_balance = alice_node.get_balance("QSD")
    initial_qeth_balance = alice_node.get_balance("qETH")
    
    collateral_to_lock = 2.0
    qsd_to_mint = 2500.0 # e.g., 2 ETH * $2000/ETH / 1.6 collateral ratio (example)

    # --- This method needs to be implemented ---
    try:
        chain.mint_qsd(
            node_id="Alice", 
            collateral_asset="qETH", 
            collateral_amount=collateral_to_lock, 
            qsd_amount_to_mint=qsd_to_mint
        )
    except AttributeError:
        pytest.fail("QuantumBlockchain.mint_qsd() method not implemented yet.")
    except Exception as e:
         pytest.fail(f"mint_qsd raised an unexpected exception: {e}")

    # --- Assertions (will fail initially) ---
    # 1. QSD balance increased
    assert alice_node.get_balance("QSD") == initial_qsd_balance + qsd_to_mint, "QSD balance did not increase correctly."

    # 2. Collateral balance decreased (or moved to a 'locked' state)


def test_mint_qsd_insufficient_collateral(qsd_chain):
    """
    Test that minting fails if the node does not have enough collateral asset.
    """
    chain = qsd_chain
    alice_node = chain.nodes["Alice"]
    initial_qeth_balance = alice_node.get_balance("qETH") # Should be 10.0 from fixture
    
    collateral_to_lock = initial_qeth_balance + 1.0 # More than Alice has
    qsd_to_mint = 1000.0 # Amount doesn't matter as much as collateral check

    with pytest.raises(ValueError):
        chain.mint_qsd(
            node_id="Alice",
            collateral_asset="qETH",
            collateral_amount=collateral_to_lock,
            qsd_amount_to_mint=qsd_to_mint
        )

    # Verify balances haven't changed
    # Get initial QSD balance *before* the try block might be safer, but should be 0 here.
    initial_qsd_balance = 0.0
    assert alice_node.get_balance("qETH") == initial_qeth_balance, "Collateral balance changed unexpectedly."
    assert alice_node.get_balance("QSD") == initial_qsd_balance, "QSD balance changed unexpectedly."

    # 3. Vault/Position created (needs implementation)
    # assert "Alice_qETH_vault_1" in chain.vaults # Example assertion
    # vault = chain.vaults["Alice_qETH_vault_1"]
    # assert vault.collateral_asset == "qETH"
    # assert vault.collateral_amount == collateral_to_lock
    # assert vault.debt_amount == qsd_to_mint


def test_mint_qsd_insufficient_collateral_value(qsd_chain):
    """
    Test that minting fails if the provided collateral value is too low
    for the requested QSD amount based on the collateral ratio.
    """
    chain = qsd_chain
    alice_node = chain.nodes["Alice"]
    initial_qeth_balance = alice_node.get_balance("qETH") # 10.0
    initial_qsd_balance = alice_node.get_balance("QSD") # 0.0
    collateral_price = chain.oracle_prices["qETH"] # 2000.0
    # Get the required ratio (should be 1.5 from fixture)
    required_ratio = next(iter(chain.params['collateral_ratio']['value'].values()), 1.5)

    collateral_to_lock = 1.0 # Alice has 10.0, so balance is sufficient
    collateral_value = collateral_to_lock * collateral_price # 1.0 * 2000.0 = 2000.0
    
    # Calculate max QSD mintable with this collateral value and ratio
    max_qsd_mintable = collateral_value / required_ratio # 2000.0 / 1.5 = 1333.33
    qsd_to_mint = max_qsd_mintable + 100.0 # Request slightly more than allowed

    with pytest.raises(ValueError, match="Insufficient collateral value provided"):
        chain.mint_qsd(
            node_id="Alice",
            collateral_asset="qETH",
            collateral_amount=collateral_to_lock,
            qsd_amount_to_mint=qsd_to_mint
        )

    # Verify balances haven't changed
    assert alice_node.get_balance("qETH") == initial_qeth_balance
    assert alice_node.get_balance("QSD") == initial_qsd_balance

def test_redeem_qsd_basic(qsd_chain):
    """
    Test basic QSD redemption to unlock collateral.
    Assumes a simple 1:1 redemption rate against the collateral's oracle price.
    """
    chain = qsd_chain
    alice_node = chain.nodes["Alice"]
    collateral_asset = "qETH"
    collateral_price = chain.oracle_prices[collateral_asset] # 2000.0

    # --- Setup: Mint some QSD first ---
    collateral_to_lock = 2.0
    # Mint slightly less than max possible to avoid rounding issues in test
    # Max = (2.0 * 2000.0) / 1.5 = 2666.66
    qsd_to_mint = 2500.0 
    chain.mint_qsd(
        node_id="Alice", 
        collateral_asset=collateral_asset, 
        collateral_amount=collateral_to_lock, 
        qsd_amount_to_mint=qsd_to_mint
    )
    initial_qeth_balance = alice_node.get_balance(collateral_asset) # Should be 10.0
    initial_qsd_balance = alice_node.get_balance("QSD") # Should be 0.0

    initial_qeth_balance_after_mint = alice_node.get_balance(collateral_asset) # Should be 10.0 - 2.0 = 8.0
    initial_qsd_balance_after_mint = alice_node.get_balance("QSD") # Should be 2500.0
    assert initial_qsd_balance_after_mint > 0

    initial_qeth_balance = initial_qeth_balance_after_mint # Use balance *after* minting as initial for redemption test
    initial_qsd_balance = initial_qsd_balance_after_mint
    assert initial_qsd_balance > 0

    # --- Test Redemption ---
    qsd_to_redeem = 1000.0
    # Calculate expected collateral to unlock (simplified: QSD / price)
    # A real system might use vault data and collateral ratio
    expected_collateral_unlocked = qsd_to_redeem / collateral_price # 1000.0 / 2000.0 = 0.5

    # --- This method needs to be implemented ---
    try:
        chain.redeem_qsd(
            node_id="Alice", 
            collateral_asset=collateral_asset, # Need to specify which collateral type to receive
            qsd_amount_to_redeem=qsd_to_redeem
            # In a vault system, might specify vault ID instead of asset/amount
        )
    except AttributeError:
        pytest.fail("QuantumBlockchain.redeem_qsd() method not implemented yet.")
    except Exception as e:
         pytest.fail(f"redeem_qsd raised an unexpected exception: {e}")

    # --- Assertions (will fail initially) ---
    # 1. QSD balance decreased
    assert alice_node.get_balance("QSD") == initial_qsd_balance - qsd_to_redeem, "QSD balance did not decrease correctly."

    # 2. Collateral balance increased
    assert alice_node.get_balance(collateral_asset) == initial_qeth_balance + expected_collateral_unlocked, "Collateral balance did not increase correctly."

    # 3. Vault/Position updated (needs implementation)
    # assert chain.vaults["Alice_qETH_vault_1"].debt_amount == qsd_to_mint - qsd_to_redeem
    # assert chain.vaults["Alice_qETH_vault_1"].collateral_amount == collateral_to_lock - expected_collateral_unlocked



def test_qsd_vault_creation(qsd_chain):
    """
    Test basic QSD vault creation.
    """
    chain = qsd_chain
    alice_node = chain.nodes["Alice"]
    collateral_asset = "qETH"

    # --- Test Vault Creation ---
    vault_id = "alice_qeth_vault_1" # Example vault ID
    collateral_amount = 5.0

    try:
        chain.create_vault(
            node_id="Alice",
            vault_id=vault_id,
            collateral_asset=collateral_asset,
            collateral_amount=collateral_amount
        )
    except AttributeError:
        pytest.fail("QuantumBlockchain.create_vault() method not implemented yet.")
    except Exception as e:
         pytest.fail(f"create_vault raised an unexpected exception: {e}")

    # --- Assertions (will fail initially) ---
    # 1. Vault exists in chain.vaults
    assert vault_id in chain.vaults, f"Vault {vault_id} not found in chain.vaults."

    vault = chain.vaults[vault_id]

    # 2. Vault attributes are set correctly
    assert vault.owner == "Alice"
    assert vault.collateral_asset == collateral_asset
    assert vault.collateral_amount == collateral_amount
    assert vault.debt_amount == 0.0 # Initially zero debt
    # assert vault.creation_time is not None # Add time check later if needed

    # 3. Vault is associated with the node (e.g., in node.vaults list)


def test_mint_qsd_against_vault(qsd_chain):
    """
    Test minting QSD against a specific, existing vault.
    """
    chain = qsd_chain
    alice_node = chain.nodes["Alice"]
    collateral_asset = "qETH"
    collateral_price = chain.oracle_prices[collateral_asset]
    required_ratio = next(iter(chain.params['collateral_ratio']['value'].values()), 1.5)

    # --- Setup: Create a vault first ---
    vault_id = "alice_qeth_vault_2"
    collateral_to_lock = 5.0
    initial_qeth_balance = alice_node.get_balance(collateral_asset) # 10.0
    chain.create_vault(
        node_id="Alice",
        vault_id=vault_id,
        collateral_asset=collateral_asset,
        collateral_amount=collateral_to_lock
    )
    qeth_balance_after_vault = alice_node.get_balance(collateral_asset) # Should be 5.0
    assert qeth_balance_after_vault == initial_qeth_balance - collateral_to_lock
    vault = chain.vaults[vault_id]
    assert vault.collateral_amount == collateral_to_lock
    assert vault.debt_amount == 0.0

    # --- Test Minting Against Vault ---
    initial_qsd_balance = alice_node.get_balance("QSD") # Should be 0.0
    
    # Calculate max mintable for this vault
    max_qsd_mintable = vault.collateral_amount * collateral_price / required_ratio # 5.0 * 2000 / 1.5 = 6666.66
    qsd_to_mint = max_qsd_mintable / 2 # Mint half the max possible

    # --- This method needs modification to accept vault_id ---
    try:
        chain.mint_qsd(
            node_id="Alice", 
            vault_id=vault_id, # NEW parameter
            qsd_amount_to_mint=qsd_to_mint
            # collateral_asset/amount no longer needed if using vault_id?
        )
    except TypeError: # Expecting TypeError initially due to wrong arguments
        pytest.xfail("mint_qsd signature needs update for vault_id.") 
    except AttributeError: # Fallback if method doesn't exist (shouldn't happen)
         pytest.fail("QuantumBlockchain.mint_qsd() method not found.")
    except Exception as e:
         pytest.fail(f"mint_qsd raised an unexpected exception: {e}")

    # --- Assertions (will fail initially) ---
    # 1. QSD balance increased
    assert alice_node.get_balance("QSD") == initial_qsd_balance + qsd_to_mint, "QSD balance did not increase correctly."

    # 2. Vault debt increased
    assert vault.debt_amount == qsd_to_mint, "Vault debt did not increase correctly."

    # 3. Vault collateral unchanged
    assert vault.collateral_amount == collateral_to_lock, "Vault collateral changed unexpectedly."

    # 4. Node's *available* collateral balance unchanged during mint
    assert alice_node.get_balance(collateral_asset) == qeth_balance_after_vault, "Node available collateral changed unexpectedly during mint."


    # assert vault_id in alice_node.vaults # Example assertion - node.vaults not yet implemented


    # assert vault.owner == "Alice"