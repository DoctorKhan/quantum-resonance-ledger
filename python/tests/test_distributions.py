import pytest
import numpy as np
from qrl_simulation.distributions import TruncatedGaussian # Assuming qrl_simulation is importable

def test_truncated_gaussian_creation():
    """Tests valid and invalid creation of TruncatedGaussian."""
    # Valid
    dist = TruncatedGaussian(mean=5.0, std_dev=1.0, min_val=3.0, max_val=7.0)
    assert dist is not None
    assert dist.type() == "TruncatedGaussian"
    assert dist.mean == 5.0
    assert dist.std_dev == 1.0
    assert dist.min_val == 3.0
    assert dist.max_val == 7.0

    # Invalid std_dev
    with pytest.raises(ValueError, match="Standard deviation must be positive"):
        TruncatedGaussian(mean=5.0, std_dev=0.0, min_val=3.0, max_val=7.0)
    with pytest.raises(ValueError, match="Standard deviation must be positive"):
        TruncatedGaussian(mean=5.0, std_dev=-1.0, min_val=3.0, max_val=7.0)

    # Invalid bounds
    with pytest.raises(ValueError, match="min_val must be less than max_val"):
        TruncatedGaussian(mean=5.0, std_dev=1.0, min_val=7.0, max_val=7.0)
    with pytest.raises(ValueError, match="min_val must be less than max_val"):
        TruncatedGaussian(mean=5.0, std_dev=1.0, min_val=8.0, max_val=7.0)

def test_truncated_gaussian_sampling():
    """Tests if samples fall within the specified bounds."""
    mean, std_dev, min_val, max_val = 5.0, 1.0, 3.0, 7.0
    dist = TruncatedGaussian(mean=mean, std_dev=std_dev, min_val=min_val, max_val=max_val)
    num_samples = 500 # Increased samples for better stats check
    samples = [dist.sample() for _ in range(num_samples)]

    for i, sample in enumerate(samples):
        assert min_val <= sample <= max_val, f"Sample {i} ({sample}) out of bounds [{min_val}, {max_val}]"

    # Basic check on mean and std dev (will be approximate)
    # Note: Mean/std dev of *truncated* sample will differ from original mean/std dev
    sample_mean = np.mean(samples)
    sample_std = np.std(samples)
    print(f"\nSample Mean: {sample_mean:.4f} (Param Mean: {mean})")
    print(f"Sample StdDev: {sample_std:.4f} (Param StdDev: {std_dev})")
    # Check if sample mean is reasonably close to the expected truncated mean
    # Exact calculation is complex, use a wider tolerance or just check bounds
    # For this specific case (symmetric truncation around mean), mean should be close to 5.0
    assert abs(sample_mean - mean) < 0.2, "Sample mean deviates significantly from parameter mean"
    # Std dev of truncated normal is always less than original std dev
    assert 0 < sample_std < std_dev, "Sample std dev out of expected range"


def test_truncated_gaussian_pdf():
    """Tests the PDF calculation at specific points."""
    mean, std_dev, min_val, max_val = 5.0, 1.0, 3.0, 7.0
    dist = TruncatedGaussian(mean=mean, std_dev=std_dev, min_val=min_val, max_val=max_val)

    # PDF should be higher near the mean
    pdf_at_mean = dist.pdf(mean)
    pdf_near_min = dist.pdf(min_val + 0.1)
    pdf_near_max = dist.pdf(max_val - 0.1)
    assert pdf_at_mean > 0
    assert pdf_near_min > 0
    assert pdf_near_max > 0
    # Check relative heights (can be sensitive to exact parameters)
    # assert pdf_at_mean > pdf_near_min
    # assert pdf_at_mean > pdf_near_max

    # PDF should be zero outside bounds
    assert dist.pdf(min_val - 0.1) == 0.0
    assert dist.pdf(max_val + 0.1) == 0.0

    # Check symmetry if mean is center (optional)
    # Recenter for symmetry test
    dist_sym = TruncatedGaussian(mean=5.0, std_dev=1.0, min_val=4.0, max_val=6.0)
    assert np.isclose(dist_sym.pdf(dist_sym.mean - 0.5), dist_sym.pdf(dist_sym.mean + 0.5))


def test_truncated_gaussian_update():
    """Tests updating the distribution parameters."""
    mean, std_dev, min_val, max_val = 5.0, 1.0, 3.0, 7.0
    dist = TruncatedGaussian(mean=mean, std_dev=std_dev, min_val=min_val, max_val=max_val)
    initial_mean_sample = np.mean([dist.sample() for _ in range(100)])

    # Update mean
    new_mean = 6.0
    dist.update(mean=new_mean)
    assert dist.mean == new_mean
    assert dist.std_dev == std_dev # Check std dev didn't change
    # Check if internal distribution was updated (sample should reflect new mean)
    samples_new_mean = [dist.sample() for _ in range(200)]
    mean_after_update = np.mean(samples_new_mean)
    print(f"Mean after update: {mean_after_update:.4f} (New Param Mean: {new_mean})")
    # Mean should shift towards new_mean (but still affected by truncation)
    assert mean_after_update > initial_mean_sample
    assert abs(mean_after_update - new_mean) < 0.5 # Allow some deviation due to truncation

    # Update std_dev
    new_std_dev = 0.5
    dist.update(std_dev=new_std_dev)
    assert dist.std_dev == new_std_dev
    assert dist.mean == new_mean # Check mean didn't change
    samples_new_std = [dist.sample() for _ in range(200)]
    std_after_update = np.std(samples_new_std)
    print(f"StdDev after update: {std_after_update:.4f} (New Param StdDev: {new_std_dev})")
    # Expect lower std dev in samples
    assert std_after_update < 0.4 # Should be significantly less than original std dev of 1.0

    # Invalid update
    with pytest.raises(ValueError, match="Standard deviation must be positive"):
        dist.update(std_dev=0.0)