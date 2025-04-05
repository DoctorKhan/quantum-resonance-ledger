from typing import Dict, List, Set, Optional, Tuple
from .parameters import Parameter
from .distributions import Distribution # May be needed for type hints or defaults

class Node:
    """Represents a node in the QRL network simulation."""

    def __init__(self, node_id: str, position: Optional[Tuple[float, float]] = None, initial_parameters: Optional[Dict[str, Parameter]] = None):
        if not node_id:
            raise ValueError("Node ID cannot be empty.")

        self.node_id: str = node_id
        self.position: Optional[Tuple[float, float]] = position # Added position attribute
        # Parameters governing this node's behavior (e.g., fee rate dist)
        self.parameters: Dict[str, Parameter] = initial_parameters if initial_parameters else {}

        # Network state / relationships
        self.neighbors: Set[str] = set() # IDs of neighboring nodes

        # Ledger state (simplified for now)
        # Balances for different token types (initialized to zero per TDD plan)
        self.balances: Dict[str, float] = {"QUSD": 0.0, "QRG": 0.0, "Gas": 0.0}
        # TODO: Track imbalance per token type
        self.quantity_imbalance: float = 0.0

        # Simulation state
        # TODO: Add inbox for event-driven simulation
        # self.inbox: List[Event] = []
        # self.processed_event_count: int = 0

        # Add other state as needed (e.g., local CUT commitments)

    def add_neighbor(self, neighbor_id: str):
        """Adds a neighbor relationship."""
        if neighbor_id and neighbor_id != self.node_id:
            self.neighbors.add(neighbor_id)

    def get_parameter(self, name: str) -> Optional[Parameter]:
        """Gets a parameter object by name."""
        return self.parameters.get(name)

    def get_parameter_value(self, name: str) -> Optional[float]:
        """Gets the current mean value of a parameter's distribution."""
        param = self.get_parameter(name)
        # Use get_mean() which handles missing attribute gracefully
        return param.get_mean() if param else None

    def sample_parameter(self, name: str) -> Optional[float]:
        """Draws a sample from a parameter's distribution."""
        param = self.get_parameter(name)
        return param.sample() if param else None


    def get_balance(self, token: str) -> float:
        """Returns the balance for a specific token, defaulting to 0.0."""
        return self.balances.get(token, 0.0)

    def increase_balance(self, token: str, amount: float):
        """Increases the balance for a specific token."""
        if amount < 0:
            raise ValueError("Cannot increase balance by a negative amount.")
        self.balances[token] = self.balances.get(token, 0.0) + amount

    def decrease_balance(self, token: str, amount: float):
        """Decreases the balance for a specific token."""
        if amount < 0:
            raise ValueError("Cannot decrease balance by a negative amount. Use increase_balance.")
        current_balance = self.balances.get(token, 0.0)
        # For now, allow negative balances, but ideally check for sufficient funds
        # if current_balance < amount:
        #     raise ValueError(f"Insufficient balance of {token} ({current_balance}) to decrease by {amount}.")
        self.balances[token] = current_balance - amount

    # TODO: Add event handling methods (deliver, process_inbox) later
    # def deliver(self, event: Event): ...
    # def process_inbox(self): ...
    # def _handle_..._event(self, event): ...

    def __str__(self) -> str:
        balances_str = ", ".join(f"{token}: {bal:.2f}" for token, bal in self.balances.items())
        return f"Node(ID: {self.node_id}, Pos: {self.position}, Balances: [{balances_str}], Imbalance: {self.quantity_imbalance:.4f})"

    def __repr__(self) -> str:
        return f"Node('{self.node_id}')"