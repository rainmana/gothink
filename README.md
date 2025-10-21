# GoThink

Advanced MCP server combining systematic thinking, mental models, and stochastic algorithms for enhanced AI decision-making.

**Author**: W. Alec Akin ([@rainmana](https://github.com/rainmana)) | [Website](https://alecakin.com/about)

> **Acknowledgments**: This project was inspired by the excellent work of the original GoThink project. We've significantly expanded and refactored it into a comprehensive MCP server with community-driven mental models and advanced thinking frameworks.

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

### Build from Source

```bash
# Clone the repository
git clone https://github.com/waldzellai/gothink.git
cd gothink

# Build the application
go build -o gothink main.go

# Run the server
./gothink
```

### Using Go Modules

```bash
go mod tidy
go run main.go
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

## API Usage

### Sequential Thinking

```bash
curl -X POST http://localhost:8080/api/v1/thinking/sequential \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "thought": "Let me analyze this problem step by step",
    "thought_number": 1,
    "total_thoughts": 5,
    "next_thought_needed": true
  }'
```

### Mental Models

```bash
curl -X POST http://localhost:8080/api/v1/thinking/mental-model \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "model_name": "first_principles",
    "problem": "How to optimize system performance?",
    "steps": ["Break down the system", "Identify bottlenecks", "Optimize each component"],
    "reasoning": "Starting from basic principles...",
    "conclusion": "Focus on database queries first"
  }'
```

### Stochastic Algorithms

#### Markov Decision Process

```bash
curl -X POST http://localhost:8080/api/v1/stochastic/mdp \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "problem": "Optimize robot navigation policy",
    "states": 100,
    "actions": ["up", "down", "left", "right"],
    "gamma": 0.9,
    "learning_rate": 0.1,
    "epsilon": 0.1,
    "max_iterations": 1000
  }'
```

#### Monte Carlo Tree Search

```bash
curl -X POST http://localhost:8080/api/v1/stochastic/mcts \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "problem": "Find optimal game strategy",
    "simulations": 1000,
    "exploration_constant": 1.4,
    "max_depth": 10
  }'
```

#### Multi-Armed Bandit

```bash
curl -X POST http://localhost:8080/api/v1/stochastic/bandit \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "problem": "Optimize ad placement",
    "arms": 5,
    "strategy": "epsilon-greedy",
    "epsilon": 0.1
  }'
```

### Decision Frameworks

```bash
curl -X POST http://localhost:8080/api/v1/decision/framework \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "decision_statement": "Choose the best technology stack",
    "options": [
      {"name": "React + Node.js", "description": "Full-stack JavaScript"},
      {"name": "Vue + Python", "description": "Python backend with Vue frontend"},
      {"name": "Angular + Java", "description": "Enterprise-grade solution"}
    ],
    "criteria": [
      {"name": "Development Speed", "weight": 0.4, "evaluation_method": "quantitative"},
      {"name": "Maintainability", "weight": 0.3, "evaluation_method": "qualitative"},
      {"name": "Performance", "weight": 0.3, "evaluation_method": "quantitative"}
    ],
    "analysis_type": "multi-criteria"
  }'
```

### Visualization

```bash
curl -X POST http://localhost:8080/api/v1/visual/concept-map \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "diagram_id": "concept-map-1",
    "operation": "create",
    "elements": [
      {"id": "node-1", "type": "node", "label": "Problem", "properties": {"color": "red"}},
      {"id": "node-2", "type": "node", "label": "Solution", "properties": {"color": "green"}},
      {"id": "edge-1", "type": "edge", "source": "node-1", "target": "node-2", "label": "leads to"}
    ],
    "iteration": 1,
    "next_operation_needed": false
  }'
```

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
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── internal/
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   ├── server/            # Server implementation
│   ├── storage/           # Data storage layer
│   └── types/             # Type definitions
└── README.md              # This file
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gothink-linux main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gothink.exe main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gothink-macos main.go
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
