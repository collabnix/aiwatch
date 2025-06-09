# Redis Observability Enhancement Guide

This branch contains comprehensive Redis observability enhancements for the AIWatch project, providing real-time analytics, monitoring, and alerting capabilities.

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Frontend UI    │◄──►│  Backend API    │◄──►│ Model Runner    │
│  (React + TS)   │    │     (Go)        │    │ (Llama 3.2)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Analytics UI    │    │     Redis       │    │ Redis TimeSeries│
│ Components      │    │   (Primary)     │    │   (Enhanced)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Token Analytics │    │ Redis Exporter  │    │ Grafana + Alert │
│    Service      │    │  (Prometheus)   │    │   Dashboard     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 What's New in Redis Observability

### 1. **Comprehensive Redis Infrastructure**
- **Redis Primary**: Optimized for token metrics storage
- **Redis TimeSeries**: Advanced time-series data handling
- **Redis Insight**: Web-based Redis management interface
- **Redis Exporters**: Prometheus metrics collection

### 2. **Token Analytics Service** (Port: 8081)
- Real-time user activity tracking
- Model usage analytics
- Token processing metrics
- Performance monitoring

### 3. **Redis TimeSeries Service** (Port: 8082)
- Advanced time-series data storage
- Efficient data retention policies
- High-performance querying
- Automated metrics collection

### 4. **Enhanced Observability Stack**
- **Grafana Dashboards**: Visual analytics and monitoring
- **Prometheus Alerting**: Intelligent threshold-based alerts
- **Jaeger Tracing**: Distributed request tracing
- **Alertmanager**: Alert routing and management

### 5. **Frontend Analytics Integration**
- Real-time analytics dashboard
- Interactive charts and visualizations
- User behavior insights
- Model performance metrics

## 📊 Available Dashboards & Metrics

### Redis Performance Metrics
- Memory usage and fragmentation
- Commands per second
- Connection count
- Slow log monitoring
- Key statistics

### Token Analytics
- Active users (5m, 1h, 24h windows)
- Token processing rates
- Model usage distribution
- Response time percentiles
- Error rates and patterns

### System Metrics
- CPU and memory utilization
- Disk space monitoring
- Network performance
- Container health status

## 🔧 Quick Start Guide

### Prerequisites
```bash
# Ensure Docker and Docker Compose are installed
docker --version
docker-compose --version

# Pull the required model
docker model pull ai/llama3.2:1B-Q8_0
```

### Start the Enhanced Stack
```bash
# Basic stack (Redis + Analytics)
docker-compose -f compose.enhanced.yaml up -d

# With monitoring (includes Alertmanager)
docker-compose -f compose.enhanced.yaml --profile monitoring up -d

# With debugging tools (includes Redis Insight)
docker-compose -f compose.enhanced.yaml --profile debug up -d

# Full stack with all components
docker-compose -f compose.enhanced.yaml --profile monitoring --profile debug up -d
```

### Environment Configuration
Create a `.env` file with:
```bash
# Redis Configuration
REDIS_PASSWORD=your_secure_password
LLM_MODEL_NAME=ai/llama3.2:1B-Q8_0

# Optional: Custom ports
ANALYTICS_PORT=8081
TIMESERIES_PORT=8082
```

## 📱 Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **Frontend** | http://localhost:3000 | Main application interface |
| **Analytics Dashboard** | http://localhost:3000/analytics | Real-time analytics |
| **Backend API** | http://localhost:8080 | Main API endpoints |
| **Token Analytics** | http://localhost:8081 | Analytics service API |
| **TimeSeries API** | http://localhost:8082 | Time-series data API |
| **Grafana** | http://localhost:3001 | Dashboards (admin/admin) |
| **Prometheus** | http://localhost:9091 | Metrics & alerting |
| **Jaeger** | http://localhost:16686 | Distributed tracing |
| **Redis Insight** | http://localhost:8001 | Redis management |
| **Alertmanager** | http://localhost:9093 | Alert management |

## 🎯 Key API Endpoints

### Token Analytics Service
```bash
# Get comprehensive analytics
GET http://localhost:8081/analytics

# Health check
GET http://localhost:8081/health

# Prometheus metrics
GET http://localhost:8081/metrics
```

### Redis TimeSeries Service
```bash
# Query time-series data
POST http://localhost:8082/query
Content-Type: application/json
{
  "key": "metrics:tokens:input_rate",
  "start_time": 1640995200000,
  "end_time": 1640998800000,
  "aggregation": "avg",
  "bucket_duration": 60000
}

# Get latest value
GET http://localhost:8082/latest?key=metrics:users:active_5m

# Multi-query
POST http://localhost:8082/multi-query
Content-Type: application/json
[
  {
    "key": "metrics:tokens:input_rate",
    "start_time": 1640995200000,
    "end_time": 1640998800000
  },
  {
    "key": "metrics:tokens:output_rate",
    "start_time": 1640995200000,
    "end_time": 1640998800000
  }
]
```

