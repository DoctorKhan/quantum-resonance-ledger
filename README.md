# Quantum Resonance Ledger (QRL)

Production implementation of the quantum-resistant blockchain from the [QRL Whitepaper](docs/qrl_whitepaper.md), featuring:

- **Probabilistic Quantity Conservation** with Laplacian/D'Alembertian operators
- **Cryptographic Uniqueness Tokens (CUTs)** with zero-knowledge proofs
- **Hamiltonian-driven parameter optimization**
- **Path Integral Consensus** with probabilistic finality
- **Relativistic node modeling** with self-healing balance fields

- **Probabilistic Quantity Conservation** with Laplacian/D'Alembertian correction
- **Cryptographic Uniqueness Tokens (CUTs)** preventing double-spending
- **Dynamic parameter management** using uncertainty relations
- **Path Integral Consensus** for probabilistic finality
- **Relativistic network modeling** with self-correcting balance fields

## Project Structure

```
qrl/
├── cmd/                    # Command-line applications in Go
│   ├── example/            # Example application showcasing parameter space concepts
│   └── simulation/         # Simulation application for the full blockchain framework
├── docs/                   # Documentation files in Markdown format
│   └── images/             # Images used within the documentation
├── internal/               # Internal Go packages, not intended for external use
├── scripts/                # Utility scripts for development and testing
├── src/                    # Python source code for the core framework logic
└── test/                   # Test files for both Go and Python components
```
=======
## Project Structure

```
qrl/
├── cmd/                    # Go command-line applications
│   ├── example/            # Example Go application for parameter space
│   └── simulation/         # Example Go application for full simulation
├── docs/                   # Documentation files (Markdown)
│   └── images/             # Images used in documentation
├── internal/               # Go internal packages
├── scripts/                # Utility scripts
├── src/                    # Python source code
└── test/                   # Test files
```

## Test-Driven Development (TDD)

This project rigorously adheres to Test-Driven Development (TDD) principles to ensure code quality, reliability, and maintainability.  TDD is applied throughout the development process, guiding the design and implementation of all features.

**Key TDD Practices:**

1.  **Red-Green-Refactor Cycle**:  The core of our TDD process involves the classic Red-Green-Refactor cycle:
    *   **Red**:  Start by writing a test that fails. This test defines a specific piece of functionality or behavior that needs to be implemented. Running the test should clearly show that it fails because the functionality is not yet present.
    *   **Green**: Write the minimum amount of code necessary to make the test pass. The focus at this stage is solely on satisfying the test requirements, not on perfect design or efficiency.
    *   **Refactor**: Once the test is passing (Green), refactor the code to improve its structure, readability, and efficiency, while ensuring that all tests remain passing. This step helps to clean up the code and prepare it for future extensions without breaking existing functionality.

2.  **Comprehensive Test Coverage**: We aim for high test coverage to ensure that all parts of the system are thoroughly tested. This includes:
    *   **Unit Tests**:  Detailed unit tests for individual components and modules to verify their correct behavior in isolation.
    *   **Integration Tests**: Integration tests to ensure that different components of the system work correctly together and that interactions between them are well-handled.

## Core Components

### Spacetime & Node Model

This component models the fundamental aspects of a distributed network, representing nodes in a simulated spacetime. It accounts for the physical positions of nodes and the relativistic latencies in communication between them.

**Key Features:**

*   **Node Representation**: Nodes are defined by their spatial coordinates (x, y, z) within the simulation space.
*   **Relativistic Latency**: Communication latency between nodes is calculated based on the distance between them, incorporating relativistic effects to simulate realistic network delays.
*   **Network Abstraction**: Provides an abstraction for managing a network of nodes and calculating distances and latencies between any two nodes in the network.

**Code Example:**

```go
// Create a network instance
network := simulation.NewNetwork()

// Instantiate nodes with specific positions in spacetime
node1 := simulation.NewNode("node1", 0.0, 0.0, 0.0) // Node at origin
node2 := simulation.NewNode("node2", 3.0, 0.0, 0.0) // Node at x=3.0

// Incorporate nodes into the network
network.AddNode(node1)
network.AddNode(node2)

// Determine the spatial distance between node1 and node2
distance := network.Distance(0, 1) // Distance between node indices 0 and 1

// Calculate the communication latency between node1 and node2, considering relativistic effects
latency := network.Latency(0, 1)   // Latency between node indices 0 and 1
```
=======
### Spacetime & Node Model

Represents nodes in a network with positions and latencies:

```go
// Create a network
network := simulation.NewNetwork()

// Create nodes at different positions
node1 := simulation.NewNode("node1", 0.0, 0.0, 0.0)
node2 := simulation.NewNode("node2", 3.0, 0.0, 0.0)

// Add nodes to the network
network.AddNode(node1)
network.AddNode(node2)

// Calculate distance and latency
distance := network.Distance(0, 1)
latency := network.Latency(0, 1)
```

### Event System

