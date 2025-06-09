#!/bin/bash

echo "=== Quick Light Mode Update for AIWatch ==="
echo

# Quick one-liner updates for existing installations

echo "🎨 Updating all dashboards to light mode..."

# Update dashboard themes
find grafana/dashboards -name "*.json" -exec sed -i 's/"style": "dark"/"style": "light"/g' {} \;

echo "✅ Updated dashboard themes"

echo
echo "🔧 Add these environment variables to your Grafana service in docker-compose.yml:"
echo
echo "    environment:"
echo "      - GF_USERS_DEFAULT_THEME=light"
echo "      - GF_UI_DEFAULT_THEME=light"
echo "      - GF_AUTH_ANONYMOUS_THEME=light"

echo
echo "🚀 Then restart:"
echo "docker compose down && docker compose up -d"

echo
echo "✅ Light mode update complete!"
