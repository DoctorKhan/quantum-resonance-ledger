import pytest
import numpy as np # For potential future use if needed
from qrl_simulation.node import Node
from qrl_simulation.parameters import Parameter
from qrl_simulation.distributions import TruncatedGaussian

@pytest.fixture
def sample_params():
    """Fixture for creating sample Parameter objects."""
    p1 = Parameter("fee_rate", TruncatedGaussian(0.01, 0.005, 0.001, 0.1))
    p2 = Parameter("block_size", TruncatedGaussian(1.0, 0.1, 0.5, 2.0))
    return {"fee_rate": p1, "block_size": p2}

def test_node_creation(sample_params):
    """Tests creating Node objects according to TDD plan."""
    node_id = "NodeA"
    position = (10.0, 20.5)
    node = Node(node_id, position=position, initial_parameters=sample_params)

    assert node.node_id == node_id
    assert node.position == position
    # Check initial zero balances per TDD plan
    assert node.balances == {"QUSD": 0.0, "QRG": 0.0, "Gas": 0.0}
    assert node.quantity_imbalance == 0.0 # Changed from {} to 0.0
    assert node.buy_propensity_fields == {}
    assert node.sell_propensity_fields == {}
    assert len(node.neighbors) == 0
    assert len(node.parameters) == 2
    assert node.parameters["fee_rate"] == sample_params["fee_rate"]
    assert node.parameters["block_size"] == sample_params["block_size"]
    # Check updated string representation
    expected_str = (f"Node(ID: {node_id}, Pos: {position}, "
    				            f"Balances: [Gas: 0.00, QRG: 0.00, QUSD: 0.00], Q Imbalance: 0.0000)") # Updated imbalance format
    assert str(node) == expected_str
    assert repr(node) == f"Node('{node_id}')"

    # Test creation with no initial params
    node_no_params = Node("NodeB", position=(1, 1)) # Position is optional but good to test
    assert node_no_params.node_id == "NodeB"
    assert node_no_params.position == (1, 1)
    assert node_no_params.balances == {"QUSD": 0.0, "QRG": 0.0, "Gas": 0.0}
    assert len(node_no_params.parameters) == 0

    # Test invalid creation
    with pytest.raises(ValueError, match="Node ID cannot be empty"):
        Node("", position=(0,0), initial_parameters=sample_params)

def test_node_add_neighbor(sample_params):
    """Tests adding neighbors."""
    node = Node("NodeA", position=(0,0), initial_parameters=sample_params)
    node.add_neighbor("NodeB")
    node.add_neighbor("NodeC")
    node.add_neighbor("NodeB") # Add duplicate
    node.add_neighbor("NodeA") # Add self (should be ignored)
    node.add_neighbor("")      # Add empty (should be ignored)

    assert node.neighbors == {"NodeB", "NodeC"}

def test_node_parameter_access(sample_params):
    """Tests getting parameter objects and values."""
    node = Node("NodeA", position=(0,0), initial_parameters=sample_params)

    # Get Parameter object
    fee_param = node.get_parameter("fee_rate")
    assert isinstance(fee_param, Parameter)
    assert fee_param.name == "fee_rate"
    assert node.get_parameter("non_existent") is None

    # Get parameter mean value using the specific getter
    assert node.get_parameter_value("fee_rate") == 0.01
    assert node.get_parameter_value("block_size") == 1.0
    assert node.get_parameter_value("non_existent") is None

    # Sample parameter
    fee_sample = node.sample_parameter("fee_rate")
    assert isinstance(fee_sample, float)
    assert 0.001 <= fee_sample <= 0.1 # Check bounds
    assert node.sample_parameter("non_existent") is None
   
# --- RTT Tests ---

def test_node_rtt_initialization():
    """Tests that RTT fields are initialized correctly."""
    node = Node("NodeRTT")
    assert node.buy_propensity_fields == {}
    assert node.sell_propensity_fields == {}
    assert node.quantity_imbalance == 0.0 # Changed from {} to 0.0

def test_node_perturb_propensity_field():
    """Tests the perturb_propensity_field method."""
    node = Node("NodeRTT")
    asset_id = "QRG"
    price_range = (10.0, 11.0)
    key = node._get_price_range_key(price_range)
    magnitude1 = 0.5
    magnitude2 = 0.3

    # 1. Perturb buy field
    node.perturb_propensity_field(asset_id, price_range, magnitude1, is_buy=True)
    assert asset_id in node.buy_propensity_fields
    assert key in node.buy_propensity_fields[asset_id]
    assert node.buy_propensity_fields[asset_id][key] == pytest.approx(magnitude1)
    assert node.sell_propensity_fields == {} # Sell field unchanged

    # 2. Perturb buy field again (accumulation)
    node.perturb_propensity_field(asset_id, price_range, magnitude2, is_buy=True)
    assert node.buy_propensity_fields[asset_id][key] == pytest.approx(magnitude1 + magnitude2)

    # 3. Perturb sell field
    magnitude3 = 0.4
    node.perturb_propensity_field(asset_id, price_range, magnitude3, is_buy=False)
    assert asset_id in node.sell_propensity_fields
    assert key in node.sell_propensity_fields[asset_id]
    assert node.sell_propensity_fields[asset_id][key] == pytest.approx(magnitude3)
    assert node.buy_propensity_fields[asset_id][key] == pytest.approx(magnitude1 + magnitude2) # Buy field unchanged

    # 4. Test invalid inputs
    with pytest.raises(ValueError, match="Perturbation magnitude cannot be negative"):
        node.perturb_propensity_field(asset_id, price_range, -0.1, is_buy=True)

    with pytest.raises(ValueError, match="Invalid price_range format"):
        node.perturb_propensity_field(asset_id, (10.0,), 0.1, is_buy=True) # Not a tuple of 2

    with pytest.raises(ValueError, match="Invalid price_range format"):
        node.perturb_propensity_field(asset_id, (11.0, 10.0), 0.1, is_buy=True) # min > max

    with pytest.raises(ValueError, match="Asset ID cannot be empty"):
        node.perturb_propensity_field("", price_range, 0.1, is_buy=True)

# TODO: Add tests for attempt_local_settlement once implemented
# TODO: Add tests for propagate_rtt_state once implemented

# --- End RTT Tests ---