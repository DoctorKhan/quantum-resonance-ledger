import os
import pytest
import numpy as np
import hashlib
import networkx as nx

from python.src.quantum_blockchain import QuantumBlockchain, Bridge

@pytest.fixture(scope="class")
def test_setup(request):
    os.environ['ENV'] = 'test'
    initial_nodes = ["A", "B", "C", "D", "E", "F", "G"]
    chain = QuantumBlockchain(initial_nodes=initial_nodes)
    bridge = Bridge()
    np.random.seed(0)

    chain.network.add_edges_from([
        ("A", "B", {'latency': 0.1, 'fee': 0.01}), ("B", "C", {'latency': 0.2, 'fee': 0.02}),
        ("C", "D", {'latency': 0.15, 'fee': 0.015}), ("D", "E", {'latency': 0.25, 'fee': 0.025}),
        ("A", "E", {'latency': 0.5, 'fee': 0.005}), ("B", "D", {'latency': 0.3, 'fee': 0.01}),
        ("A", "C", {'latency': 0.4, 'fee': 0.03}), ("C", "F", {'latency': 0.1, 'fee': 0.01}),
        ("F", "G", {'latency': 0.2, 'fee': 0.02}), ("G", "E", {'latency': 0.3, 'fee': 0.015}),
        ("A", "G", {'latency': 0.6, 'fee': 0.008})
    ])

    request.cls.chain = chain
    request.cls.bridge = bridge

@pytest.mark.usefixtures("test_setup")
class TestQuantumBlockchain:
    def test_initialization(self):
        """Ensures proper initialization within bounds."""
        for param in ['block_size', 'fee_rate']:
            for node in self.chain.network.nodes():
                value = self.chain.params[param]['value'][node]
                assert self.chain.params[param]['min'] <= value <= self.chain.params[param]['max']

    def test_parameter_bounds_dynamic_updates(self):
        """Verifies parameters stay within hard bounds."""
        for _ in range(10):
            self.chain.update_parameters()
            for param in self.chain.params.values():
                for node in self.chain.network.nodes():
                    assert param['min'] <= param['value'][node] <= param['max']

    def test_transaction_creation_cut_structure(self):
        """Validates transaction structure and CUT commitments."""
        self.chain.create_transaction("A", "B", 1.0)
        tx = self.chain.pending_transactions[0]
        assert 'partial_reveal' in tx and len(tx['partial_reveal']) == 16
        assert 'commitment' in tx and len(tx['commitment']) == 64
        assert tx['commitment'] == hashlib.sha256(b'\xaa' * 32).hexdigest()

    def test_transaction_processing(self):
        """Tests transaction processing with probabilistic errors."""
        initial_A, initial_B = self.chain.token_balances["A"], self.chain.token_balances["B"]
        self.chain.create_transaction("A", "B", 10.0)
        self.chain.update_parameters()
        # Check that the balances have been updated with some probabilistic error
        assert abs(self.chain.token_balances["A"] - (initial_A - 10.0)) < 1.5
        assert abs(self.chain.token_balances["B"] - (initial_B + 10.0)) < 1.5
        assert -0.1 < self.chain.quantity_imbalance["B"] < 0.1

    def test_path_selection_optimal_routing(self):
        """Verifies routing selects a valid path."""
        path = self.chain.path_integral_routing("A", "G")
        valid_paths = [["A", "B", "C", "F", "G"], ["A", "G"], ["A", "E", "G"]]
        assert path in valid_paths

    def test_uncertainty_relation_block_validation(self):
        """Checks uncertainty-based block validation."""
        valid_block = {'params': {'block_size': 1.3, 'fee_rate': 0.017}}
        invalid_block = {'params': {'block_size': 1.001, 'fee_rate': 0.011}}
        assert self.chain.validate_block(valid_block)
        assert not self.chain.validate_block(invalid_block)

    def test_cross_chain_bridge_atomic_swap(self):
        """Tests atomic swap functionality."""
        token_a, token_b = self.bridge.create_entangled_pair()
        status, secret = self.bridge.atomic_swap(token_a)
        assert status and secret == token_b
        status, _ = self.bridge.atomic_swap(token_a)
        assert not status

    def test_network_topology_management(self):
        """Verifies network structure."""
        assert len(self.chain.network.nodes()) == 7
        assert len(self.chain.network.edges()) == 11
        assert self.chain.network["A"]["B"]['latency'] == 0.1

    def test_hamiltonian_adaptive_optimization(self):
        """Tests optimization under stress."""
        for node in self.chain.network.nodes():
            self.chain.params['block_size']['value'][node] = 2.5
        for _ in range(3):
            self.chain.update_parameters()
            self.chain.adapt_weights()
        assert all(self.chain.params['block_size']['value'][node] <= 2.0 for node in self.chain.network.nodes())

    def test_edge_case_no_path_routing_failure(self):
        """Ensures routing fails gracefully."""
        self.chain.network.remove_edges_from(list(self.chain.network.in_edges("G")))
        with pytest.raises(nx.NetworkXNoPath):
            self.chain.path_integral_routing("A", "G")

    def test_laplacian_smoothing_parameter_coherence(self):
        """Verifies smoothing reduces parameter differences."""
        self.chain.params['fee_rate']['value']['B'] = 0.02
        self.chain.params['fee_rate']['value']['C'] = 0.08
        initial_diff = abs(self.chain.params['fee_rate']['value']['B'] - self.chain.params['fee_rate']['value']['C'])
        for _ in range(3):
            self.chain.update_parameters()
        final_diff = abs(self.chain.params['fee_rate']['value']['B'] - self.chain.params['fee_rate']['value']['C'])
        assert final_diff < initial_diff

    def test_scalability_high_transaction_load(self):
        """Tests performance under load."""
        initial_time = self.chain.performance_history[-1]['confirmation_time'] if self.chain.performance_history else 5.0
        for _ in range(10):
            self.chain.create_transaction("A", "E", 1.0)
            self.chain.update_parameters()
        final_time = np.mean([p['confirmation_time'] for p in self.chain.performance_history[-5:]])
        assert final_time < initial_time * 1.5

    def test_confirmation_time_dependency(self):
        """Verifies confirmation time varies with block size."""
        self.chain.params['block_size']['value'] = {node: 1.5 for node in self.chain.network.nodes()}
        self.chain.update_parameters()
        time_high_bs = self.chain.performance_history[-1]['confirmation_time']
        self.chain.params['block_size']['value'] = {node: 0.5 for node in self.chain.network.nodes()}
        self.chain.update_parameters()
        time_low_bs = self.chain.performance_history[-1]['confirmation_time']
        # Higher block size should result in higher confirmation time
        assert time_high_bs > time_low_bs