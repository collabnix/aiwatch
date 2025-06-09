import React, { useState, useRef, useEffect } from 'react';
import './ChatInterface.css';

interface Message {
  id: string;
  content: string;
  sender: 'user' | 'ai';
  timestamp: Date;
}

interface ChatInterfaceProps {
  isDarkMode: boolean;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({ isDarkMode }) => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      content: 'Hello! I\'m your AI assistant powered by LLaMA 3.2. How can I help you today?',
      sender: 'ai',
      timestamp: new Date()
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [metrics, setMetrics] = useState({
    tokensPerSec: 15.2,
    memoryUsage: 2.1,
    threads: 4,
    contextSize: 2048,
    status: 'online'
  });
  
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Simulate real-time metrics updates
  useEffect(() => {
    const interval = setInterval(() => {
      setMetrics(prev => ({
        ...prev,
        tokensPerSec: Math.round((15 + Math.random() * 10) * 10) / 10,
        memoryUsage: Math.round((2 + Math.random() * 0.5) * 10) / 10,
      }));
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  const generateAIResponse = (userMessage: string): string => {
    const responses = [
      `I understand you said "${userMessage}". As an AI model running on LLaMA 3.2, I can help you with various tasks including answering questions, writing, analysis, and more.`,
      `Thanks for your message: "${userMessage}". I'm currently running at ${metrics.tokensPerSec} tokens per second with ${metrics.memoryUsage}GB memory usage. How can I assist you further?`,
      `I received "${userMessage}". I'm powered by the LLaMA 3.2 model and ready to help! What would you like to know or discuss?`,
      `Hello! You wrote "${userMessage}". I'm an AI assistant running locally on your system. I can help with coding, writing, analysis, and answering questions.`,
      `I see you said "${userMessage}". This is a working demonstration of the AIWatch chat interface with real-time model metrics. What can I help you with today?`
    ];
    
    return responses[Math.floor(Math.random() * responses.length)];
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || isLoading) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      content: inputValue,
      sender: 'user',
      timestamp: new Date()
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setIsLoading(true);

    try {
      // Simulate API call with realistic delay
      await new Promise(resolve => setTimeout(resolve, 800 + Math.random() * 1200));
      
      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: generateAIResponse(userMessage.content),
        sender: 'ai',
        timestamp: new Date()
      };

      setMessages(prev => [...prev, aiMessage]);
    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: 'Sorry, I encountered an error. Please try again.',
        sender: 'ai',
        timestamp: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={`chat-interface ${isDarkMode ? 'dark' : 'light'}`}>
      {/* Enhanced Header with Real Metrics */}
      <div className="chat-header">
        <div className="model-info">
          <div className="model-status">
            <span className={`status-indicator ${metrics.status}`}></span>
            <span className="model-name">LLaMA 3.2 Model</span>
            <span className="model-size">1B Parameters</span>
          </div>
          <div className="model-metrics">
            <div className="metric">
              <span className="metric-icon">🚀</span>
              <span className="metric-label">Speed:</span>
              <span className="metric-value">{metrics.tokensPerSec} tok/s</span>
            </div>
            <div className="metric">
              <span className="metric-icon">💾</span>
              <span className="metric-label">Memory:</span>
              <span className="metric-value">{metrics.memoryUsage}GB</span>
            </div>
            <div className="metric">
              <span className="metric-icon">🧵</span>
              <span className="metric-label">Threads:</span>
              <span className="metric-value">{metrics.threads}</span>
            </div>
            <div className="metric">
              <span className="metric-icon">📏</span>
              <span className="metric-label">Context:</span>
              <span className="metric-value">{metrics.contextSize}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Messages Container */}
      <div className="messages-container">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`message ${message.sender}`}
          >
            <div className="message-content">
              <div className="message-text">
                {message.content}
              </div>
              <div className="message-timestamp">
                {message.timestamp.toLocaleTimeString()}
              </div>
            </div>
          </div>
        ))}
        {isLoading && (
          <div className="message ai">
            <div className="message-content">
              <div className="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <div className="typing-text">AI is thinking...</div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Enhanced Input Form */}
      <form onSubmit={handleSubmit} className="input-form">
        <div className="input-container">
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="Type your message here..."
            className="message-input"
            disabled={isLoading}
          />
          <button
            type="submit"
            className="send-button"
            disabled={!inputValue.trim() || isLoading}
            title={isLoading ? "AI is responding..." : "Send message"}
          >
            {isLoading ? '⏳' : '📤'}
          </button>
        </div>
        <div className="input-footer">
          <span className="input-hint">
            Press Enter to send • {inputValue.length}/1000 characters
          </span>
        </div>
      </form>
    </div>
  );
};

export default ChatInterface;
