# GoThink

Advanced MCP server combining systematic thinking, mental models, and stochastic algorithms for enhanced AI decision-making.

**Author**: W. Alec Akin ([@rainmana](https://github.com/rainmana)) | [Website](https://alecakin.com/about)

> **Acknowledgments**: This project is based on the excellent work from the [Waldzell MCP](https://github.com/waldzellai/waldzell-mcp) repository, specifically the [clear-thought](https://github.com/waldzellai/waldzell-mcp/tree/main/servers) and [stochastic-thinking](https://github.com/waldzellai/waldzell-mcp/tree/main/servers) servers. We've significantly expanded and refactored these components into a comprehensive MCP server with community-driven mental models and advanced thinking frameworks.

## Overview

GoThink is a comprehensive Model Context Protocol (MCP) server written in Go that combines the best of systematic thinking approaches with advanced stochastic algorithms. It provides AI assistants with powerful tools for:

- **Systematic Thinking**: Sequential reasoning, mental models, debugging approaches
- **Stochastic Algorithms**: Markov Decision Processes, Monte Carlo Tree Search, Multi-Armed Bandits, Bayesian Optimization, Hidden Markov Models
- **Decision Frameworks**: Expected utility, multi-criteria analysis, risk assessment
- **Visualization**: Concept maps, decision trees, probability trees, Bayesian networks
- **Hybrid Reasoning**: Adaptive approaches that combine systematic and stochastic methods

## Features

### Systematic Thinking Tools

- **Sequential Thinking**: Structured thought processes with branching and revision support
- **Mental Models**: First principles, opportunity cost, Bayesian thinking, systems thinking
- **Debugging Approaches**: Binary search, reverse engineering, root cause analysis
- **Collaborative Reasoning**: Multi-perspective problem solving
- **Socratic Method**: Question-based inquiry and discovery
- **Creative Thinking**: Ideation and innovation techniques
- **Systems Thinking**: Holistic analysis of complex systems
- **Scientific Method**: Hypothesis-driven investigation

### Stochastic Algorithms

- **Markov Decision Processes (MDPs)**: Policy optimization for sequential decisions
- **Monte Carlo Tree Search (MCTS)**: Strategic planning and game playing
- **Multi-Armed Bandit**: Exploration vs exploitation in decision making
- **Bayesian Optimization**: Efficient parameter optimization
- **Hidden Markov Models (HMMs)**: State inference and pattern recognition
- **Reinforcement Learning**: Adaptive learning from experience

### Decision Frameworks

- **Expected Utility**: Rational decision making under uncertainty
- **Multi-Criteria Analysis**: Weighted evaluation of multiple factors
- **Risk Analysis**: Comprehensive risk assessment and management
- **Stochastic Decision Making**: Probabilistic decision frameworks

### Visualization Tools

- **Concept Maps**: Knowledge representation and organization
- **Mind Maps**: Creative brainstorming and idea mapping
- **Flowcharts**: Process visualization and workflow design
- **Decision Trees**: Decision path visualization
- **Probability Trees**: Probabilistic outcome visualization
- **Bayesian Networks**: Causal relationship modeling

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### Install from Source

```bash
# Clone the repository
git clone https://github.com/rainmana/gothink.git
cd gothink

# Install the MCP server
go install .

# The gothink binary will be installed to $GOPATH/bin
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/rainmana/gothink.git
cd gothink

# Build the application
go build -o gothink .

# Run the server
./gothink
```

### Using Make

```bash
# Build the application
make build

# Run in development mode
make dev

# Run the built application
make run
```

## Configuration

GoThink can be configured via environment variables or a configuration file:

### Environment Variables

```bash
export GOTHINK_PORT=8080
export GOTHINK_HOST=localhost
export GOTHINK_LOG_LEVEL=info
export GOTHINK_ENABLE_STOCHASTIC=true
export GOTHINK_ENABLE_SYSTEMATIC=true
export GOTHINK_ENABLE_VISUALIZATION=true
export GOTHINK_ENABLE_HYBRID=true
```

### Configuration File

Create a `config.json` file:

```json
{
  "port": "8080",
  "host": "localhost",
  "log_level": "info",
  "enable_stochastic_algorithms": true,
  "enable_systematic_thinking": true,
  "enable_visualization": true,
  "enable_hybrid_thinking": true,
  "max_thoughts_per_session": 100,
  "session_timeout": "30m",
  "max_stochastic_iterations": 1000,
  "default_confidence_threshold": 0.8
}
```

## MCP Server Usage

GoThink is an MCP (Model Context Protocol) server that communicates via stdio. It provides AI assistants with powerful thinking tools through the MCP protocol.

### Available Tools

The server exposes the following tools:

#### Thinking Tools
- **sequential_thinking**: Perform structured thought progression
- **mental_model**: Apply mental models to solve problems
- **debugging_approach**: Apply systematic debugging approaches
- **list_mental_models**: List all available mental models

#### Stochastic Algorithms
- **markov_decision_process**: Run MDP optimization for sequential decisions
- **monte_carlo_tree_search**: Run MCTS for game tree exploration
- **multi_armed_bandit**: Run bandit algorithms for exploration vs exploitation

#### Decision Frameworks
- **decision_framework**: Apply decision frameworks for structured decision making

#### Visualization Tools
- **concept_map**: Create and manipulate concept maps for visual thinking

#### Session Management
- **session_stats**: Get statistics for a session
- **session_export**: Export all data for a session

#### Intelligence Tools
- **query_attack**: Query MITRE ATT&CK techniques and tactics
- **query_nvd**: Query NVD CVE data for security vulnerabilities
- **query_owasp**: Query OWASP testing procedures and guidelines
- **refresh_intelligence**: Refresh all intelligence data from external sources
- **intelligence_stats**: Get statistics about available intelligence data

### Testing the MCP Server

You can test the server using JSON-RPC messages:

```bash
# Initialize the server
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | gothink

# List available tools
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' | gothink
```

### Integration with MCP Clients

The server is designed to work with MCP-compatible clients like Claude Desktop. To integrate:

1. Install the server using `go install .`
2. Configure your MCP client to use the `gothink` binary
3. The server will communicate via stdio as required by the MCP protocol

## Mental Models

GoThink includes several built-in mental models:

### First Principles Thinking
Break down complex problems into fundamental components and build up from there.

### Opportunity Cost Analysis
Consider what you give up when making a choice between alternatives.

### Bayesian Thinking
Update beliefs based on new evidence using probabilistic reasoning.

### Systems Thinking
Understand how parts of a system interact and consider emergent properties.

### Error Propagation
Understand how errors compound through complex systems.

### Rubber Duck Debugging
Explain your problem to an inanimate object to gain clarity.

### Pareto Principle
Focus on the 20% of efforts that produce 80% of results.

### Occam's Razor
Prefer simpler explanations when multiple explanations are available.

## Algorithm Selection Guide

### When to Use Systematic Thinking
- **Problem Understanding**: Initial analysis and decomposition
- **Strategic Planning**: Long-term decision making
- **Debugging**: Troubleshooting complex issues
- **Learning**: Understanding new concepts

### When to Use Stochastic Algorithms
- **Optimization**: Finding optimal solutions in complex spaces
- **Decision Making**: Choosing between multiple options under uncertainty
- **Pattern Recognition**: Identifying hidden patterns in data
- **Game Playing**: Strategic decision making in competitive environments

### When to Use Hybrid Approaches
- **Complex Problems**: Multi-faceted issues requiring both systematic and probabilistic reasoning
- **Uncertain Environments**: Situations with high uncertainty and multiple variables
- **Adaptive Systems**: Problems that require learning and adaptation over time

## Development

### Project Structure

```
gothink/
├── main.go                 # MCP server entry point
├── go.mod                  # Go module definition
├── internal/
│   ├── config/            # Configuration management
│   ├── handlers/          # MCP tool handlers
│   ├── models/            # Mental models loader
│   ├── storage/           # Data storage layer
│   ├── types/             # Type definitions
│   └── intelligence/      # Intelligence data services
├── examples/              # Example mental models
├── docs/                  # Documentation
└── README.md              # This file
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gothink-linux .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gothink.exe .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gothink-macos .

# Or use the Makefile
make build-linux
make build-windows
make build-macos
```

### Docker Support

```bash
# Build Docker image
docker build -t gothink .

# Run container
docker run -p 8080:8080 gothink
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Based on the Model Context Protocol (MCP) by Anthropic
- Built upon the [clear-thought](https://github.com/waldzellai/waldzell-mcp/tree/main/servers) and [stochastic-thinking](https://github.com/waldzellai/waldzell-mcp/tree/main/servers) servers from [Waldzell MCP](https://github.com/waldzellai/waldzell-mcp)
- Inspired by classic works in decision theory and cognitive science
- Combines insights from systematic thinking and stochastic optimization
- Built with Go for performance and reliability

## Roadmap

- [ ] Advanced reinforcement learning algorithms
- [ ] Real-time collaboration features
- [ ] Machine learning model integration
- [ ] Advanced visualization capabilities
- [ ] Performance optimization
- [ ] Comprehensive test coverage
- [ ] Documentation improvements
- [ ] Plugin architecture