## 🔔 Alerting Configuration

### Prometheus Alerts
The system includes predefined alerts for:
- Redis memory usage > 90%
- High error rates
- Model response time > 30s
- Low user activity
- System resource exhaustion

### Alertmanager Integration
Configure notifications in `alertmanager/alertmanager.yml`:
```yaml
global:
  smtp_smarthost: 'localhost:587'
  smtp_from: 'alerts@yourcompany.com'

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
- name: 'web.hook'
  email_configs:
  - to: 'admin@yourcompany.com'
    subject: 'AIWatch Alert: {{ .GroupLabels.alertname }}'
    body: |
      {{ range .Alerts }}
      Alert: {{ .Annotations.summary }}
      Description: {{ .Annotations.description }}
      {{ end }}
```

## 📈 Performance Optimization

### Redis Configuration
The Redis configuration is optimized for:
- **Memory efficiency**: Hash/ZipList optimizations
- **Persistence**: AOF + RDB snapshots
- **Performance**: Appropriate timeout and backlog settings
- **Monitoring**: Slow log and latency tracking

### TimeSeries Optimization
- **Retention policies**: 24h for real-time, 7d for historical
- **Compression**: Automatic data compression
- **Aggregation**: Efficient downsampling
- **Labels**: Structured metadata for querying

## 🧪 Testing & Validation

### Health Checks
```bash
# Check all services
docker-compose -f compose.enhanced.yaml ps

# Verify Redis connectivity
docker-compose -f compose.enhanced.yaml exec redis redis-cli ping

# Check analytics service
curl http://localhost:8081/health

# Verify TimeSeries service
curl http://localhost:8082/health
```

### Load Testing
```bash
# Simple load test for analytics endpoint
for i in {1..100}; do
  curl -s http://localhost:8081/analytics > /dev/null &
done
wait

# Monitor Redis performance
docker-compose -f compose.enhanced.yaml exec redis redis-cli --latency-history
```

## 🔍 Troubleshooting

### Common Issues

1. **Redis Connection Failed**
   ```bash
   # Check Redis logs
   docker-compose -f compose.enhanced.yaml logs redis
   
   # Verify network connectivity
   docker-compose -f compose.enhanced.yaml exec backend ping redis
   ```

2. **Analytics Service Not Responding**
   ```bash
   # Check service logs
   docker-compose -f compose.enhanced.yaml logs token-analytics
   
   # Restart service
   docker-compose -f compose.enhanced.yaml restart token-analytics
   ```

3. **Grafana Dashboard Issues**
   ```bash
   # Check Grafana logs
   docker-compose -f compose.enhanced.yaml logs grafana
   
   # Verify Prometheus data source
   curl http://localhost:9091/api/v1/targets
   ```

### Performance Issues
- **High Memory Usage**: Adjust Redis `maxmemory` settings
- **Slow Queries**: Enable Redis slow log monitoring
- **Connection Limits**: Increase `maxclients` in Redis config
- **Disk Space**: Monitor Docker volume usage

## 🔮 Future Enhancements

### Planned Features
- [ ] **Machine Learning Integration**: Predictive analytics for usage patterns
- [ ] **Multi-tenant Support**: Organization-level analytics separation
- [ ] **API Rate Limiting**: Redis-based intelligent throttling
- [ ] **Advanced Dashboards**: Custom user-configurable dashboards
- [ ] **Data Export**: CSV/JSON export functionality
- [ ] **Webhook Notifications**: Real-time event notifications

### Scaling Considerations
- **Redis Cluster**: Horizontal scaling for high-volume deployments
- **Load Balancing**: Multiple analytics service instances
- **Data Archival**: Long-term storage strategies
- **Cross-Region Replication**: Global deployment support

## 🤝 Contributing

To contribute to the Redis Observability enhancements:

1. Create a feature branch from `Redis_Observability`
2. Implement your enhancements
3. Add appropriate tests and documentation
4. Submit a pull request with detailed description

### Development Setup
```bash
# Clone and switch to branch
git clone https://github.com/collabnix/aiwatch.git
cd aiwatch
git checkout Redis_Observability

# Start development environment
docker-compose -f compose.enhanced.yaml up --build

# Run tests
cd tests
go test -v
```

## 📚 Additional Resources

- [Redis TimeSeries Documentation](https://redis.io/docs/data-types/timeseries/)
- [Prometheus Alerting Rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/best-practices/)
- [OpenTelemetry Tracing](https://opentelemetry.io/docs/instrumentation/go/)

---

**Note**: This Redis Observability enhancement provides a solid foundation for monitoring and analytics. The implementation is production-ready and includes comprehensive error handling, health checks, and performance optimizations.
