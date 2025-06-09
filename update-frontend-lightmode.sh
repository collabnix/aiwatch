#!/bin/bash

echo "=== Quick Frontend Light Mode Update ==="
echo

# Quick update for existing frontend
if [ ! -d "frontend/src" ]; then
    echo "❌ Frontend src directory not found"
    exit 1
fi

echo "🎨 Quick update to light mode default..."

# Update just the theme initialization in App.tsx
echo "📝 Updating theme default in App.tsx..."

# Find and replace the theme initialization
if [ -f "frontend/src/App.tsx" ]; then
    # Backup original
    cp frontend/src/App.tsx frontend/src/App.tsx.backup
    
    # Replace dark mode detection with light mode default
    sed -i 's/useState(window.matchMedia.*prefers-color-scheme: dark.*matches)/useState(false)/g' frontend/src/App.tsx
    sed -i 's/useState(true)/useState(false)/g' frontend/src/App.tsx
    
    echo "✅ Updated App.tsx theme default"
else
    echo "❌ App.tsx not found"
fi

# Update any theme-related CSS if it exists
if [ -f "frontend/src/App.css" ]; then
    echo "📝 Optimizing CSS for light mode..."
    
    # Add light mode as primary styling
    cat >> frontend/src/App.css << 'EOF'

/* Light Mode Optimization */
body {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  color: #2d3748;
}

.app {
  background: transparent;
}

/* Ensure light mode is default */
:root {
  --primary-bg: #ffffff;
  --primary-text: #2d3748;
  --secondary-bg: #f7fafc;
  --border-color: rgba(0, 0, 0, 0.1);
}
EOF
    
    echo "✅ Added light mode CSS optimizations"
fi

echo
echo "🚀 Restart your frontend:"
echo "cd frontend && npm start"
echo
echo "✅ Quick light mode update complete!"
