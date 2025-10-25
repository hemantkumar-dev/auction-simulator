# Auction Simulator

A concurrent Go-based auction simulator that models multiple bidders participating in simultaneous auctions with resource constraints.

## Project Structure
```
auction-simulator/
├── cmd/
│   └── simulator/      # Main application entrypoint
│       └── main.go
├── internal/
│   ├── auction/       # Auction execution logic
│   │   └── auction.go
│   ├── bidder/        # Bidder behavior & types
│   │   └── bidder.go
│   ├── model/         # Shared data types
│   │   └── types.go
│   ├── resources/     # Resource management
│   │   └── limiter.go
│   └── util/          # Helper functions
│       └── attributes.go
└── tests/             # Test files
    └── auction_test.go
```

## Quick Start

1. Clone the repository
```bash
git clone https://github.com/hemantkumar-dev/auction-simulator.git
cd auction-simulator
```

2. Initialize and update dependencies
```bash
go mod tidy
```

3. Run the simulator
```bash
go run ./cmd/simulator
```

## Configuration Flags

The simulator accepts various command-line flags:

- `-bidders int`: Number of bidders (default: 100)
- `-auctions int`: Number of concurrent auctions (default: 40)
- `-timeout-ms int`: Auction timeout in milliseconds (default: 800)
- `-vcpu int`: vCPU count for resource standardization (default: system CPU count)
- `-ram-mb int`: RAM in MB for resource standardization (default: 2048)
- `-out string`: Output directory for auction results (default: "sample-outputs")
- `-seed int64`: Random seed (default: current time)

Example with custom settings:
```bash
go run ./cmd/simulator -bidders 50 -auctions 20 -timeout-ms 500 -vcpu 4 -ram-mb 1024
```

## Output Files

The simulator generates:
- Individual JSON files for each auction in the output directory (`auction_NNN.json`)
- A summary file (`summary.json`) with overall statistics

