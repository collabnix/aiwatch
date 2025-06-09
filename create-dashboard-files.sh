#!/bin/bash

echo "=== Creating AIWatch Dashboard JSON Files ==="
echo

# Check if directories exist, create if they don't
echo "📁 Ensuring directory structure exists..."
mkdir -p grafana/dashboards/llm
mkdir -p grafana/dashboards/infrastructure

# Create LLaMA.cpp Performance Dashboard
echo "📊 Creating LLaMA.cpp Performance Dashboard..."
cat > grafana/dashboards/llm/llamacpp-performance.json << 'EOF'
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "description": "AIWatch - LLaMA.cpp Model Performance Monitoring",
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "panels": [
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 0
      },
      "id": 1,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "rate(genai_app_llamacpp_tokens_per_second[5m])",
          "interval": "",
          "legendFormat": "Tokens/sec",
          "refId": "A"
        }
      ],
      "title": "Token Generation Speed",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "yellow",
                "value": 0.5
              },
              {
                "color": "red",
                "value": 1
              }
            ]
          },
          "unit": "s"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 0
      },
      "id": 2,
      "options": {
        "orientation": "auto",
        "reduceOptions": {
          "values": false,
          "calcs": [
            "lastNotNull"
          ],
          "fields": ""
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "text": {}
      },
      "pluginVersion": "8.0.0",
      "targets": [
        {
          "expr": "genai_app_llamacpp_prompt_eval_seconds",
          "interval": "",
          "legendFormat": "Prompt Evaluation Time",
          "refId": "A"
        }
      ],
      "title": "Prompt Evaluation Time",
      "type": "gauge"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "bytes"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 8
      },
      "id": 3,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "genai_app_llamacpp_memory_per_token_bytes",
          "interval": "",
          "legendFormat": "Memory per Token",
          "refId": "A"
        }
      ],
      "title": "Memory Usage per Token",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "yellow",
                "value": 2048
              },
              {
                "color": "red",
                "value": 4096
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 6,
        "x": 12,
        "y": 8
      },
      "id": 4,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "auto",
        "reduceOptions": {
          "values": false,
          "calcs": [
            "lastNotNull"
          ],
          "fields": ""
        },
        "text": {},
        "textMode": "auto"
      },
      "pluginVersion": "8.0.0",
      "targets": [
        {
          "expr": "genai_app_llamacpp_context_size",
          "interval": "",
          "legendFormat": "Context Size",
          "refId": "A"
        }
      ],
      "title": "Context Window Size",
      "type": "stat"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "yellow",
                "value": 4
              },
              {
                "color": "red",
                "value": 8
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 6,
        "x": 18,
        "y": 8
      },
      "id": 5,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "auto",
        "reduceOptions": {
          "values": false,
          "calcs": [
            "lastNotNull"
          ],
          "fields": ""
        },
        "text": {},
        "textMode": "auto"
      },
      "pluginVersion": "8.0.0",
      "targets": [
        {
          "expr": "genai_app_llamacpp_threads_used",
          "interval": "",
          "legendFormat": "Threads Used",
          "refId": "A"
        }
      ],
      "title": "Thread Utilization",
      "type": "stat"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 24,
        "x": 0,
        "y": 16
      },
      "id": 6,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "genai_app_llamacpp_batch_size",
          "interval": "",
          "legendFormat": "Batch Size",
          "refId": "A"
        }
      ],
      "title": "Batch Size",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 27,
  "style": "dark",
  "tags": [
    "aiwatch",
    "llama.cpp",
    "llm",
    "performance"
  ],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "AIWatch - LLaMA.cpp Performance",
  "uid": "aiwatch-llamacpp",
  "version": 1
}
EOF

# Create API Health Dashboard
echo "📊 Creating API Health Dashboard..."
cat > grafana/dashboards/api-health.json << 'EOF'
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "description": "AIWatch - API Health and Performance Monitoring",
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "panels": [
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [
            {
              "options": {
                "0": {
                  "color": "red",
                  "index": 0,
                  "text": "Down"
                },
                "1": {
                  "color": "green",
                  "index": 1,
                  "text": "Up"
                }
              },
              "type": "value"
            }
          ],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "red",
                "value": null
              },
              {
                "color": "green",
                "value": 1
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 6,
        "x": 0,
        "y": 0
      },
      "id": 1,
      "options": {
        "colorMode": "background",
        "graphMode": "none",
        "justifyMode": "auto",
        "orientation": "auto",
        "reduceOptions": {
          "values": false,
          "calcs": [
            "lastNotNull"
          ],
          "fields": ""
        },
        "text": {},
        "textMode": "auto"
      },
      "pluginVersion": "8.0.0",
      "targets": [
        {
          "expr": "up{job=\"aiwatch-backend\"}",
          "interval": "",
          "legendFormat": "Backend Status",
          "refId": "A"
        }
      ],
      "title": "API Status",
      "type": "stat"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "reqps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 18,
        "x": 6,
        "y": 0
      },
      "id": 2,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "rate(http_requests_total{job=\"aiwatch-backend\"}[5m])",
          "interval": "",
          "legendFormat": "{{method}} {{handler}}",
          "refId": "A"
        }
      ],
      "title": "Request Rate",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "s"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 8
      },
      "id": 3,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{job=\"aiwatch-backend\"}[5m]))",
          "interval": "",
          "legendFormat": "95th percentile",
          "refId": "A"
        },
        {
          "expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket{job=\"aiwatch-backend\"}[5m]))",
          "interval": "",
          "legendFormat": "50th percentile",
          "refId": "B"
        }
      ],
      "title": "Response Time",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "reqps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 8
      },
      "id": 4,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "rate(http_requests_total{job=\"aiwatch-backend\",code!~\"2..\"}[5m])",
          "interval": "",
          "legendFormat": "Error Rate {{code}}",
          "refId": "A"
        }
      ],
      "title": "Error Rate",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 16
      },
      "id": 5,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "go_memstats_alloc_bytes{job=\"aiwatch-backend\"}",
          "interval": "",
          "legendFormat": "Memory Usage",
          "refId": "A"
        }
      ],
      "title": "Memory Usage",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "vis": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 16
      },
      "id": 6,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom"
        },
        "tooltip": {
          "mode": "single"
        }
      },
      "targets": [
        {
          "expr": "go_goroutines{job=\"aiwatch-backend\"}",
          "interval": "",
          "legendFormat": "Goroutines",
          "refId": "A"
        }
      ],
      "title": "Active Goroutines",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 27,
  "style": "dark",
  "tags": [
    "aiwatch",
    "api",
    "health",
    "backend"
  ],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "AIWatch - API Health",
  "uid": "aiwatch-api-health",
  "version": 1
}
EOF

