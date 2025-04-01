from typing import Any
# Import the base Distribution class (or specific ones if needed)
from .distributions import Distribution, TruncatedGaussian

class Parameter:
    """Represents a simulation parameter governed by a probability distribution."""

    def __init__(self, name: str, distribution: Distribution):
        if not name:
            # Decide if empty names are allowed - let's disallow for now
            raise ValueError("Parameter name cannot be empty.")
        if distribution is None:
            # Decide if nil distribution is allowed - let's disallow for now
            raise ValueError("Parameter must have a distribution.")

        self.name = name
        self.distribution = distribution

    def sample(self) -> float:
        """Draw a sample from the parameter's distribution."""
        return self.distribution.sample()

    def pdf(self, x: float) -> float:
        """Get the PDF value from the parameter's distribution."""
        return self.distribution.pdf(x)

    def update_distribution(self, **kwargs: Any):
        """
        Update the parameters of the underlying distribution.
        Delegates to the distribution's update method if it exists.
        """
        if hasattr(self.distribution, 'update') and callable(self.distribution.update):
             # Pass only the arguments relevant to the distribution's update method
             # This requires knowing the signature or using inspection,
             # or simply passing all kwargs and letting the distribution handle it.
             # For now, assume distribution.update handles arbitrary kwargs.
             self.distribution.update(**kwargs)
        # If update doesn't exist, we might not raise an error, just do nothing silently,
        # or log a warning, depending on desired behavior.
        # else:
        #     print(f"Warning: Distribution type {type(self.distribution)} does not support update.")


    def get_mean(self) -> float:
        """Returns the mean of the distribution, if available."""
        if hasattr(self.distribution, 'mean'):
            return self.distribution.mean
        # May need to calculate if not stored (e.g., for Beta)
        # For now, return NaN or raise error if not directly available
        return float('nan')

    def get_std_dev(self) -> float:
        """Returns the standard deviation of the distribution, if available."""
        if hasattr(self.distribution, 'std_dev'):
            return self.distribution.std_dev
        # May need to calculate if not stored
        return float('nan')

    def __str__(self) -> str:
        dist_type = self.distribution.type() if self.distribution else "None"
        return f"Parameter(Name: {self.name}, Distribution: {dist_type})"