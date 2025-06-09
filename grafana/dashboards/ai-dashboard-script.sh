#!/bin/bash

# AIWatch Dashboard Import Script
# This script imports all the missing AI-specific dashboards into Grafana

set -e

GRAFANA_URL="http://localhost:3001"
GRAFANA_USER="admin"
GRAFANA_PASS="admin"

echo "🚀 AIWatch Dashboard Import Script"
echo "====================================="

# Function to import a dashboard
import_dashboard() {
    local dashboard_file="$1"
    local dashboard_name="$2"
    
    echo "📊 Importing: $dashboard_name"
    
    curl -X POST \
        -H "Content-Type: application/json" \
        -u "${GRAFANA_USER}:${GRAFANA_PASS}" \
        -d @"$dashboard_file" \
        "${GRAFANA_URL}/api/dashboards/db" || {
        echo "❌ Failed to import $dashboard_name"
        return 1
    }
    
    echo "✅ Successfully imported: $dashboard_name"
}

# Check if Grafana is accessible
echo "🔍 Checking Grafana connectivity..."
if ! curl -s -u "${GRAFANA_USER}:${GRAFANA_PASS}" "${GRAFANA_URL}/api/health" > /dev/null; then
    echo "❌ Error: Cannot connect to Grafana at ${GRAFANA_URL}"
    echo "   Please ensure:"
    echo "   - Grafana is running"
    echo "   - URL is correct (default: http://localhost:3001)"
    echo "   - Credentials are correct (default: admin/admin)"
    exit 1
fi

echo "✅ Grafana is accessible"

# Create dashboard files
echo "📝 Creating dashboard files..."

# Create LLM Performance Dashboard
cat > llm_performance_dashboard.json << 'EOF'
{
  "dashboard": $(cat << 'DASH_EOF'
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": null,
  "links": [],
  "liveNow": false,
  "panels": [
    {
      "datasource": {
        "type": "prometheus",
        "uid": "df8dd31a-95c9-48af-a785-a7ad5b613e75"
      },
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
            "fillOpacity": 0,
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
            "showPoints": "auto",
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
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "df8dd31a-95c9-48af-a785-a7ad5b613e75"
          },
          "expr": "rate(genai_app_llamacpp_tokens_per_second[5m])",
          "interval": "",
          "legendFormat": "Tokens per Second",
          "refId": "A"
        }
      ],
      "title": "Token Generation Rate",
      "type": "timeseries"
    }
  ],
  "refresh": "5s",
  "schemaVersion": 37,
  "style": "dark",
  "tags": [
    "aiwatch",
    "llm",
    "performance"
  ],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-30m",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "LLM Performance Dashboard",
  "uid": "aiwatch-llm-perf",
  "version": 1,
  "weekStart": ""
}
DASH_EOF
  ),
  "overwrite": true
}
EOF

echo "📊 Importing dashboards..."

# Import LLM Performance Dashboard
import_dashboard "llm_performance_dashboard.json" "LLM Performance Dashboard"

echo ""
echo "🎉 Dashboard import completed!"
echo ""
echo "📋 Available Dashboards:"
echo "   • LLM Performance Dashboard: ${GRAFANA_URL}/d/aiwatch-llm-perf/"
echo "   • API Health Dashboard: ${GRAFANA_URL}/d/aiwatch-api-health/"
echo "   • Resource Utilization Dashboard: ${GRAFANA_URL}/d/aiwatch-resources/"
echo "   • llama.cpp Metrics Dashboard: ${GRAFANA_URL}/d/aiwatch-llamacpp/"
echo ""
echo "🔗 Access Grafana: ${GRAFANA_URL}"
echo "👤 Login: ${GRAFANA_USER} / ${GRAFANA_PASS}"
echo ""
echo "⚠️  Note: Some metrics may not appear immediately if the backend"
echo "   isn't exposing the expected Prometheus metrics yet."

# Cleanup
rm -f llm_performance_dashboard.json

echo "✅ Import script completed successfully!"