# Create Resource Utilization Dashboard
echo "📊 Creating Resource Utilization Dashboard..."
cat > grafana/dashboards/infrastructure/resource-utilization.json << 'EOF'
{
  "annotations": {
    "list": []
  },
  "description": "AIWatch - System Resource Monitoring",
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "panels": [
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {"mode": "palette-classic"},
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {"legend": false, "tooltip": false, "vis": false},
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {"type": "linear"},
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {"group": "A", "mode": "none"},
            "thresholdsStyle": {"mode": "off"}
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "red", "value": 80}
            ]
          },
          "unit": "percent"
        },
        "overrides": []
      },
      "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0},
      "id": 1,
      "options": {
        "legend": {"calcs": [], "displayMode": "list", "placement": "bottom"},
        "tooltip": {"mode": "single"}
      },
      "targets": [
        {
          "expr": "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)",
          "interval": "",
          "legendFormat": "CPU Usage",
          "refId": "A"
        }
      ],
      "title": "CPU Usage",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {"mode": "palette-classic"},
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {"legend": false, "tooltip": false, "vis": false},
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {"type": "linear"},
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {"group": "A", "mode": "none"},
            "thresholdsStyle": {"mode": "off"}
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "red", "value": 80}
            ]
          },
          "unit": "percent"
        },
        "overrides": []
      },
      "gridPos": {"h": 8, "w": 12, "x": 12, "y": 0},
      "id": 2,
      "options": {
        "legend": {"calcs": [], "displayMode": "list", "placement": "bottom"},
        "tooltip": {"mode": "single"}
      },
      "targets": [
        {
          "expr": "(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100",
          "interval": "",
          "legendFormat": "Memory Usage",
          "refId": "A"
        }
      ],
      "title": "Memory Usage",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {"mode": "palette-classic"},
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {"legend": false, "tooltip": false, "vis": false},
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {"type": "linear"},
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {"group": "A", "mode": "none"},
            "thresholdsStyle": {"mode": "off"}
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "red", "value": 80}
            ]
          },
          "unit": "bytes"
        },
        "overrides": []
      },
      "gridPos": {"h": 8, "w": 12, "x": 0, "y": 8},
      "id": 3,
      "options": {
        "legend": {"calcs": [], "displayMode": "list", "placement": "bottom"},
        "tooltip": {"mode": "single"}
      },
      "targets": [
        {
          "expr": "go_memstats_alloc_bytes",
          "interval": "",
          "legendFormat": "Go Memory Allocated",
          "refId": "A"
        }
      ],
      "title": "Application Memory",
      "type": "timeseries"
    },
    {
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "color": {"mode": "palette-classic"},
          "custom": {
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 10,
            "gradientMode": "none",
            "hideFrom": {"legend": false, "tooltip": false, "vis": false},
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {"type": "linear"},
            "showPoints": "never",
            "spanNulls": false,
            "stacking": {"group": "A", "mode": "none"},
            "thresholdsStyle": {"mode": "off"}
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {"color": "green", "value": null},
              {"color": "red", "value": 80}
            ]
          },
          "unit": "short"
        },
        "overrides": []
      },
      "gridPos": {"h": 8, "w": 12, "x": 12, "y": 8},
      "id": 4,
      "options": {
        "legend": {"calcs": [], "displayMode": "list", "placement": "bottom"},
        "tooltip": {"mode": "single"}
      },
      "targets": [
        {
          "expr": "go_goroutines",
          "interval": "",
          "legendFormat": "Active Goroutines",
          "refId": "A"
        }
      ],
      "title": "Go Goroutines",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 27,
  "style": "dark",
  "tags": ["aiwatch", "infrastructure", "resources"],
  "templating": {"list": []},
  "time": {"from": "now-1h", "to": "now"},
  "timepicker": {},
  "timezone": "",
  "title": "AIWatch - Resource Utilization",
  "uid": "aiwatch-resources",
  "version": 1
}
EOF

echo "✅ All dashboard JSON files created successfully!"
echo
echo "📁 Dashboard files created:"
echo "├── grafana/dashboards/llm/llamacpp-performance.json"
echo "├── grafana/dashboards/api-health.json"
echo "└── grafana/dashboards/infrastructure/resource-utilization.json"
echo
echo "🎯 Next steps:"
echo "1. Update your docker-compose.yml with proper Grafana volumes"
echo "2. Test: docker compose down && docker compose up -d --build"
echo "3. Visit: http://localhost:3001 (admin/admin)"
echo
echo "🔍 Run troubleshooting if needed: ./scripts/troubleshoot-dashboards.sh"
