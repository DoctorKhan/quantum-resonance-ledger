import pytest
import numpy as np # For potential future use if needed
from src.node import Node
from src.parameters import Parameter
from src.distributions import TruncatedGaussian

@pytest.fixture
def sample_params():
    """Fixture for creating sample Parameter objects."""
    p1 = Parameter("fee_rate", TruncatedGaussian(0.01, 0.005, 0.001, 0.1))
    p2 = Parameter("block_size", TruncatedGaussian(1.0, 0.1, 0.5, 2.0))
    return {"fee_rate": p1, "block_size": p2}

def test_node_creation(sample_params):
    """Tests creating Node objects."""
    node_id = "NodeA"
    node = Node(node_id, sample_params)

    assert node.node_id == node_id
    assert node.token_balance == 100.0 # Default
    assert node.quantity_imbalance == 0.0 # Default
    assert len(node.neighbors) == 0
    assert len(node.parameters) == 2
    assert node.parameters["fee_rate"] == sample_params["fee_rate"]
    assert node.parameters["block_size"] == sample_params["block_size"]
    assert str(node) == f"Node(ID: {node_id}, Balance: 100.00, Imbalance: 0.0000)"
    assert repr(node) == f"Node('{node_id}')"

    # Test creation with no initial params
    node_no_params = Node("NodeB")
    assert node_no_params.node_id == "NodeB"
    assert len(node_no_params.parameters) == 0

    # Test invalid creation
    with pytest.raises(ValueError, match="Node ID cannot be empty"):
        Node("", sample_params)

def test_node_add_neighbor(sample_params):
    """Tests adding neighbors."""
    node = Node("NodeA", sample_params)
    node.add_neighbor("NodeB")
    node.add_neighbor("NodeC")
    node.add_neighbor("NodeB") # Add duplicate
    node.add_neighbor("NodeA") # Add self (should be ignored)
    node.add_neighbor("")      # Add empty (should be ignored)

    assert node.neighbors == {"NodeB", "NodeC"}

def test_node_parameter_access(sample_params):
    """Tests getting parameter objects and values."""
    node = Node("NodeA", sample_params)

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