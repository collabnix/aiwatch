import { useState, useEffect } from 'react';
import './App.css';
import { ChatBox, Header } from './components';
import { ModelMetadata } from './types';

function App() {
  // Initialize darkMode from localStorage or system preference
  const [darkMode, setDarkMode] = useState(() => {
    try {
      const savedDarkMode = localStorage.getItem('darkMode');
      
      // If we have a saved preference, use it
      if (savedDarkMode !== null) {
        return savedDarkMode === 'true';
      }
      
      // Otherwise, check system preference
      if (window.matchMedia) {
        return window.matchMedia('(prefers-color-scheme: dark)').matches;
      }
      
      // Fallback to light mode
      return false;
    } catch (error) {
      console.warn('Error reading dark mode preference:', error);
      return false;
    }
  });
  
  const [modelInfo, setModelInfo] = useState<ModelMetadata | null>(null);

  // Apply dark mode class immediately when component mounts
  useEffect(() => {
    try {
      if (darkMode) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    } catch (error) {
      console.warn('Error applying dark mode class:', error);
    }
  }, [darkMode]);

  // Listen for changes to system preference
  useEffect(() => {
    try {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      
      const handleChange = (e: MediaQueryListEvent) => {
        // Only update if user hasn't explicitly set a preference
        const savedPreference = localStorage.getItem('darkMode');
        if (savedPreference === null) {
          setDarkMode(e.matches);
        }
      };
      
      // Check if addEventListener is available (modern browsers)
      if (mediaQuery.addEventListener) {
        mediaQuery.addEventListener('change', handleChange);
        return () => mediaQuery.removeEventListener('change', handleChange);
      } else if (mediaQuery.addListener) {
        // Fallback for older browsers
        // @ts-ignore - For older browsers
        mediaQuery.addListener(handleChange);
        return () => {
          // @ts-ignore - For older browsers
          mediaQuery.removeListener(handleChange);
        };
      }
    } catch (error) {
      console.warn('Error setting up system preference listener:', error);
    }
  }, []);

  // Fetch model information on component mount
  useEffect(() => {
    fetchModelInfo();
  }, []);

  // Save dark mode preference and apply class when darkMode state changes
  useEffect(() => {
    try {
      // Save preference to localStorage
      localStorage.setItem('darkMode', darkMode.toString());
      
      // Apply or remove dark class
      if (darkMode) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    } catch (error) {
      console.warn('Error saving dark mode preference:', error);
    }
  }, [darkMode]);

  const fetchModelInfo = async () => {
    try {
      const response = await fetch('http://localhost:8080/health');
      if (response.ok) {
        const data = await response.json();
        if (data.model_info) {
          setModelInfo(data.model_info);
        }
      }
    } catch (e) {
      console.error('Failed to fetch model info:', e);
    }
  };

  const toggleDarkMode = () => {
    setDarkMode(prevMode => !prevMode);
  };

  // Check if this is a llama.cpp model
  const isLlamaCppModel = 
    modelInfo?.modelType === 'llama.cpp' || 
    modelInfo?.model?.toLowerCase().includes('llama');

  return (
    <div className="min-h-screen flex flex-col bg-white dark:bg-gray-900 dark:text-white transition-colors duration-200">
      <Header toggleDarkMode={toggleDarkMode} darkMode={darkMode} />
      <div className="flex-1 p-4">
        <ChatBox />
      </div>
      <footer className="text-center p-4 text-sm text-gray-500 dark:text-gray-400">
        <p>
          Powered by <span className="font-semibold">Docker Model Runner</span>
          {modelInfo && (
            <> running <span className="font-semibold">{modelInfo.model}</span></>
          )}
          {isLlamaCppModel && (
            <> with <span className="text-blue-500 font-semibold">llama.cpp</span> metrics</>
          )}
        </p>
      </footer>
    </div>
  );
}

export default App;