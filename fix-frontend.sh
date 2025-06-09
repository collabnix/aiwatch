#!/bin/bash

echo "=== Fixing Frontend Light Mode (Manual Method) ==="
echo

# Check if frontend directory exists
if [ ! -d "frontend/src" ]; then
    echo "❌ Frontend src directory not found"
    exit 1
fi

echo "🔍 Let's check your current App.tsx..."
if [ -f "frontend/src/App.tsx" ]; then
    echo "📄 Current App.tsx content:"
    head -20 frontend/src/App.tsx
    echo "..."
    echo
else
    echo "❌ App.tsx not found. Let's create it."
fi

echo "📝 Creating updated App.tsx with light mode default..."

# Create the corrected App.tsx
cat > frontend/src/App.tsx << 'EOF'
import React, { useState, useEffect } from 'react';
import ChatInterface from './components/ChatInterface';
import './App.css';

function App() {
  // Default to light mode instead of dark mode
  const [isDarkMode, setIsDarkMode] = useState(false);

  useEffect(() => {
    // Check if user has a saved preference, otherwise default to light
    const savedTheme = localStorage.getItem('aiwatch-theme');
    if (savedTheme) {
      setIsDarkMode(savedTheme === 'dark');
    } else {
      // Default to light mode
      setIsDarkMode(false);
      localStorage.setItem('aiwatch-theme', 'light');
    }
  }, []);

  const toggleTheme = () => {
    const newTheme = !isDarkMode;
    setIsDarkMode(newTheme);
    localStorage.setItem('aiwatch-theme', newTheme ? 'dark' : 'light');
  };

  return (
    <div className={`app ${isDarkMode ? 'dark' : 'light'}`}>
      <div className="app-header">
        <div className="header-content">
          <div className="logo-section">
            <h1 className="app-title">
              🤖 AIWatch
            </h1>
            <span className="app-subtitle">
              Intelligent AI Model Monitoring & Chat
            </span>
          </div>
          <button 
            onClick={toggleTheme}
            className="theme-toggle"
            aria-label="Toggle theme"
          >
            {isDarkMode ? '☀️' : '🌙'}
          </button>
        </div>
      </div>
      
      <main className="app-main">
        <ChatInterface isDarkMode={isDarkMode} />
      </main>
      
      <footer className="app-footer">
        <div className="footer-content">
          <span>Powered by AIWatch • </span>
          <a href="http://localhost:3001" target="_blank" rel="noopener noreferrer">
            📊 Monitoring Dashboard
          </a>
          <span> • </span>
          <a href="http://localhost:9091" target="_blank" rel="noopener noreferrer">
            📈 Metrics
          </a>
        </div>
      </footer>
    </div>
  );
}

export default App;
EOF

echo "✅ Created new App.tsx with light mode default"

echo
echo "📝 Creating optimized App.css for light mode..."

# Create light mode optimized CSS
cat > frontend/src/App.css << 'EOF'
/* AIWatch App Styles - Light Mode Optimized */

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: all 0.3s ease;
}

/* Light Mode (Default) */
.app.light {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  color: #2d3748;
}

.app.light .app-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.app.light .app-title {
  color: #2b6cb0;
}

.app.light .app-subtitle {
  color: #4a5568;
}

.app.light .theme-toggle {
  background: rgba(0, 0, 0, 0.05);
  color: #4a5568;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.app.light .theme-toggle:hover {
  background: rgba(0, 0, 0, 0.1);
  transform: scale(1.05);
}

.app.light .app-footer {
  background: rgba(255, 255, 255, 0.9);
  color: #4a5568;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.app.light .footer-content a {
  color: #2b6cb0;
}

.app.light .footer-content a:hover {
  color: #2c5aa0;
}

