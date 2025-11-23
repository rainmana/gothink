# GoThink

Advanced MCP server combining systematic thinking, mental models, and stochastic algorithms for enhanced AI decision-making.

**Author**: W. Alec Akin ([@rainmana](https://github.com/rainmana)) | [Website](https://alecakin.com/about)

> **Acknowledgments**: This project is based on the excellent work from the [Waldzell MCP](https://github.com/waldzellai/waldzell-mcp) repository, specifically the [clear-thought](https://github.com/waldzellai/waldzell-mcp/tree/main/servers) and [stochastic-thinking](https://github.com/waldzellai/waldzell-mcp/tree/main/servers) servers. We've significantly expanded and refactored these components into a comprehensive MCP server with community-driven mental models and advanced thinking frameworks.

## Overview

GoThink is a comprehensive Model Context Protocol (MCP) server written in Go that focuses on systematic thinking approaches. It provides AI assistants with powerful tools for:

- **Systematic Thinking**: Sequential reasoning, mental models, debugging approaches
- **Collaborative Reasoning**: Multi-perspective problem solving

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [MCP Server Usage](#mcp-server-usage)
- [Mental Models](#mental-models)
- [Algorithm Selection Guide](#algorithm-selection-guide)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)


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



## Installation

### Prerequisites

- Go 1.21 or later
- Git

### Install from Source

```bash
# Install the latest version
go install github.com/rainmana/gothink/cmd/gothink@latest
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
export GOTHINK_MENTAL_MODELS_PATH=/path/to/models
```

### Configuration File

Create a `config.json` file:

```json
{
  "port": "8080",
  "host": "localhost",
  "log_level": "info",
  "max_thoughts_per_session": 100,
  "session_timeout": "30m",
  "mental_models_path": "/path/to/models"
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

#### Session Management
- **session_stats**: Get statistics for a session
- **session_export**: Export all data for a session


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



## Development

### Project Structure

```
gothink/
├── cmd/
│   └── gothink/
│       └── main.go        # MCP server entry point
├── go.mod                  # Go module definition
├── internal/
│   ├── config/            # Configuration management
│   ├── handlers/          # MCP tool handlers
│   ├── models/            # Mental models loader
│   ├── storage/           # Data storage layer
│   ├── types/             # Type definitions
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
GOOS=linux GOARCH=amd64 go build -o gothink-linux ./cmd/gothink

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gothink.exe ./cmd/gothink

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gothink-macos ./cmd/gothink

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
- Combines insights from systematic thinking and cognitive science
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
