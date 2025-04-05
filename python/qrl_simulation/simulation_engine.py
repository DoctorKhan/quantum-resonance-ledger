import time
import random # Added for transaction generation
from typing import Dict, List # Added List
import matplotlib.pyplot as plt # Added for plotting

from .quantum_blockchain import QuantumBlockchain

class SimulationEngine:
    """Manages the execution of the QRL simulation."""

    def __init__(self, blockchain: QuantumBlockchain = None, config: Dict = None):
        """
        Initializes the simulation engine.

        Args:
            blockchain: An instance of the QuantumBlockchain. If None, a default instance with 4 nodes is created.
            config: A dictionary containing simulation configuration,
                    e.g., {'steps': 100}. If None, a default config with 10 steps is used.
        """
        # Create default blockchain if none provided
        if blockchain is None:
            initial_nodes = ["Node1", "Node2", "Node3", "Node4"]
            blockchain = QuantumBlockchain(initial_nodes=initial_nodes)
            # Add some basic network connections
            blockchain.network.add_edges_from([("Node1", "Node2"), ("Node2", "Node3"), ("Node1", "Node3"), ("Node3", "Node4")])
        
        # Create default config if none provided
        if config is None:
            config = {'steps': 10}
        if not isinstance(blockchain, QuantumBlockchain):
            raise TypeError("blockchain must be an instance of QuantumBlockchain")
        if not isinstance(config, dict) or 'steps' not in config:
            raise ValueError("config must be a dict with at least a 'steps' key")

        self.blockchain = blockchain
        self.total_steps = int(config['steps'])
        self.current_step = 0
        self.metrics: Dict[str, List] = {
            "step": [],
            "transactions_processed": [],
            "total_token_balance": []
        }
        # Add more config parameters as needed (e.g., transaction rate)

    def run(self):
        """Runs the simulation for the configured number of steps."""
        print(f"Starting simulation for {self.total_steps} steps...")
        start_time = time.time()

        for step in range(self.total_steps):
            self.current_step = step + 1
            print(f"--- Step {self.current_step}/{self.total_steps} ---")

            # 1. Update blockchain state (e.g., parameters)
            # In a more complex simulation, this would involve event processing
            self.blockchain.update_parameters()
            print("Updated blockchain parameters.")

            # 2. Generate events (e.g., new transactions)
            self._generate_transactions()
            # print("Generated transactions.") # Optional print

            # 3. Process events/transactions
            processed_count = self.blockchain.process_pending_transactions()
            # print("Processed transactions.") # Optional print

            # 4. Collect metrics
            self._collect_metrics(processed_count)
            # print("Collected metrics.") # Optional print

            # Optional: Add a small delay or check for termination conditions
            # time.sleep(0.1)

        end_time = time.time()
        print(f"\nSimulation finished in {end_time - start_time:.2f} seconds.")
        self._print_metrics_summary()
        self._plot_metrics() # Call the new plotting method


    def _print_metrics_summary(self):
        """Prints a summary of the collected metrics."""
        print("\n--- Simulation Metrics Summary ---")
        if not self.metrics["step"]:
            print("No metrics collected.")
            return

        total_processed = sum(self.metrics["transactions_processed"])
        avg_processed = total_processed / len(self.metrics["step"])
        final_balance = self.metrics["total_token_balance"][-1]

        print(f"Total Steps: {self.metrics['step'][-1]}")
        print(f"Total Transactions Processed: {total_processed}")
        print(f"Average Transactions Processed per Step: {avg_processed:.2f}")
        print(f"Final Total Token Balance: {final_balance:.2f}")
        # Add more summary stats as needed
        print("---------------------------------")


    def _plot_metrics(self):
        """Plots the collected simulation metrics."""
        print("\nGenerating plots...")
        if not self.metrics["step"]:
            print("No metrics to plot.")
            return

        steps = self.metrics["step"]
        tx_processed = self.metrics["transactions_processed"]
        total_balance = self.metrics["total_token_balance"]

        fig, axs = plt.subplots(2, 1, figsize=(10, 8), sharex=True)
        fig.suptitle('Simulation Metrics Over Time')

        # Plot Transactions Processed
        axs[0].plot(steps, tx_processed, marker='o', linestyle='-', color='b')
        axs[0].set_ylabel('Transactions Processed')
        axs[0].set_title('Transactions Processed per Step')
        axs[0].grid(True)

        # Plot Total Token Balance
        axs[1].plot(steps, total_balance, marker='x', linestyle='--', color='r')
        axs[1].set_xlabel('Simulation Step')
        axs[1].set_ylabel('Total Token Balance')
        axs[1].set_title('Total Token Balance Across All Nodes')
        axs[1].grid(True)

        plt.tight_layout(rect=[0, 0.03, 1, 0.95]) # Adjust layout to prevent title overlap
        # Save the plot to a file instead of showing interactively
        plot_filename = "simulation_metrics.png"
        plt.savefig(plot_filename)
        print(f"Plots saved to {plot_filename}")
        # plt.show() # Avoid showing interactive plot in this context


    # --- Placeholder methods for future expansion ---

    def _generate_transactions(self):
        """Generates one new random transaction per step."""
        node_ids = list(self.blockchain.nodes.keys())
        if len(node_ids) < 2:
            # print("Not enough nodes to create a transaction.")
            return # Need at least two nodes

        # Ensure sender and receiver are different
        sender_id = random.choice(node_ids)
        receiver_id = random.choice(node_ids)
        while receiver_id == sender_id:
            receiver_id = random.choice(node_ids)

        # Generate a random amount (e.g., between 1.0 and 10.0)
        amount = round(random.uniform(1.0, 10.0), 2)

        print(f"Generating transaction: {sender_id} -> {receiver_id} ({amount})")
        self.blockchain.create_transaction(sender_id, receiver_id, amount)


    def _collect_metrics(self, processed_count: int):
        """Collects basic performance and state metrics during the simulation."""
        # processed_count is now passed in from the run loop
        total_balance = sum(self.blockchain.token_balances.values())

        self.metrics["step"].append(self.current_step)
        self.metrics["transactions_processed"].append(processed_count) # Use the passed-in count
        self.metrics["total_token_balance"].append(total_balance)
        # print(f"Metrics collected: Step {self.current_step}, Processed: {processed_count}, Total Balance: {total_balance:.2f}")


if __name__ == '__main__':
    # Example usage (requires running from the project root directory typically)
    print("Running basic simulation engine example...")

    # Basic setup
    initial_nodes = ["Node1", "Node2", "Node3", "Node4"]
    qrl_chain = QuantumBlockchain(initial_nodes=initial_nodes)

    # Add some basic network connections for context (optional for this basic run)
    qrl_chain.network.add_edges_from([("Node1", "Node2"), ("Node2", "Node3"), ("Node1", "Node3"), ("Node3", "Node4")])

    simulation_config = {'steps': 10}
    engine = SimulationEngine(blockchain=qrl_chain, config=simulation_config)

    engine.run()

    print("\nExample finished.")
    # You could print some final state from qrl_chain here if needed
    # print(f"Final balances: {qrl_chain.token_balances}")