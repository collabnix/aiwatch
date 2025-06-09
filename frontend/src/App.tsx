import React, { useState, useEffect } from 'react';
import ChatInterface from './components/ChatInterface';
import './App.css';

function App() {
  // Start with light mode as default
  const [isDarkMode, setIsDarkMode] = useState(false);

  useEffect(() => {
    // Check localStorage for saved preference
    const savedTheme = localStorage.getItem('aiwatch-theme');
    if (savedTheme) {
      const isDark = savedTheme === 'dark';
      setIsDarkMode(isDark);
    } else {
      // Default to light mode and save it
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
            aria-label={`Switch to ${isDarkMode ? 'light' : 'dark'} mode`}
            title={`Currently in ${isDarkMode ? 'dark' : 'light'} mode. Click to switch.`}
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
          <span>Powered by AIWatch</span>
          <span className="separator">•</span>
          <a href="http://localhost:3001" target="_blank" rel="noopener noreferrer">
            📊 Grafana Dashboard
          </a>
          <span className="separator">•</span>
          <a href="http://localhost:9091" target="_blank" rel="noopener noreferrer">
            📈 Prometheus Metrics
          </a>
        </div>
      </footer>
    </div>
  );
}

export default App;
