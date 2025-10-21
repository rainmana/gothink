# Contributing to GoThink

Thank you for your interest in contributing to GoThink! We welcome contributions from the community, whether it's mental models, bug fixes, documentation improvements, or new features.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Types of Contributions](#types-of-contributions)
- [Mental Models Contributions](#mental-models-contributions)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Community Guidelines](#community-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gothink.git
   cd gothink
   ```
3. **Create a new branch** for your contribution:
   ```bash
   git checkout -b feature/your-feature-name
   # or for mental models:
   git checkout -b add-mental-model-model-name
   ```

## Types of Contributions

### 🧠 Mental Models
- Submit new thinking frameworks
- Improve existing models
- Add examples and use cases
- Categorize and tag models

### 🐛 Bug Reports
- Report issues with clear reproduction steps
- Include system information
- Provide error messages and logs

### ✨ Feature Requests
- Suggest new functionality
- Propose improvements
- Request new mental model categories

### 📚 Documentation
- Improve existing documentation
- Add examples and tutorials
- Fix typos and clarify instructions

### 🔧 Code Contributions
- Fix bugs in the Go codebase
- Add new features
- Improve performance
- Add tests

## Mental Models Contributions

### What Makes a Good Mental Model?

A good mental model should be:

- **Actionable**: Clear, step-by-step process that others can follow
- **Tested**: You've used it successfully in real situations
- **Original**: Either your creation or a significant improvement on existing models
- **Well-documented**: Clear description, examples, and use cases
- **Categorized**: Fits into one of our established categories

### Categories

We organize mental models into these categories:

- **🧠 Analytical**: Problem analysis, data interpretation, systematic thinking
- **🎯 Decision-making**: Frameworks for making better decisions
- **🎨 Creative**: Innovation, brainstorming, creative problem-solving
- **🧮 Strategic**: Long-term planning, competitive analysis, business strategy
- **🧪 Scientific**: Evidence-based approaches, research methods, hypothesis testing
- **🤝 Collaborative**: Team dynamics, communication, group decision-making
- **⚡ Performance**: Productivity, optimization, efficiency
- **🌍 Systems**: Complex systems thinking, interdependencies, emergent behavior

### Submission Format

When submitting a mental model, use this format:

```markdown
<div class="model-card">
  <h3>Your Model Name</h3>
  <p class="model-category">Category</p>
  <p class="model-author">by <a href="https://github.com/yourusername" target="_blank">@yourusername</a></p>
  <p>Brief description of what this model does and when to use it.</p>
  <ul class="model-steps">
    <li>Step 1: Clear, actionable step</li>
    <li>Step 2: Another clear step</li>
    <li>Step 3: Continue with more steps</li>
  </ul>
  <div class="model-meta">
    <span class="model-tags">#tag1 #tag2 #tag3</span>
  </div>
</div>
```

### Example Submission

```markdown
<div class="model-card">
  <h3>Root Cause Analysis</h3>
  <p class="model-category">Analytical</p>
  <p class="model-author">by <a href="https://github.com/johndoe" target="_blank">@johndoe</a></p>
  <p>Systematically identify the underlying causes of problems rather than just treating symptoms.</p>
  <ul class="model-steps">
    <li>Define the problem clearly and specifically</li>
    <li>Collect data and evidence about the problem</li>
    <li>Ask "why" five times to drill down to root causes</li>
    <li>Identify the fundamental root cause(s)</li>
    <li>Develop solutions that address root causes</li>
    <li>Implement and monitor the solutions</li>
  </ul>
  <div class="model-meta">
    <span class="model-tags">#problem-solving #analysis #systematic #troubleshooting</span>
  </div>
</div>
```

## Development Setup

### Prerequisites

- Go 1.23 or later
- Git
- Basic understanding of Go and MCP (Model Context Protocol)

### Building the Project

```bash
# Install dependencies
go mod tidy

# Build the MCP server
go build -o gothink-mcp mcp_main.go

# Run tests
go test ./...

# Test the MCP server
./gothink-mcp
```

### Testing Mental Models

```bash
# Test with custom mental models
export GOTHINK_MENTAL_MODELS_PATH="/path/to/your/mental_models.yaml"
./gothink-mcp
```

## Pull Request Process

### Before Submitting

1. **Test your changes** thoroughly
2. **Update documentation** if needed
3. **Add tests** for new functionality
4. **Follow coding standards** (run `gofmt` and `golint`)
5. **Check that tests pass**: `go test ./...`

### PR Requirements

1. **Clear title** describing the change
2. **Detailed description** of what was changed and why
3. **Reference any issues** that this PR addresses
4. **Screenshots** for UI changes
5. **Testing instructions** for reviewers

### Review Process

1. **Automated checks** must pass (tests, linting, formatting)
2. **Community review** for feedback and suggestions
3. **Maintainer review** for final approval
4. **Merge** once approved

### PR Templates

We provide templates for different types of contributions:

- **Mental Model**: `.github/pull_request_template_mental_model.md`
- **Bug Fix**: `.github/pull_request_template_bug_fix.md`
- **Feature**: `.github/pull_request_template_feature.md`
- **Documentation**: `.github/pull_request_template_documentation.md`

## Community Guidelines

### Communication

- **Be respectful** and constructive in all interactions
- **Ask questions** if you're unsure about something
- **Help others** when you can
- **Use inclusive language** that welcomes everyone

### Recognition

Contributors are recognized in several ways:

- **Attribution** on mental models with GitHub profile links
- **Contributors page** listing all community members
- **Release notes** mentioning significant contributions
- **Special badges** for major contributors

### Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and community chat
- **Documentation**: Check our comprehensive guides
- **Examples**: Review existing mental models for inspiration

## License

By contributing to GoThink, you agree that your contributions will be licensed under the same license as the project (see [LICENSE](LICENSE) file).

## Thank You

Thank you for contributing to GoThink! Your contributions help make this project better for everyone in the community.

---

*Questions? Feel free to open an issue or start a discussion!*
