import numpy as np
from scipy.stats import truncnorm
from typing import Optional

class Distribution:
    """Base class for probability distributions (optional)."""
    def sample(self) -> float:
        raise NotImplementedError

    def pdf(self, x: float) -> float:
        raise NotImplementedError

    def type(self) -> str:
        raise NotImplementedError

class TruncatedGaussian(Distribution):
    """
    Represents a Gaussian distribution truncated to [min_val, max_val].
    Uses scipy.stats.truncnorm internally.
    """
    def __init__(self, mean: float, std_dev: float, min_val: float, max_val: float):
        if std_dev <= 0:
            raise ValueError("Standard deviation must be positive.")
        if min_val >= max_val:
            raise ValueError("min_val must be less than max_val.")

        self.mean = mean
        self.std_dev = std_dev
        self.min_val = min_val
        self.max_val = max_val

        # Parameters for truncnorm: a, b are (min-mean)/std, (max-mean)/std
        self.a = (min_val - mean) / std_dev
        self.b = (max_val - mean) / std_dev

        # Create the frozen distribution object
        self._dist = truncnorm(a=self.a, b=self.b, loc=mean, scale=std_dev)

    def sample(self) -> float:
        """Draw a random sample from the distribution."""
        return self._dist.rvs()

    def pdf(self, x: float) -> float:
        """Calculate the probability density function at point x."""
        return self._dist.pdf(x)

    def type(self) -> str:
        """Return the type identifier."""
        return "TruncatedGaussian"

    def update(self, mean: Optional[float] = None, std_dev: Optional[float] = None):
        """Update distribution parameters (mean, std_dev). Min/Max remain fixed."""
        if mean is not None:
            self.mean = mean
        if std_dev is not None:
            if std_dev <= 0:
                raise ValueError("Standard deviation must be positive.")
            self.std_dev = std_dev

        # Recalculate a, b and recreate the frozen distribution
        self.a = (self.min_val - self.mean) / self.std_dev
        self.b = (self.max_val - self.mean) / self.std_dev
        self._dist = truncnorm(a=self.a, b=self.b, loc=self.mean, scale=self.std_dev)

# Add BetaDistribution, etc. later