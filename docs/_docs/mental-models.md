---
layout: doc
title: "Mental Models"
description: "Comprehensive guide to GoThink's mental models system including core models, custom models, and YAML configuration."
---

# Mental Models

GoThink's mental models system provides a powerful framework for structured thinking and problem-solving. The system includes built-in core models and supports custom models via YAML configuration.

## Core Mental Models

GoThink comes with four built-in mental models that are always available:

### 1. First Principles Thinking
**Category:** Analytical  
**Description:** Break down complex problems into fundamental components

**Steps:**
1. Identify the problem clearly
2. Break it down into basic components
3. Question assumptions
4. Build up from the basics

**Use Case:** When you need to understand the fundamental nature of a problem and build solutions from the ground up.

### 2. Opportunity Cost Analysis
**Category:** Decision-making  
**Description:** Consider what you give up when making a choice

**Steps:**
1. Identify all available options
2. List the benefits of each option
3. Identify what you give up with each choice
4. Compare opportunity costs

**Use Case:** When making decisions between multiple options and need to understand the trade-offs.

### 3. Bayesian Thinking
**Category:** Probabilistic  
**Description:** Update beliefs based on new evidence

**Steps:**
1. Start with prior beliefs
2. Gather new evidence
3. Update beliefs using Bayes' theorem
4. Consider alternative explanations

**Use Case:** When dealing with uncertainty and need to update your understanding based on new information.

### 4. Systems Thinking
**Category:** Holistic  
**Description:** Understand how parts of a system interact

**Steps:**
1. Identify system boundaries
2. Map system components
3. Identify relationships and feedback loops
4. Consider emergent properties

**Use Case:** When dealing with complex systems and need to understand how different parts interact.

## Custom Mental Models

You can create your own mental models using YAML configuration. This allows you to:

- Define custom problem-solving frameworks
- Override core models with enhanced versions
- Organize models by category
- Set priority levels for model selection

### YAML Configuration Format

```yaml
models:
  model_key:
    name: "Display Name"
    description: "Description of the model"
    steps:
      - "Step 1: Description"
      - "Step 2: Description"
      - "Step 3: Description"
    category: "category-name"
    priority: 10  # Higher = more priority
```

### Configuration Options

- **`name`** (required): The display name of the mental model
- **`description`** (required): A description of what the model does
- **`steps`** (required): Array of step descriptions
- **`category`** (required): Category for organization (e.g., "analytical", "creative", "decision-making")
- **`priority`** (optional): Priority level (default: 1 for custom models, 0 for core models)

### Example Custom Models

#### Problem-Solving Framework
```yaml
models:
  custom_problem_solving:
    name: "Custom Problem Solving Framework"
    description: "A comprehensive approach to solving complex problems"
    steps:
      - "Step 1: Define the problem clearly and specifically"
      - "Step 2: Gather all relevant information and data"
      - "Step 3: Identify constraints and limitations"
      - "Step 4: Generate multiple potential solutions"
      - "Step 5: Evaluate each solution against criteria"
      - "Step 6: Select the best solution and implement"
      - "Step 7: Monitor results and iterate if needed"
    category: "problem-solving"
    priority: 15
```

#### Enhanced First Principles
```yaml
models:
  first_principles:
    name: "Enhanced First Principles Thinking"
    description: "An enhanced version of first principles thinking with additional steps"
    steps:
      - "Step 1: Identify the problem or question"
      - "Step 2: Break down into fundamental components"
      - "Step 3: Question all assumptions"
      - "Step 4: Research the fundamentals"
      - "Step 5: Build up from first principles"
      - "Step 6: Test the new understanding"
      - "Step 7: Apply and iterate"
    category: "analytical"
    priority: 12
```

## Using Mental Models

### Via MCP Server

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "mental_model",
    "arguments": {
      "session_id": "my-session",
      "model_name": "first_principles",
      "problem": "How to optimize team productivity"
    }
  }
}
```

### Response Format

```json
{
  "status": "success",
  "model_id": "unique-model-id",
  "model_info": {
    "name": "First Principles Thinking",
    "description": "Break down complex problems into fundamental components",
    "category": "analytical",
    "priority": 0
  },
  "steps_used": [
    "Identify the problem clearly",
    "Break it down into basic components",
    "Question assumptions",
    "Build up from the basics"
  ],
  "has_steps": true,
  "session_context": {
    "session_id": "my-session",
    "total_mental_models": 1
  }
}
```

## Configuration

### Environment Variable
```bash
export GOTHINK_MENTAL_MODELS_PATH=/path/to/mental_models.yaml
```

### Configuration File
```json
{
  "mental_models_path": "/path/to/mental_models.yaml"
}
```

## Best Practices

### 1. Model Design
- Keep steps clear and actionable
- Use consistent language and formatting
- Include enough detail to be useful
- Test your models with real problems

### 2. Categories
- Use consistent category names
- Group related models together
- Consider your use cases when categorizing

### 3. Priorities
- Use higher priorities (10+) for your most important models
- Core models have priority 0
- Custom models default to priority 1
- Consider the frequency of use when setting priorities

### 4. Naming
- Use descriptive, memorable names
- Avoid special characters in model keys
- Use snake_case for model keys
- Use Title Case for display names

## Troubleshooting

### Model Not Found
If you get a "model not found" error:
1. Check the model name spelling
2. Verify the YAML file is valid
3. Ensure the file path is correct
4. Check that the model is defined in the YAML

### Invalid YAML
If you get YAML parsing errors:
1. Validate your YAML syntax
2. Check for proper indentation
3. Ensure all required fields are present
4. Verify string formatting

### Priority Issues
If models aren't appearing in the expected order:
1. Check priority values
2. Ensure priorities are integers
3. Remember: higher numbers = higher priority

## Advanced Features

### Model Override
Custom models can override core models by using the same key name. The custom model will be used instead of the core model.

### Category Organization
Models are automatically organized by category, making it easy to find related models.

### Priority Sorting
Models are sorted by priority (highest first), then by name for models with the same priority.

### Session Tracking
Each mental model usage is tracked per session, allowing you to see which models are being used most frequently.
