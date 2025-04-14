# Installation, Setup, and Running

## Prerequisites

- Go 1.22 or later
- Git

## Installation

1.  Clone the repository:

    ```bash
    git clone git@github.com:DoctorKhan/quantum-resonance-ledger.git
    cd quantum-resonance-ledger
    ```

2.  Navigate to the `node` directory:

    ```bash
    cd node
    ```

3.  Install the dependencies:

    ```bash
    go mod download
    ```

## Setup

1.  Configure the environment variables:

    -   Create a `.env` file in the `node` directory.
    -   Add the following environment variables:

        ```
        # Example environment variables
        # Replace with your actual values
        API_KEY=your_api_key
        DATABASE_URL=your_database_url
        ```

2.  Configure the node:

    -   Create a `config.yaml` file in the `node/cmd/qrlnode` directory.
    -   Add the following configuration options:

        ```yaml
        # Example configuration options
        # Replace with your actual values
        network: mainnet
        db_path: /path/to/your/database
        ```

## Running

1.  Run the node:

    ```bash
    go run cmd/qrlnode/main.go --config cmd/qrlnode/config.yaml
    ```

2.  Run the tests:

    ```bash
    go test ./...