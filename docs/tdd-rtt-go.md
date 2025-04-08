# TDD Implementation Plan for Reflexive Resonance Trading (RRT) in Go Node

This plan outlines the steps to implement the RRT protocol as described in `docs/rtt.md` using a Test-Driven Development approach within the `node/internal/core` package of the Go implementation.

**Guiding Principles:**

*   **Test First:** Write failing tests before writing implementation code.
*   **Small Steps:** Implement the minimum code required to pass the current test.
*   **Refactor:** Improve code structure after tests pass.
*   **Focus on Behavior:** Tests should verify the *behavior* described in the RRT specification.
*   **Modularity:** Design components (fields, settlement, propagation) with clear interfaces for testability.

**Proposed File Structure:**

*   Create a new package: `node/internal/core/rtt/`
*   Core RRT logic files: `node/internal/core/rtt/field.go`, `node/internal/core/rtt/settlement.go`, `node/internal/core/rtt/propagation.go`, `node/internal/core/rtt/types.go`, `node/internal/core/rtt/hamiltonian.go`
*   Corresponding test files: `node/internal/core/rtt/field_test.go`, `node/internal/core/rtt/settlement_test.go`, `node/internal/core/rtt/propagation_test.go`, `node/internal/core/rtt/hamiltonian_test.go`
*   Integrate RRT state into `node/internal/core/state.go`.
*   Integrate RRT processing into `node/internal/core/node.go` or a dedicated RRT processor called by the node.

**TDD Steps:**

1.  **Define Core Data Structures (in `rtt/types.go`):**
    *   **Test:** Write tests to verify the creation, initialization, and basic properties of structures representing:
        *   `PropensityField` (`Ψ_buy`, `Ψ_sell`): Needs to represent probability density over asset/price ranges. Consider using histograms, splines, or other suitable representations. Test serialization/deserialization.
        *   `LocalNodeStateRTT`: Contains local fields, relevant CUTs, recent neighbor data snapshots.
        *   `SettlementRecord`: Details of a probabilistic local settlement (amount, price range, involved CUTs).
        *   `PropagationPacket`: Data exchanged between nodes (field updates, settlement info).
    *   **Implement:** Create the basic struct definitions.

2.  **Trading Intent Perturbation (in `rtt/field.go`):**
    *   **Test:** Write tests for a function `PerturbPropensityField(field *PropensityField, asset, priceRange, magnitude)`:
        *   Verify that perturbing increases probability density correctly within the specified range.
        *   Test accumulation of multiple perturbations.
        *   Test edge cases (zero/negative magnitude, invalid ranges).
        *   Test that perturbations are reflected immediately in the local state structure.
    *   **Implement:** Implement the `PerturbPropensityField` function and update the local state.

3.  **Local Probabilistic Reflection & Settlement (in `rtt/settlement.go`):**
    *   **Test:** Write tests for `CalculateOverlap(buyField, sellField *PropensityField)`:
        *   Mock fields with known overlaps and verify the calculated overlap integral/metric.
    *   **Test:** Write tests for `AttemptLocalSettlement(state *LocalNodeStateRTT)`:
        *   Mock `LocalNodeStateRTT` with varying degrees of buy/sell field overlap.
        *   Simulate multiple runs and statistically verify that `SettlementRecord` is generated with probability proportional to the calculated overlap (above a `θ_trade_threshold`).
        *   Verify the generated `SettlementRecord` contains correct details (amount proportional to overlap, price range).
        *   Verify that local balances (within the test state) are updated correctly upon settlement.
    *   **Implement:** Implement `CalculateOverlap` and `AttemptLocalSettlement`.

4.  **CUTs Management Integration (in `rtt/settlement.go` and potentially `core/cut.go`):**
    *   **Test:** Extend `AttemptLocalSettlement` tests:
        *   Verify that CUTs specified in the `LocalNodeStateRTT` are marked or consumed upon successful settlement generation.
        *   Write tests ensuring that attempting settlement with already-marked/consumed CUTs fails or is ignored.
    *   **Implement:** Add logic to `AttemptLocalSettlement` to interact with CUT representations (potentially requiring interfaces to mock `core/cut` behavior if direct modification is complex).

5.  **Propagation Logic (in `rtt/propagation.go`):**
    *   **Test:** Write tests for `GeneratePropagationPacket(state *LocalNodeStateRTT)`:
        *   Verify packet contains necessary field updates and recent settlement records.
    *   **Test:** Write tests for `ApplyPropagationPacket(localState *LocalNodeStateRTT, packet *PropagationPacket, latency time.Duration)`:
        *   Mock incoming packets from neighbors.
        *   Verify local propensity fields are updated based on neighbor data (implementing Laplacian/D'Alembertian smoothing logic).
        *   Verify received settlement information updates local knowledge (e.g., marking relevant CUTs if seen settled elsewhere).
        *   Test latency effects (information is applied based on simulated arrival time).
    *   **Implement:** Implement packet generation and application logic, including the core field propagation equations.

6.  **Quantity Imbalance (`Q`) Management (integrated in `state.go`, `rtt/settlement.go`, `rtt/propagation.go`):**
    *   **Test:** Modify `AttemptLocalSettlement` tests:
        *   Create scenarios where settlement occurs based on incomplete (local-only) information, verify `Q` imbalance is recorded in the state.
    *   **Test:** Modify `ApplyPropagationPacket` tests:
        *   Verify that the propagation/smoothing logic includes terms that gradually reduce the `Q` imbalance over time as more global information arrives.
        *   Test network-wide convergence: simulate multiple nodes exchanging packets and verify total `Q` trends towards zero.
    *   **Implement:** Add `Q` field to state, update it during settlement, and add correction terms to the propagation logic.

7.  **Hamiltonian Integration (in `rtt/hamiltonian.go`):**
    *   **Test:** Write tests for `CalculateHamiltonian(state *GlobalStateSnapshot)` (requires defining what constitutes a relevant snapshot):
        *   Mock states with high `Q` imbalance, high field divergence between nodes, high field volatility. Verify the Hamiltonian calculation reflects these penalties.
    *   **Test:** Write tests for `OptimizeParameters(hamiltonianValue, currentState)`:
        *   Verify that based on the Hamiltonian, system parameters (like `θ_trade_threshold`, field constraints) are adjusted in the direction that should minimize H. (This might involve mocking an optimization algorithm).
    *   **Implement:** Implement Hamiltonian calculation and parameter adjustment logic.

8.  **Integration Tests (in `rtt/rtt_integration_test.go` or `core/node_test.go`):**
    *   **Test:** Set up multi-node simulations:
        *   Scenario: Node A perturbs buy, Node B perturbs sell. Verify local settlements near A and B, propagation, and eventual balance consistency.
        *   Scenario: High volume, concurrent perturbations. Verify system stability and `Q` management.
        *   Scenario: Network latency/partitions. Verify behavior aligns with RRT principles.
        *   Measure metrics: time to first local settlement, time to global state convergence (probabilistic), throughput.
    *   **Implement:** Write test harnesses to orchestrate multi-node interactions and assertions.

---

This plan provides a structured approach to implementing RRT with continuous testing, ensuring each component behaves as expected before integrating it into the larger system.