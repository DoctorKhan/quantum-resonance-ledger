# Test Driven Development Plan for QRL Go Node Implementation

This document outlines the Test-Driven Development (TDD) strategy for developing the Quantum Resonance Ledger (QRL) **Go node implementation**. Unlike the Python simulation framework, this Go code represents the **actual, operational node software**. The focus is on ensuring the correctness and reliability of its core physics-inspired mechanisms and native functionalities through rigorous testing.

## Goal

To ensure the reliability, correctness, and robustness of the **operational Go QRL node** by writing tests before implementing the actual functionality, adhering to TDD best practices.

## Scope

This plan covers the Go packages responsible for the core node logic, primarily within:

-   `go/internal/simulation` (Note: While named "simulation" for historical reasons, this package contains core node logic, not just simulation components).
-   Bridging, WSI, voting, and verification primitives as defined in the **multi-function implementation plan**.

## TDD Principles

We will adhere to the following TDD principles:

*   **Red-Green-Refactor:** Write a failing test *first* defining the desired behavior, then write the minimal code to make it pass, then refactor for clarity and efficiency.
*   **Small, Focused Tests:** Each test should verify a single, specific aspect of the code's behavior.
*   **Isolation:** Tests should be independent, minimizing dependencies on external systems or other tests where possible. Mocks and stubs will be used as needed.
*   **Repeatability:** Tests must produce the same deterministic results every time they are run in the same environment.
*   **Fast Execution:** Tests should run quickly to provide rapid feedback during development.
*   **Comprehensive Coverage:** Aim for high code coverage, especially for critical components like consensus, cryptography, and state transitions (>90% target).

## Detailed Plan

*(Note: Test names are illustrative and may evolve. The focus is on the functionality being tested.)*

### Phase 1: Foundational Components

*(Focus: Basic building blocks, parameters, events, cryptography)*

**1.1 Node & Network Primitives:**
    *   `TestNodeInitialization`: Verify correct initialization of node state (ID, initial parameters, zero balances).
    *   `TestNetworkRepresentation`: (If applicable) Test basic network state management within a node (peer connections, etc.).
    *   `TestLatencyCalculation`: Test calculation of communication latency based on node properties or network conditions (relevant for message passing).

**1.2 Event System (if applicable for internal node logic):**
    *   `TestEventCreation`: Test creating different internal event types (e.g., ParameterUpdateTrigger, MessageReceived).
    *   `TestEventHandling`: Verify correct handler execution for internal events.

**1.3 Parameter Management & Distributions:**
    *   `TestParameterInitialization`: Verify correct setup of node parameters with names, distributions (TruncatedGaussian, etc.), and bounds.
    *   `TestDistributionSampling`: **(Current)** Verify `Sample()` returns values within bounds for `TruncatedGaussian`. *[Test partially exists]*
    *   `TestDistributionMean`: **(Next)** Verify `Mean()` method returns the correct configured mean for `TruncatedGaussian`. *[Test added, needs implementation]*
    *   `TestDistributionStdDev`: **(Next)** Verify `StdDev()` method returns the correct configured standard deviation for `TruncatedGaussian`. *[Test added, needs implementation]*
    *   `TestParameterUpdateRule`: Test the core parameter update logic (e.g., Langevin dynamics if used, influence of Hamiltonian gradients, smoothing). Ensure updates respect bounds.
    *   `TestUncertaintyRelationHandling`: Test logic related to defined uncertainty relations between parameters.

**1.4 Cryptographic Uniqueness Tokens (CUTs):**
    *   `TestCUT_Creation`: Test generation of secret key `sk` and commitment `C`.
    *   `TestCUT_CommitmentVerification`: Verify `C` matches `sk`.
    *   `TestCUT_SpendProofGeneration`: Test generating valid ZKPs for spending CUTs.
    *   `TestCUT_SpendProofVerification`: Test verification of valid ZKPs and rejection of invalid ones.
    *   `TestCUT_NoCloning`: Test mechanisms preventing duplication/reuse.
    *   `TestCUT_Representation`: Test correct representation of QRG, Gas, and bridged assets like `qUSDC`.

**1.5 Basic Transaction Handling:**
    *   `TestTransactionSerialization`: Test encoding/decoding of base transaction structures.
    *   `TestSignatureVerification`: Test signing and verifying transactions using node keys.
    *   `TestTransactionValidation_Basic`: Test initial validation rules (format, signature, nonce).

### Phase 2: Core QRL Dynamics & Hamiltonian

*(Focus: Implementing the core physics-inspired mechanics within the node)*

**2.1 Relativistic Latency Effects (if modeled directly in node):**
    *   `TestGammaCalculation`: Test Lorentz factor `γ` calculation if used.
    *   `TestObservedTimeCalculation`: Test time dilation effects on message timestamps or event processing if applicable.

**2.2 Path Integral Core Logic (as applied to node decision-making):**
    *   `TestActionCalculation`: Test calculating the action `S` for potential state transitions or operations, including relevant cost terms (latency, fees, security, validity, uncertainty penalties, native function costs like WSI peg deviation).
    *   `TestAmplitudeCalculation`: Test calculating transition amplitudes from actions.
    *   `TestProbabilityCalculation`: Test calculating probabilities from amplitudes.

**2.3 Monte Carlo Decision Making (if used for path selection/consensus):**
    *   `TestMonteCarloSampling`: Test the core MC sampling algorithm (e.g., Metropolis-Hastings) used for decision-making.
    *   `TestAcceptanceCriteria`: Test the acceptance logic within the MC sampler.