The Event System is crucial for managing the simulation's progression. It handles events in chronological order, allowing for asynchronous simulation of network activities.

**Key Features:**

*   **Chronological Event Processing**: Events are processed based on their timestamps, ensuring that the simulation unfolds in a realistic, time-ordered manner.
*   **Asynchronous Simulation**: Enables the simulation of concurrent activities and interactions within the network, such as transaction creation and propagation, without requiring strict sequential processing.
*   **Event Queue Management**: Provides an efficient event queue to store, prioritize, and retrieve events for processing.

**Code Example:**

```go
// Initialize a new event queue to manage simulation events
eventQueue := simulation.NewEventQueue()

// Construct a new simulation event, e.g., a transaction creation event
event := simulation.NewEvent(
    simulation.EventTypeTransactionCreated, // Event type: Transaction creation
    10.0,                                 // Timestamp: Time at which the event occurs
    node1,                                // From Node: Node initiating the transaction
    node2,                                // To Node: Node intended to receive the transaction
    payload,                              // Payload: Data associated with the event (e.g., transaction details)
)

// Enqueue the event for processing
eventQueue.AddEvent(event)

// Event processing loop: Continue as long as there are events in the queue
for eventQueue.HasEvents() {
    event := eventQueue.ProcessNextEvent() // Dequeue and retrieve the next event to process
    // Handle the event: Execute actions based on the event type and associated data
    // e.g., process transaction, update node state, etc.
}
```
=======
### Event System

Manages events in chronological order:

```go
// Create an event queue
eventQueue := simulation.NewEventQueue()

// Create events
event := simulation.NewEvent(simulation.EventTypeTransactionCreated, 10.0, node1, node2, payload)

// Add events to the queue
eventQueue.AddEvent(event)

// Process events
for eventQueue.HasEvents() {
    event := eventQueue.ProcessNextEvent()
    // Handle the event
}
```

### Parameter Management

Manages parameters with uncertainty relations:

```go
// Create a parameter manager
paramManager := simulation.NewParameterManager()

// Create parameters
param1 := simulation.NewParameter("param1", 0.0, 10.0)
param2 := simulation.NewParameter("param2", 0.0, 10.0)

// Add parameters to the manager
paramManager.AddParameter(param1)
paramManager.AddParameter(param2)

// Create distributions
dist1 := simulation.NewUniformDistribution(param1)
dist2 := simulation.NewNormalDistribution(param2, 5.0, 1.0)

// Set distributions
paramManager.SetDistribution(param1, dist1)
paramManager.SetDistribution(param2, dist2)

// Create uncertainty relations
relation := simulation.NewUncertaintyRelation(param1, param2, 1.0)

// Add relations to the manager
paramManager.AddUncertaintyRelation(relation)

// Validate uncertainty relations
valid, violations := paramManager.ValidateUncertaintyRelations()
```

### Transaction Processing

Manages transactions with receiver-pays model:

```go
// Create a transaction manager
txManager := simulation.NewTransactionManager()

// Create a transaction
tx := simulation.NewTransaction("sender", "receiver", 10.0, 1.0, 0.5)

// Add the transaction to the manager
txManager.AddTransaction(tx)

// Process the transaction
txManager.ProcessTransaction(tx)
```

## Running the Examples

To run the parameter space example:

```bash
go run cmd/example/main.go
```

To run the full simulation example:

```bash
go run cmd/simulation/main.go
```

## Running Tests

The project includes a custom test formatter that provides Jest-like output for Go tests:

```bash
# Run tests with pretty formatting
./scripts/pretty-test ./pkg/simulation

# Run tests with verbose output and pretty formatting
./scripts/pretty-test -v ./pkg/simulation

# Run a specific test with pretty formatting
./scripts/pretty-test -v ./pkg/simulation -run TestParameterManager

# Run tests with skipped tests
./scripts/pretty-test -v -short ./pkg/simulation
```

The formatter displays test results in a Jest-like format with colored output, making it easier to read and understand test results.

### Watch Mode

The project also includes a watch mode for tests, similar to Jest's watch mode:

```bash
# Run tests in watch mode
./scripts/pretty-test-watch -v ./pkg/simulation
```

This will run the tests and then watch for file changes. When a file is modified, the tests will automatically run again.

## Next Steps

Future development will focus on:

1. **Phase 2: Path Integral & Relativistic Modeling**
   - Relativistic Latency
   - Path Integral Core
   - Integrating Path Integral with Transaction Propagation
   - Receiver-Pays Fee Model

2. **Phase 3: Blockchain Components & Consensus**
   - Advanced Transaction / "No-Cloning" Tokens
   - Blocks & State
   - Path Integral Consensus

3. **Phase 4: Advanced & Optional Features**
   - Simulated "Entanglement"
   - Visualization / Analysis Tools
   - Performance Benchmarking
   - Experimentation & Parameter Tuning

## License

This project is licensed under the MIT License - see the LICENSE file for details.