/* Dark Mode */
.app.dark {
  background: linear-gradient(135deg, #0c0c0c 0%, #1a1a1a 100%);
  color: #e2e8f0;
}

.app.dark .app-header {
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.app.dark .app-title {
  color: #63b3ed;
}

.app.dark .app-subtitle {
  color: #a0aec0;
}

.app.dark .theme-toggle {
  background: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.app.dark .theme-toggle:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

.app.dark .app-footer {
  background: rgba(0, 0, 0, 0.8);
  color: #a0aec0;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.app.dark .footer-content a {
  color: #63b3ed;
}

.app.dark .footer-content a:hover {
  color: #90cdf4;
}

/* Header Styles */
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  padding: 1rem 2rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.logo-section {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.app-title {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.025em;
}

.app-subtitle {
  font-size: 0.875rem;
  font-weight: 500;
  opacity: 0.8;
}

.theme-toggle {
  padding: 0.5rem;
  border-radius: 0.5rem;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  transition: all 0.2s ease;
}

.theme-toggle:hover {
  transform: scale(1.05);
}

/* Main Content */
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

/* Footer */
.app-footer {
  padding: 1rem 2rem;
  backdrop-filter: blur(10px);
}

.footer-content {
  display: flex;
  justify-content: center;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  font-size: 0.875rem;
}

.footer-content a {
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s ease;
}

/* Responsive Design */
@media (max-width: 768px) {
  .app-header {
    padding: 1rem;
  }
  
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .app-main {
    padding: 1rem;
  }
  
  .app-footer {
    padding: 1rem;
  }
  
  .footer-content {
    flex-direction: column;
    gap: 0.5rem;
    text-align: center;
  }
}

/* Smooth transitions for theme switching */
* {
  transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease;
}
EOF

echo "✅ Created optimized App.css for light mode"

echo
echo "🔍 Checking if ChatInterface component exists..."
if [ ! -f "frontend/src/components/ChatInterface.tsx" ]; then
    echo "📝 Creating ChatInterface component directory..."
    mkdir -p frontend/src/components
    
    echo "📝 Creating basic ChatInterface component..."
    cat > frontend/src/components/ChatInterface.tsx << 'EOF'
import React, { useState } from 'react';

interface ChatInterfaceProps {
  isDarkMode: boolean;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({ isDarkMode }) => {
  const [message, setMessage] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Message:', message);
    setMessage('');
  };

  return (
    <div style={{
      background: isDarkMode ? '#1a1a1a' : '#ffffff',
      color: isDarkMode ? '#e2e8f0' : '#2d3748',
      padding: '2rem',
      borderRadius: '1rem',
      border: `1px solid ${isDarkMode ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)'}`,
      boxShadow: '0 8px 32px rgba(0,0,0,0.1)',
      minHeight: '400px',
      display: 'flex',
      flexDirection: 'column',
      gap: '1rem'
    }}>
      <div style={{
        padding: '1rem',
        background: isDarkMode ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.02)',
        borderRadius: '0.5rem',
        textAlign: 'center'
      }}>
        <h2>🤖 AI Chat Interface</h2>
        <p style={{ opacity: 0.7, fontSize: '0.875rem' }}>
          Light mode is now the default! Toggle with the 🌙 button above.
        </p>
      </div>
      
      <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div style={{ textAlign: 'center', opacity: 0.6 }}>
          <p>Chat interface ready!</p>
          <p style={{ fontSize: '0.875rem' }}>Start typing your message below...</p>
        </div>
      </div>
      
      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '0.5rem' }}>
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type your message..."
          style={{
            flex: 1,
            padding: '0.75rem',
            borderRadius: '0.5rem',
            border: `1px solid ${isDarkMode ? 'rgba(255,255,255,0.2)' : 'rgba(0,0,0,0.2)'}`,
            background: isDarkMode ? 'rgba(255,255,255,0.05)' : '#ffffff',
            color: isDarkMode ? '#e2e8f0' : '#2d3748',
            outline: 'none'
          }}
        />
        <button
          type="submit"
          disabled={!message.trim()}
          style={{
            padding: '0.75rem 1.5rem',
            borderRadius: '0.5rem',
            border: 'none',
            background: '#3b82f6',
            color: 'white',
            cursor: message.trim() ? 'pointer' : 'not-allowed',
            opacity: message.trim() ? 1 : 0.5
          }}
        >
          Send
        </button>
      </form>
    </div>
  );
};

export default ChatInterface;
EOF
    
    echo "✅ Created basic ChatInterface component"
else
    echo "✅ ChatInterface component already exists"
fi

echo
echo "🎨 Light Mode Features Applied:"
echo "  ✅ Default light theme (no more dark mode on startup)"
echo "  ✅ Beautiful gradient background"
echo "  ✅ Professional header with AIWatch branding"
echo "  ✅ Theme toggle button (🌙/☀️)"
echo "  ✅ Links to monitoring dashboard and metrics"
echo "  ✅ Responsive design"
echo "  ✅ Smooth transitions"

echo
echo "🚀 Now restart your frontend:"
echo "cd frontend"
echo "npm start"
echo
echo "🌐 Visit http://localhost:3000 to see your light mode interface!"
echo
echo "✅ Light mode fix complete!"