**2.4 Extended Hamiltonian Calculation:**
    *   `TestHamiltonianComponentCalculation`: Test calculating individual cost terms relevant to the node's operation (network costs, WSI stability penalty, voting costs, bridging costs, verification costs).
    *   `TestHamiltonianTotalCalculation`: Test combining component costs with weights (`w_f`) and penalties.
    *   `TestHamiltonianGradient`: Test calculating gradients `∇H` with respect to node parameters (`θ_i`).
    *   `TestPenalty_WSI_Peg`: Test calculation of the penalty for WSI peg deviation.

**2.5 Advanced Bridging Mechanisms:**
    *   `TestNettingFlowCalculation`: Test the netting flow optimization algorithm for bridging intents.
    *   `TestProbabilisticRelease`: Test the probabilistic release mechanism for bridged assets.
    *   `TestInventoryManagement`: Test inventory pool updates and consistency during bridging operations.
    *   `TestBridgeSecurity`: Test mechanisms preventing double-spends and ensuring consistency during asynchronous bridging.

**2.6 Verification Primitives:**
    *   `TestVerificationDataStructures`: Test the data structures used for anchoring cryptographic proofs.
    *   `TestVerificationTransactionValidation`: Test validation rules for anchoring transactions.
    *   `TestVerificationThroughput`: Test the system's ability to handle high verification loads.

### Phase 3: Native Function Implementation

*(Focus: Testing the logic of each native function as implemented within the node)*

**3.1 Wavefunction Stability Index (WSI) Mechanism:**
    *   `TestWSI_ParameterHandling`: Test setting/updating WSI target weights $\theta_{w,i}$.
    *   `TestWSI_OracleProcessing`: Test consuming and validating oracle price feeds.
    *   `TestWSI_ValueCalculation`: Test local calculation of the WSI value.
    *   `TestWSI_PegDeviationPenalty`: Test calculating the local `Penalty_WSI_Peg`.
    *   `TestWSI_WeightUpdateTrigger`: Test triggering updates to $\theta_{w,i}$ based on Hamiltonian gradients.

**3.2 Advanced Bridging Module:**
    *   `TestBridgeIntentHandling`: Test handling of bridging intents and their aggregation.
    *   `TestBridgeNettingOptimization`: Test the optimization of netting flows over epochs.
    *   `TestBridgeConfirmationDepth`: Test dynamic adjustment of confirmation depth (`θ_confirmation_depth`).
    *   `TestBridgeRiskAssessment`: Test risk assessment logic using the Hamiltonian and `Q` field dynamics.

**3.3 Verification Primitives:**
    *   `TestVerificationMerkleTree`: Test the integrity of Merkle trees used for proof anchoring.
    *   `TestVerificationCommitmentScheme`: Test the correctness of commitment schemes for off-chain data.

**3.4 Quantum Stable Dollar (QSD) Mechanism:**
    *   `TestQSD_Minting`: Verify correct QSD minting against locked collateral (CUTs), checking collateral ratio (`θ_collateral_ratio`).
    *   `TestQSD_Burning`: Verify correct QSD burning and collateral redemption.
    *   `TestQSD_CollateralRatioCheck`: Test accurate calculation and enforcement of the minimum collateral ratio.
    *   `TestQSD_StabilityFee`: Test calculation and application of the dynamic stability fee (`θ_stability_fee`).
    *   `TestQSD_LiquidationTrigger`: Test correct identification of positions eligible for liquidation based on collateral value and `θ_collateral_ratio`.
    *   `TestQSD_LiquidationProcess`: Test the basic mechanics of a liquidation (marking for auction, applying penalty `θ_liquidation_penalty`). (Auction participation might be Phase 4).
    *   `TestQSD_WSIReference`: Test how QSD logic potentially uses the WSI value (e.g., in Hamiltonian term or risk assessment).
    *   `TestQSD_TransactionValidation`: Test specific validation rules for QSD mint/burn/manage transactions.
    *   `TestQSD_StateUpdate`: Verify correct updates to node state regarding QSD vaults, debt, and collateral.
    *   `TestQSD_QRGInteraction_FeeSink`: Test mechanism for using stability fees to potentially interact with QRG (e.g., marking fees for buy/burn).


### Phase 4: Integration, Consensus & System-Level

*(Focus: Interactions between components, consensus participation, state management)*

**4.1 Integrated Transaction Processing:**
    *   `TestTxValidation_AllTypes`: Ensure node validation logic correctly handles all transaction types (WSI updates, votes, bridge ops, anchors, transfers).
    *   `TestStateUpdate_AllTypes`: Verify processing transactions correctly updates the node's local state (balances, vote tallies, bridge state, anchor data, `Q` field).

**4.2 Consensus Participation:**
    *   `TestConsensusMessageHandling`: Test processing consensus messages from peers.
    *   `TestForkResolutionLogic`: Test the node's logic for contributing to and resolving forks based on path integral probabilities.
    *   `TestConsensus_DoubleSpendingPrevention`: Verify node logic prevents inclusion/validation of double spends.
    *   `TestConsensus_ParameterConvergence`: Test participation in achieving consensus on dynamic parameters (`Θ`), including WSI weights and bridging parameters.

**4.3 State Management:**
    *   `TestStateConsistency`: Verify internal consistency of the node's state components.
    *   `TestStateStorageRetrieval`: Test efficient storage and retrieval of node state.
    *   `TestStateUpdate_Bridging`: Verify state updates for bridging operations, including inventory and netting flows.
    *   `TestStateUpdate_Verification`: Verify state updates for anchored verification proofs.

## Implementation Notes

-   Start with the tests for `distribution.go` as planned.
-   Proceed through the phases, writing failing tests first for each piece of functionality.
-   Use Go's built-in `testing` package.
-   Employ mocks/stubs for external dependencies (like network peers or oracle feeds) during unit testing.
-   Integration tests will be crucial in later phases to test interactions between components.
-   Continuously run tests using `go test ./...` or targeted package tests.