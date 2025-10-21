#!/usr/bin/env node

/**
 * Mental Model Validation Script
 * 
 * This script validates mental model submissions to ensure they meet
 * the quality standards and formatting requirements.
 * 
 * Usage: node scripts/validate-mental-model.js <model-file>
 */

const fs = require('fs');
const path = require('path');

// Validation rules
const VALIDATION_RULES = {
  requiredFields: ['name', 'description', 'steps', 'category', 'author'],
  categories: [
    'analytical', 'decision-making', 'creative', 'strategic', 
    'scientific', 'collaborative', 'performance', 'systems'
  ],
  minSteps: 3,
  maxSteps: 10,
  minDescriptionLength: 20,
  maxDescriptionLength: 200,
  maxNameLength: 50
};

function validateMentalModel(modelData) {
  const errors = [];
  const warnings = [];

  // Check required fields
  for (const field of VALIDATION_RULES.requiredFields) {
    if (!modelData[field]) {
      errors.push(`Missing required field: ${field}`);
    }
  }

  // Validate name
  if (modelData.name) {
    if (modelData.name.length > VALIDATION_RULES.maxNameLength) {
      errors.push(`Name too long (max ${VALIDATION_RULES.maxNameLength} characters)`);
    }
    if (modelData.name.length < 3) {
      errors.push('Name too short (min 3 characters)');
    }
  }

  // Validate description
  if (modelData.description) {
    if (modelData.description.length < VALIDATION_RULES.minDescriptionLength) {
      errors.push(`Description too short (min ${VALIDATION_RULES.minDescriptionLength} characters)`);
    }
    if (modelData.description.length > VALIDATION_RULES.maxDescriptionLength) {
      errors.push(`Description too long (max ${VALIDATION_RULES.maxDescriptionLength} characters)`);
    }
  }

  // Validate category
  if (modelData.category) {
    if (!VALIDATION_RULES.categories.includes(modelData.category.toLowerCase())) {
      errors.push(`Invalid category. Must be one of: ${VALIDATION_RULES.categories.join(', ')}`);
    }
  }

  // Validate steps
  if (modelData.steps) {
    if (!Array.isArray(modelData.steps)) {
      errors.push('Steps must be an array');
    } else {
      if (modelData.steps.length < VALIDATION_RULES.minSteps) {
        errors.push(`Too few steps (min ${VALIDATION_RULES.minSteps})`);
      }
      if (modelData.steps.length > VALIDATION_RULES.maxSteps) {
        errors.push(`Too many steps (max ${VALIDATION_RULES.maxSteps})`);
      }
      
      // Check each step
      modelData.steps.forEach((step, index) => {
        if (typeof step !== 'string') {
          errors.push(`Step ${index + 1} must be a string`);
        } else if (step.trim().length < 10) {
          warnings.push(`Step ${index + 1} is very short and might not be actionable`);
        } else if (step.trim().length > 200) {
          warnings.push(`Step ${index + 1} is very long and might be too complex`);
        }
      });
    }
  }

  // Validate author
  if (modelData.author) {
    if (!modelData.author.startsWith('@')) {
      warnings.push('Author should start with @ (e.g., @username)');
    }
  }

  // Check for common issues
  if (modelData.description && modelData.description.includes('TODO')) {
    warnings.push('Description contains TODO - please complete before submission');
  }

  if (modelData.steps && Array.isArray(modelData.steps)) {
    const hasEmptySteps = modelData.steps.some(step => !step.trim());
    if (hasEmptySteps) {
      errors.push('Some steps are empty');
    }
  }

  return { errors, warnings };
}

function parseModelFromMarkdown(markdownContent) {
  // This is a simple parser for the markdown format used in the community models page
  // In a real implementation, you might want to use a proper markdown parser
  
  const modelData = {};
  
  // Extract name from h3 tag
  const nameMatch = markdownContent.match(/<h3>(.*?)<\/h3>/);
  if (nameMatch) {
    modelData.name = nameMatch[1].trim();
  }
  
  // Extract category
  const categoryMatch = markdownContent.match(/<p class="model-category">(.*?)<\/p>/);
  if (categoryMatch) {
    modelData.category = categoryMatch[1].trim().toLowerCase();
  }
  
  // Extract author
  const authorMatch = markdownContent.match(/<p class="model-author">by <a[^>]*>@([^<]+)<\/a><\/p>/);
  if (authorMatch) {
    modelData.author = `@${authorMatch[1]}`;
  }
  
  // Extract description
  const descMatch = markdownContent.match(/<p>(.*?)<\/p>/);
  if (descMatch && !descMatch[1].includes('model-category') && !descMatch[1].includes('model-author')) {
    modelData.description = descMatch[1].trim();
  }
  
  // Extract steps
  const stepsMatch = markdownContent.match(/<ul class="model-steps">(.*?)<\/ul>/s);
  if (stepsMatch) {
    const stepMatches = stepsMatch[1].match(/<li>(.*?)<\/li>/g);
    if (stepMatches) {
      modelData.steps = stepMatches.map(step => 
        step.replace(/<\/?li>/g, '').trim()
      );
    }
  }
  
  return modelData;
}

function main() {
  const args = process.argv.slice(2);
  
  if (args.length === 0) {
    console.log('Usage: node scripts/validate-mental-model.js <model-file>');
    console.log('');
    console.log('This script validates mental model submissions.');
    console.log('It can parse both JSON and Markdown formats.');
    process.exit(1);
  }
  
  const filePath = args[0];
  
  if (!fs.existsSync(filePath)) {
    console.error(`Error: File not found: ${filePath}`);
    process.exit(1);
  }
  
  const content = fs.readFileSync(filePath, 'utf8');
  let modelData;
  
  try {
    // Try to parse as JSON first
    modelData = JSON.parse(content);
  } catch (e) {
    // If not JSON, try to parse as Markdown
    modelData = parseModelFromMarkdown(content);
  }
  
  console.log(`Validating mental model: ${modelData.name || 'Unknown'}`);
  console.log('');
  
  const { errors, warnings } = validateMentalModel(modelData);
  
  if (errors.length > 0) {
    console.log('❌ VALIDATION ERRORS:');
    errors.forEach(error => console.log(`  • ${error}`));
    console.log('');
  }
  
  if (warnings.length > 0) {
    console.log('⚠️  WARNINGS:');
    warnings.forEach(warning => console.log(`  • ${warning}`));
    console.log('');
  }
  
  if (errors.length === 0 && warnings.length === 0) {
    console.log('✅ Model validation passed!');
    console.log('');
    console.log('This model is ready for review and potential inclusion.');
  } else if (errors.length === 0) {
    console.log('✅ Model validation passed with warnings.');
    console.log('');
    console.log('This model can be included, but consider addressing the warnings.');
  } else {
    console.log('❌ Model validation failed.');
    console.log('');
    console.log('Please fix the errors before submitting.');
    process.exit(1);
  }
  
  // Print model summary
  console.log('📋 Model Summary:');
  console.log(`  Name: ${modelData.name || 'N/A'}`);
  console.log(`  Category: ${modelData.category || 'N/A'}`);
  console.log(`  Author: ${modelData.author || 'N/A'}`);
  console.log(`  Steps: ${modelData.steps ? modelData.steps.length : 0}`);
  console.log(`  Description: ${modelData.description ? modelData.description.substring(0, 100) + '...' : 'N/A'}`);
}

if (require.main === module) {
  main();
}

module.exports = { validateMentalModel, parseModelFromMarkdown };
