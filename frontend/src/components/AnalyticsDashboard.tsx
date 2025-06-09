import React, { useState, useEffect } from 'react';
import { 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend, 
  ResponsiveContainer,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell
} from 'recharts';
import { Activity, Users, Zap, Clock, TrendingUp, AlertTriangle } from 'lucide-react';

interface AnalyticsData {
  active_users_5m: number;
  active_users_1h: number;
  active_sessions: number;
  token_rates: {
    input_per_minute: number;
    output_per_minute: number;
  };
  top_users: UserStats[];
  model_usage: Record<string, ModelStats>;
  response_time_p95: number;
  response_time_p99: number;
  error_rate: number;
  timestamp: number;
}

interface UserStats {
  user_id: string;
  total_input_tokens: number;
  total_output_tokens: number;
  total_sessions: number;
  avg_tokens_per_request: number;
  last_seen: string;
}

interface ModelStats {
  total_requests: number;
  total_input_tokens: number;
  total_output_tokens: number;
  avg_response_time: number;
  avg_tokens_per_second: number;
}

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];

export const AnalyticsDashboard: React.FC = () => {
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [historicalData, setHistoricalData] = useState<any[]>([]);

  useEffect(() => {
    const fetchAnalytics = async () => {
      try {
        const response = await fetch('http://localhost:8081/analytics');
        if (!response.ok) throw new Error('Failed to fetch analytics');
        
        const data = await response.json();
        setAnalyticsData(data);
        
        // Add to historical data for trends
        setHistoricalData(prev => [
          ...prev.slice(-19), // Keep last 19 entries
          {
            time: new Date(data.timestamp * 1000).toLocaleTimeString(),
            activeUsers5m: data.active_users_5m,
            activeUsers1h: data.active_users_1h,
            activeSessions: data.active_sessions,
            inputRate: data.token_rates.input_per_minute,
            outputRate: data.token_rates.output_per_minute,
            responseTimeP95: data.response_time_p95,
            errorRate: data.error_rate * 100, // Convert to percentage
          }
        ]);
        
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    // Initial fetch
    fetchAnalytics();
    
    // Set up polling every 10 seconds
    const interval = setInterval(fetchAnalytics, 10000);
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded flex items-center">
          <AlertTriangle className="mr-2" size={20} />
          Error loading analytics: {error}
        </div>
      </div>
    );
  }

  if (!analyticsData) {
    return <div>No data available</div>;
  }

  const MetricCard: React.FC<{
    title: string;
    value: string | number;
    icon: React.ReactNode;
    trend?: 'up' | 'down' | 'stable';
    subtitle?: string;
  }> = ({ title, value, icon, trend, subtitle }) => (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600 dark:text-gray-400">{title}</p>
          <p className="text-3xl font-bold text-gray-900 dark:text-white">{value}</p>
          {subtitle && <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">{subtitle}</p>}
        </div>
        <div className="flex flex-col items-center">
          <div className="text-blue-600 dark:text-blue-400">{icon}</div>
          {trend && (
            <TrendingUp 
              className={`mt-2 ${
                trend === 'up' ? 'text-green-500' : 
                trend === 'down' ? 'text-red-500 transform rotate-180' : 
                'text-gray-400'
              }`} 
              size={16} 
            />
          )}
        </div>
      </div>
    </div>
  );

  // Prepare model usage data for pie chart
  const modelUsageData = Object.entries(analyticsData.model_usage).map(([model, stats]) => ({
    name: model,
    value: stats.total_requests,
    tokens: stats.total_input_tokens + stats.total_output_tokens,
  }));

  return (
    <div className="max-w-7xl mx-auto p-6 space-y-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Analytics Dashboard
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          Real-time insights into your AI model usage and performance
        </p>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <MetricCard
          title="Active Users (5m)"
          value={analyticsData.active_users_5m}
          icon={<Users size={24} />}
          trend="stable"
        />
        <MetricCard
          title="Active Sessions"
          value={analyticsData.active_sessions}
          icon={<Activity size={24} />}
          trend="up"
        />
        <MetricCard
          title="Response Time (P95)"
          value={`${analyticsData.response_time_p95.toFixed(2)}s`}
          icon={<Clock size={24} />}
          trend={analyticsData.response_time_p95 < 2 ? 'up' : 'down'}
        />
        <MetricCard
          title="Error Rate"
          value={`${(analyticsData.error_rate * 100).toFixed(2)}%`}
          icon={<AlertTriangle size={24} />}
          trend={analyticsData.error_rate < 0.01 ? 'up' : 'down'}
          subtitle={analyticsData.error_rate < 0.01 ? 'Healthy' : 'Needs Attention'}
        />
      </div>

      {/* Token Processing Metrics */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <MetricCard
          title="Input Tokens/min"
          value={analyticsData.token_rates.input_per_minute.toFixed(0)}
          icon={<Zap size={24} />}
          trend="stable"
        />
        <MetricCard
          title="Output Tokens/min"
          value={analyticsData.token_rates.output_per_minute.toFixed(0)}
          icon={<Zap size={24} />}
          trend="stable"
        />
      </div>

      {/* Historical Trends */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            User Activity Trends
          </h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={historicalData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line 
                type="monotone" 
                dataKey="activeUsers5m" 
                stroke="#8884d8" 
                name="Active Users (5m)"
                strokeWidth={2}
              />
              <Line 
                type="monotone" 
                dataKey="activeSessions" 
                stroke="#82ca9d" 
                name="Active Sessions"
                strokeWidth={2}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            Token Processing Rate
          </h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={historicalData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line 
                type="monotone" 
                dataKey="inputRate" 
                stroke="#ff7300" 
                name="Input Tokens/min"
                strokeWidth={2}
              />
              <Line 
                type="monotone" 
                dataKey="outputRate" 
                stroke="#00ff00" 
                name="Output Tokens/min"
                strokeWidth={2}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Model Usage Distribution */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            Model Usage Distribution
          </h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={modelUsageData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {modelUsageData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            Performance Metrics
          </h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={historicalData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis yAxisId="left" orientation="left" />
              <YAxis yAxisId="right" orientation="right" />
              <Tooltip />
              <Legend />
              <Line 
                yAxisId="left"
                type="monotone" 
                dataKey="responseTimeP95" 
                stroke="#ff4444" 
                name="Response Time P95 (s)"
                strokeWidth={2}
              />
              <Line 
                yAxisId="right"
                type="monotone" 
                dataKey="errorRate" 
                stroke="#ff8800" 
                name="Error Rate (%)"
                strokeWidth={2}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Top Users Table */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
        <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
          Top Users by Token Usage
        </h3>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  User ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Input Tokens
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Output Tokens
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Sessions
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Avg Tokens/Request
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Last Seen
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {analyticsData.top_users.map((user, index) => (
                <tr key={user.user_id} className={index % 2 === 0 ? 'bg-gray-50 dark:bg-gray-700' : ''}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                    {user.user_id}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {user.total_input_tokens.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {user.total_output_tokens.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {user.total_sessions}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {user.avg_tokens_per_request.toFixed(1)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {user.last_seen}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Model Performance Table */}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
        <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
          Model Performance Summary
        </h3>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Model
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Total Requests
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Input Tokens
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Output Tokens
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Avg Response Time
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Tokens/Second
                </th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {Object.entries(analyticsData.model_usage).map(([model, stats], index) => (
                <tr key={model} className={index % 2 === 0 ? 'bg-gray-50 dark:bg-gray-700' : ''}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                    {model}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {stats.total_requests.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {stats.total_input_tokens.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {stats.total_output_tokens.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {stats.avg_response_time.toFixed(2)}s
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
                    {stats.avg_tokens_per_second.toFixed(1)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsDashboard;
