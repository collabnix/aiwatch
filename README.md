# AIWatch - AI Model Management and Observability powered by Docker Model Runner

<img width="932" alt="image" src="https://github.com/user-attachments/assets/869dc88c-fb9a-4f25-ba7a-acec2d2c8984" />

A modern, full-stack chat application demonstrating how to integrate React frontend with a Go backend and run local Large Language Models (LLMs) using Docker's Model Runner. This project features a **comprehensive Redis-powered observability stack** with real-time monitoring, analytics, and distributed tracing.

## Overview

<img width="679" alt="image" src="https://github.com/user-attachments/assets/9b3931c2-aab3-421e-a3ca-990117ee545b" />

This project showcases a complete Generative AI interface with enterprise-grade observability that includes:
- React/TypeScript frontend with a responsive chat UI
- Go backend server for API handling  
- Integration with Docker's Model Runner to run Llama 3.2 locally
- **Redis Stack** with TimeSeries for data persistence and analytics
- **Comprehensive observability** with metrics, logging, and tracing
- **NEW: Redis-powered analytics** with real-time performance monitoring
- **Enhanced Docker Compose** setup with full observability stack

## 🔧 Features

- 💬 Interactive chat interface with message history
- 🔄 Real-time streaming responses (tokens appear as they're generated)
- 🌓 Light/dark mode support based on user preference
- 🐳 Dockerized deployment for easy setup and portability
- 🏠 Run AI models locally without cloud API dependencies
- 🔒 Cross-origin resource sharing (CORS) enabled
- 🧪 Integration testing using Testcontainers
- 📊 **Redis-powered metrics** and performance monitoring
- 📝 Structured logging with zerolog
- 🔍 Distributed tracing with OpenTelemetry & Jaeger
- 📈 **Grafana dashboards** for visualization
- 🚀 Advanced llama.cpp performance metrics
- **🆕 Redis Stack** with TimeSeries, Search, and JSON support
- **🆕 Redis Exporter** for Prometheus metrics integration
- **🆕 Token Analytics Service** for usage tracking
- **🆕 Production-ready** health checks and service dependencies


## 🏗️ Enhanced Architecture

The application now consists of a comprehensive observability stack:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │ >>> │   Backend       │ >>> │  Model Runner   │
│  (React/TS)     │     │    (Go)         │     │ (Llama 3.2)     │
└─────────────────┘     └─────────────────┘     └─────────────────┘
      :3000                   :8080                   :12434
                              │  │
┌─────────────────┐    ┌──────────────┐    ┌─────────────────┐
│    Grafana      │ <<< │  Prometheus  │    │     Jaeger      │
│  Dashboards     │     │   Metrics    │    │    Tracing      │
└─────────────────┘     └──────────────┘    └─────────────────┘
      :3001                   :9091                :16686

┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Redis Stack   │    │ Redis Exporter  │    │ Token Analytics │
│ DB + Insight    │    │ (Prometheus)    │    │    Service      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
   :6379, :8001              :9121                   :8082

┌─────────────────┐
│ Redis TimeSeries│
│    Service      │
└─────────────────┘
       :8085
```

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git
- Go 1.19 or higher (for local development)
- Node.js and npm (for frontend development)

Before starting, pull the required model:

```bash
docker model pull ai/llama3.2:1B-Q8_0
```

### 🎯 One-Command Deployment

**Start the complete AIWatch observability stack:**

```bash
# Clone the repository
git clone https://github.com/collabnix/aiwatch.git
cd aiwatch

# Start the complete stack (builds and runs all services)
docker-compose -f compose.enhanced.yaml up -d --build
```

### 🌐 Access Points

After deployment, access these services:

| Service | URL | Credentials | Purpose |
|---------|-----|-------------|---------|
| **AIWatch Frontend** | http://localhost:3000 | - | Main chat interface |
| **Grafana** | http://localhost:3001 | admin/admin | Monitoring dashboards |
| **Redis Insight** | http://localhost:8001 | - | Redis database GUI |
| **Prometheus** | http://localhost:9091 | - | Metrics collection |
| **Jaeger** | http://localhost:16686 | - | Distributed tracing |
| **Token Analytics** | http://localhost:8082 | - | Usage analytics API |
| **TimeSeries API** | http://localhost:8085 | - | Redis TimeSeries service |

### 📊 Redis Observability Features

#### **Redis Stack Components**

1. **Redis Database** (Port 6379)
   - Primary data store for chat history and session management
   - Redis TimeSeries for metrics storage
   - Redis JSON for complex data structures
   - Redis Search for full-text capabilities

2. **Redis Insight** (Port 8001) 
   - Web-based Redis GUI for database inspection
   - Real-time monitoring of Redis performance
   - Key-value browser and query interface

3. **Redis Exporter** (Port 9121)
   - Exports Redis metrics to Prometheus
   - Monitors memory usage, command statistics, connection counts
   - Integration with alerting systems

4. **Token Analytics Service** (Port 8082)
   - Tracks token usage patterns and costs
   - API endpoint for analytics queries
   - Integration with frontend metrics display

5. **Redis TimeSeries Service** (Port 8085)
   - Dedicated API for time-series data operations
   - Historical performance data storage
   - Real-time metrics aggregation

#### **Monitoring & Analytics**

- **Real-time Redis Metrics**: Memory usage, commands/sec, connections
- **Token Usage Analytics**: Input/output tokens, cost tracking, usage patterns  
- **Performance Monitoring**: Response times, throughput, error rates
- **Historical Data**: Time-series storage of all metrics for trend analysis
- **Grafana Integration**: Pre-configured dashboards for Redis monitoring

## 🛠️ Development Setup

### Frontend

The frontend is built with React, TypeScript, and Vite:

```bash
cd frontend
npm install
npm run dev
```

This will start the development server at [http://localhost:3000](http://localhost:3000).

### Backend

The Go backend can be run directly:

```bash
go mod download
go run main.go
```

Make sure to set the required environment variables from `backend.env`:
- `BASE_URL`: URL for the model runner
- `MODEL`: Model identifier to use
- `API_KEY`: API key for authentication (defaults to "ollama")
- `REDIS_ADDR`: Redis connection address (redis:6379)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `LOG_PRETTY`: Whether to output pretty-printed logs
- `TRACING_ENABLED`: Enable OpenTelemetry tracing
- `OTLP_ENDPOINT`: OpenTelemetry collector endpoint

## 🔄 How It Works

1. The frontend sends chat messages to the backend API
2. The backend formats the messages and sends them to the Model Runner
3. Chat history and session data are stored in Redis
4. The LLM processes the input and generates a response
5. The backend streams the tokens back to the frontend as they're generated
6. **Token analytics** are collected and stored in Redis TimeSeries
7. **Redis metrics** are exported to Prometheus for monitoring
8. Observability components collect metrics, logs, and traces throughout the process
9. **Grafana dashboards** provide real-time visualization of system performance

## 📁 Project Structure

```
├── compose.enhanced.yaml        # Complete observability stack
├── backend.env                  # Backend environment variables
├── main.go                     # Go backend server
├── frontend/                   # React frontend application
│   ├── src/                    # Source code
│   │   ├── components/         # React components
│   │   ├── App.tsx            # Main application component
│   │   └── ...
├── pkg/                       # Go packages
│   ├── logger/                # Structured logging
│   ├── metrics/               # Prometheus metrics
│   ├── middleware/            # HTTP middleware
│   ├── tracing/               # OpenTelemetry tracing
│   └── health/                # Health check endpoints
├── prometheus/                # Prometheus configuration
├── grafana/                   # Grafana dashboards and configuration
├── redis/                     # Redis configuration
│   └── redis.conf            # Redis server configuration
├── observability/             # Observability documentation
└── ...
```

## 📈 llama.cpp Metrics Features

The application includes detailed llama.cpp metrics displayed directly in the UI:

- **Tokens per Second**: Real-time generation speed
- **Context Window Size**: Maximum tokens the model can process
- **Prompt Evaluation Time**: Time spent processing the input prompt
- **Memory per Token**: Memory usage efficiency
- **Thread Utilization**: Number of threads used for inference
- **Batch Size**: Inference batch size

These metrics help in understanding the performance characteristics of llama.cpp models and can be used to optimize configurations.

## 🔍 Observability Features

The project includes comprehensive observability features:

### Metrics

- Model performance (latency, time to first token)
- Token usage (input and output counts)
- Request rates and error rates
- Active request monitoring
- **Redis performance metrics** (memory, commands, connections)
- **Token analytics** with cost tracking
- llama.cpp specific performance metrics

### Logging

- Structured JSON logs with zerolog
- Log levels (debug, info, warn, error, fatal)
- Request logging middleware
- Error tracking

### Tracing

- Request flow tracing with OpenTelemetry
- Integration with Jaeger for visualization
- Span context propagation

For more information, see [Observability Documentation](./observability/README.md).

## 🎛️ Configuration Options

### Redis Configuration

The Redis setup includes:
- **Persistence**: RDB and AOF enabled for data durability
- **Memory Optimization**: Configured for optimal performance
- **Security**: Protected mode disabled for development (configure for production)
- **TimeSeries**: Enabled for metrics storage
- **Networking**: Bridge network for service communication

### Service Dependencies

All services include:
- **Health Checks**: Automated service health monitoring
- **Restart Policies**: Automatic restart on failure
- **Resource Limits**: Memory and CPU constraints
- **Logging**: Centralized log collection

## ⚙️ Customization

You can customize the application by:
1. Changing the model in `backend.env` to use a different LLM
2. Modifying the frontend components for a different UI experience
3. Extending the backend API with additional functionality
4. Customizing the Grafana dashboards for different metrics
5. Adjusting llama.cpp parameters for performance optimization
6. **Configuring Redis** for different persistence and performance requirements
7. **Adding custom analytics** using the Token Analytics Service API
8. **Creating custom dashboards** in Grafana for specific monitoring needs

## 🧪 Testing

The project includes integration tests using Testcontainers:

```bash
cd tests
go test -v
```

## 🚨 Troubleshooting

### Common Issues

- **Model not loading**: Ensure you've pulled the model with `docker model pull`
- **Connection errors**: Verify Docker network settings and that Model Runner is running
- **Streaming issues**: Check CORS settings in the backend code
- **Metrics not showing**: Verify that Prometheus can reach the backend metrics endpoint
- **Redis connection failed**: Check Redis container status and network connectivity
- **llama.cpp metrics missing**: Confirm that your model is indeed a llama.cpp model
- **Grafana dashboards empty**: Ensure Prometheus is collecting metrics and data source is configured correctly

### Redis-Specific Troubleshooting

- **Redis Insight not accessible**: Check if port 8001 is available and Redis container is running
- **Token analytics not working**: Verify Redis TimeSeries module is loaded and service dependencies are met
- **Performance degradation**: Monitor Redis memory usage and consider adjusting configuration
- **Data not persisting**: Check Redis volume mounts and persistence configuration

### Health Checks

Monitor service health using:
```bash
# Check all container status
docker-compose -f compose.enhanced.yaml ps

# View specific service logs
docker-compose -f compose.enhanced.yaml logs redis
docker-compose -f compose.enhanced.yaml logs grafana
docker-compose -f compose.enhanced.yaml logs token-analytics
```

## 📊 Performance Optimization

### Redis Optimization
- **Memory Management**: Configure `maxmemory` and eviction policies
- **Persistence**: Balance between RDB and AOF based on use case
- **Networking**: Use Redis clustering for high availability
- **Monitoring**: Set up alerts for memory usage and connection limits

### Model Performance
- **Thread Configuration**: Optimize thread count based on CPU cores
- **Memory Settings**: Configure context window based on available RAM
- **Batch Processing**: Adjust batch size for optimal throughput

## 🔄 Migration from Basic Setup

If upgrading from a previous version:

1. **Backup existing data** (if any)
2. **Stop current services**: `docker-compose down`
3. **Use new compose file**: `docker-compose -f compose.enhanced.yaml up -d --build`
4. **Verify all services**: Check health endpoints and Grafana dashboards
5. **Import existing data** into Redis if needed

## 📜 License

MIT

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- Docker Model Runner team for local LLM capabilities
- Redis Stack for comprehensive data management
- Grafana and Prometheus communities for observability tools
- OpenTelemetry project for distributed tracing standards
