# GoThink Documentation

This directory contains the documentation site for GoThink, built with Jekyll and designed for GitHub Pages deployment.

## Structure

```
docs/
├── _config.yml          # Jekyll configuration
├── _layouts/            # Jekyll layouts
│   └── default.html     # Main layout template
├── _docs/               # Documentation pages
│   └── mental-models.md # Mental models documentation
├── assets/              # Static assets
│   ├── css/
│   │   └── styles.css   # Custom styles
│   └── js/
│       └── script.js    # Custom JavaScript
├── index.html           # Homepage
├── Gemfile              # Jekyll dependencies
└── README.md            # This file
```

## Features

- **Dark Mode Default**: Modern dark theme with light mode toggle
- **Responsive Design**: Works on all device sizes
- **Jekyll Integration**: Full Jekyll compatibility for GitHub Pages
- **Syntax Highlighting**: Code blocks with Prism.js
- **Smooth Animations**: CSS animations and transitions
- **Copy Code**: One-click code copying
- **SEO Optimized**: Meta tags and structured data

## Local Development

1. Install Ruby and Bundler
2. Install dependencies:
   ```bash
   bundle install
   ```
3. Serve locally:
   ```bash
   bundle exec jekyll serve
   ```
4. Open http://localhost:4000

## Deployment

The site is automatically deployed to GitHub Pages when changes are pushed to the main branch. The deployment is handled by the GitHub Actions workflow in `.github/workflows/pages.yml`.

## Customization

### Adding New Pages
1. Create a new Markdown file in `_docs/`
2. Add front matter with layout and metadata
3. Add navigation links in `_config.yml`

### Styling
- Main styles: `assets/css/styles.css`
- Uses CSS custom properties for theming
- Responsive design with mobile-first approach

### JavaScript
- Main script: `assets/js/script.js`
- Modular class-based architecture
- Theme management, smooth scrolling, and animations

## GitHub Pages Configuration

The site is configured for GitHub Pages with:
- Jekyll 4.3.0
- Minima theme (as fallback)
- GitHub Pages compatible plugins
- Proper baseurl configuration for project pages

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- CSS Grid and Flexbox support required
- JavaScript ES6+ features used
