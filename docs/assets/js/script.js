// Theme Management
class ThemeManager {
  constructor() {
    this.theme = localStorage.getItem('theme') || 'dark';
    this.init();
  }

  init() {
    this.setTheme(this.theme);
    this.bindEvents();
  }

  setTheme(theme) {
    this.theme = theme;
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
    this.updateThemeIcon();
  }

  toggleTheme() {
    const newTheme = this.theme === 'dark' ? 'light' : 'dark';
    this.setTheme(newTheme);
  }

  updateThemeIcon() {
    const themeIcon = document.querySelector('.theme-icon');
    if (themeIcon) {
      themeIcon.textContent = this.theme === 'dark' ? '☀️' : '🌙';
    }
  }

  bindEvents() {
    const themeToggle = document.getElementById('themeToggle');
    if (themeToggle) {
      themeToggle.addEventListener('click', () => this.toggleTheme());
    }
  }
}

// Smooth Scrolling
class SmoothScroll {
  constructor() {
    this.init();
  }

  init() {
    // Handle anchor links
    document.addEventListener('click', (e) => {
      const link = e.target.closest('a[href^="#"]');
      if (link) {
        e.preventDefault();
        const targetId = link.getAttribute('href').substring(1);
        const targetElement = document.getElementById(targetId);
        if (targetElement) {
          this.scrollToElement(targetElement);
        }
      }
    });
  }

  scrollToElement(element) {
    const offsetTop = element.offsetTop - 80; // Account for fixed navbar
    window.scrollTo({
      top: offsetTop,
      behavior: 'smooth'
    });
  }
}

// Code Block Enhancements
class CodeEnhancer {
  constructor() {
    this.init();
  }

  init() {
    this.addCopyButtons();
    this.enhanceCodeBlocks();
  }

  addCopyButtons() {
    const codeBlocks = document.querySelectorAll('.code-block pre');
    codeBlocks.forEach(block => {
      const button = document.createElement('button');
      button.className = 'copy-button';
      button.innerHTML = '📋';
      button.title = 'Copy code';
      button.addEventListener('click', () => this.copyCode(block, button));
      
      const wrapper = document.createElement('div');
      wrapper.className = 'code-wrapper';
      wrapper.style.position = 'relative';
      
      block.parentNode.insertBefore(wrapper, block);
      wrapper.appendChild(block);
      wrapper.appendChild(button);
    });
  }

  copyCode(block, button) {
    const text = block.textContent;
    navigator.clipboard.writeText(text).then(() => {
      const originalText = button.innerHTML;
      button.innerHTML = '✅';
      button.style.color = '#10b981';
      setTimeout(() => {
        button.innerHTML = originalText;
        button.style.color = '';
      }, 2000);
    }).catch(err => {
      console.error('Failed to copy code: ', err);
    });
  }

  enhanceCodeBlocks() {
    // Add syntax highlighting if Prism is available
    if (typeof Prism !== 'undefined') {
      Prism.highlightAll();
    }
  }
}

// Navigation Enhancement
class NavigationEnhancer {
  constructor() {
    this.init();
  }

  init() {
    this.handleScroll();
    this.bindEvents();
  }

  handleScroll() {
    let lastScrollY = window.scrollY;
    
    window.addEventListener('scroll', () => {
      const currentScrollY = window.scrollY;
      const navbar = document.querySelector('.navbar');
      
      if (currentScrollY > 100) {
        navbar.style.background = 'rgba(15, 23, 42, 0.95)';
        navbar.style.backdropFilter = 'blur(20px)';
      } else {
        navbar.style.background = 'rgba(15, 23, 42, 0.8)';
        navbar.style.backdropFilter = 'blur(10px)';
      }
      
      lastScrollY = currentScrollY;
    });
  }

  bindEvents() {
    // Mobile menu toggle (if needed in future)
    const mobileMenuToggle = document.getElementById('mobileMenuToggle');
    if (mobileMenuToggle) {
      mobileMenuToggle.addEventListener('click', () => {
        const navLinks = document.querySelector('.nav-links');
        navLinks.classList.toggle('mobile-open');
      });
    }
  }
}

// Animation on Scroll
class ScrollAnimations {
  constructor() {
    this.init();
  }

  init() {
    this.observeElements();
  }

  observeElements() {
    const observerOptions = {
      threshold: 0.1,
      rootMargin: '0px 0px -50px 0px'
    };

    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('animate-in');
        }
      });
    }, observerOptions);

    // Observe elements that should animate
    const animateElements = document.querySelectorAll('.feature-card, .model-card, .tool-card, .example-card');
    animateElements.forEach(el => {
      el.classList.add('animate-element');
      observer.observe(el);
    });
  }
}

// Performance Optimizations
class PerformanceOptimizer {
  constructor() {
    this.init();
  }

  init() {
    this.lazyLoadImages();
    this.preloadCriticalResources();
  }

  lazyLoadImages() {
    if ('IntersectionObserver' in window) {
      const imageObserver = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            const img = entry.target;
            img.src = img.dataset.src;
            img.classList.remove('lazy');
            imageObserver.unobserve(img);
          }
        });
      });

      const lazyImages = document.querySelectorAll('img[data-src]');
      lazyImages.forEach(img => imageObserver.observe(img));
    }
  }

  preloadCriticalResources() {
    // Preload critical fonts
    const fontLink = document.createElement('link');
    fontLink.rel = 'preload';
    fontLink.href = 'https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap';
    fontLink.as = 'style';
    document.head.appendChild(fontLink);
  }
}

// Error Handling
class ErrorHandler {
  constructor() {
    this.init();
  }

  init() {
    window.addEventListener('error', (e) => {
      console.error('JavaScript error:', e.error);
      // Could send to analytics service here
    });

    window.addEventListener('unhandledrejection', (e) => {
      console.error('Unhandled promise rejection:', e.reason);
      // Could send to analytics service here
    });
  }
}

// Initialize everything when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  new ThemeManager();
  new SmoothScroll();
  new CodeEnhancer();
  new NavigationEnhancer();
  new ScrollAnimations();
  new PerformanceOptimizer();
  new ErrorHandler();
});

// Add CSS for animations
const style = document.createElement('style');
style.textContent = `
  .animate-element {
    opacity: 0;
    transform: translateY(20px);
    transition: opacity 0.6s ease, transform 0.6s ease;
  }
  
  .animate-element.animate-in {
    opacity: 1;
    transform: translateY(0);
  }
  
  .copy-button {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 0.25rem;
    padding: 0.5rem;
    cursor: pointer;
    font-size: 0.875rem;
    transition: all 0.3s ease;
    z-index: 10;
  }
  
  .copy-button:hover {
    background: var(--accent-primary);
    color: white;
    border-color: var(--accent-primary);
  }
  
  .code-wrapper {
    position: relative;
  }
  
  .lazy {
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  .lazy.loaded {
    opacity: 1;
  }
`;
document.head.appendChild(style);
