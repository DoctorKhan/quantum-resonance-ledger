import pytest
import numpy as np # For np.isclose
from qrl_simulation.parameters import Parameter
from qrl_simulation.distributions import TruncatedGaussian # Use a concrete distribution for testing

# Fixture for a common distribution instance
@pytest.fixture
def gaussian_dist():
    return TruncatedGaussian(mean=5.0, std_dev=1.0, min_val=3.0, max_val=7.0)

def test_parameter_creation(gaussian_dist):
    """Tests creating Parameter objects."""
    param = Parameter("TestParam", gaussian_dist)

    assert param.name == "TestParam"
    assert param.distribution == gaussian_dist
    assert isinstance(param.distribution, TruncatedGaussian)
    assert str(param) == "Parameter(Name: TestParam, Distribution: TruncatedGaussian)"

    # Test invalid creation
    with pytest.raises(ValueError, match="Parameter name cannot be empty"):
        Parameter("", gaussian_dist)
    with pytest.raises(ValueError, match="Parameter must have a distribution"):
        Parameter("NoDistParam", None)

def test_parameter_sampling_pdf(gaussian_dist):
    """Tests delegation of sample() and pdf() to the distribution."""
    param = Parameter("TestParam", gaussian_dist)

    # Test sampling (just check bounds, detailed sampling tested elsewhere)
    sample = param.sample()
    assert 3.0 <= sample <= 7.0

    # Test PDF
    pdf_val = param.pdf(5.0)
    assert pdf_val > 0
    assert param.pdf(2.0) == 0.0 # Outside bounds

def test_parameter_update_distribution(gaussian_dist):
    """Tests delegation of update to the distribution."""
    param = Parameter("TestParam", gaussian_dist)

    # Check initial mean
    assert param.get_mean() == 5.0

    # Update mean via parameter
    param.update_distribution(mean=6.0)
    assert param.get_mean() == 6.0
    # Check underlying distribution was updated
    assert gaussian_dist.mean == 6.0

    # Update std_dev via parameter
    param.update_distribution(std_dev=0.5)
    assert param.get_std_dev() == 0.5
    assert gaussian_dist.std_dev == 0.5

    # Update both
    param.update_distribution(mean=5.5, std_dev=0.8)
    assert param.get_mean() == 5.5
    assert param.get_std_dev() == 0.8

    # Test invalid update delegation
    with pytest.raises(ValueError):
        param.update_distribution(std_dev=0.0)

def test_parameter_getters(gaussian_dist):
    """Tests getters for distribution properties."""
    param = Parameter("TestParam", gaussian_dist)
    assert param.get_mean() == 5.0
    assert param.get_std_dev() == 1.0

    # TODO: Test with a distribution that doesn't store mean/std_dev directly if needed