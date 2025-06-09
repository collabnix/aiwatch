#!/bin/bash

echo "=== Creating Pull Request for AIWatch Dashboards ==="
echo

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Not in a git repository. Please run this from your aiwatch project root."
    exit 1
fi

# Check if GitHub CLI is available
if command -v gh &> /dev/null; then
    echo "✅ GitHub CLI found - will create PR automatically"
    GH_CLI_AVAILABLE=true
else
    echo "ℹ️  GitHub CLI not found - will provide manual instructions"
    GH_CLI_AVAILABLE=false
fi

# Get current branch
current_branch=$(git branch --show-current)
echo "📍 Current branch: $current_branch"

# Create feature branch if not already on one
if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
    echo "🔀 Creating feature branch for dashboard changes..."
    git checkout -b feature/add-grafana-dashboards
    echo "✅ Switched to feature/add-grafana-dashboards branch"
else
    echo "✅ Already on feature branch: $current_branch"
fi

echo
echo "📁 Checking files to commit..."

# Run our setup scripts if files don't exist
if [ ! -f "grafana/provisioning/datasources/prometheus.yml" ]; then
    echo "🔧 Creating Grafana configuration files..."
    ./git-setup.sh
fi

if [ ! -f "grafana/dashboards/llm/llamacpp-performance.json" ]; then
    echo "📊 Creating dashboard JSON files..."
    ./create-dashboard-files.sh
fi

# Add all files
echo "📝 Adding files to git..."
git add .

# Show what's being committed
echo "📋 Files to be committed:"
git status --porcelain

echo
# Create commit
echo "💬 Creating commit..."
git commit -m "Add predefined Grafana dashboards and monitoring tools

🎯 Features Added:
- Grafana provisioning configuration for automatic dashboard loading
- LLaMA.cpp performance dashboard with token generation metrics
- API health monitoring dashboard with request/response metrics
- Resource utilization dashboard for system monitoring
- Troubleshooting and setup scripts
- Comprehensive documentation

📊 Dashboards Include:
- Token generation speed and context window monitoring
- Prompt evaluation time and memory usage per token
- API request rates, response times, and error tracking
- CPU, memory, and goroutine monitoring
- Thread utilization and batch size metrics

🔧 Tools Added:
- Automated troubleshooting script
- Dashboard setup automation
- Docker Compose configuration updates
- Prometheus configuration for metrics collection

Fixes: Missing predefined dashboards for Model Runner, llama.cpp and token metrics
Resolves: Grafana showing empty state with no pre-configured monitoring

Co-authored-by: Claude <claude@anthropic.com>"

# Push to GitHub
echo "🚀 Pushing to GitHub..."
git push origin $(git branch --show-current)

echo
if [ "$GH_CLI_AVAILABLE" = true ]; then
    echo "🎯 Creating Pull Request with GitHub CLI..."
    
    # Create PR with GitHub CLI
    gh pr create \
        --title "Add predefined Grafana dashboards for AIWatch monitoring" \
        --body "## 🎯 Overview

This PR adds the missing predefined Grafana dashboards for AIWatch, providing comprehensive monitoring for Model Runner, llama.cpp, and token metrics.

## 📊 What's Added

### **Predefined Dashboards**
- **LLaMA.cpp Performance**: Token generation speed, context size, memory usage, thread utilization
- **API Health Monitoring**: Request rates, response times, error tracking, memory usage
- **Resource Utilization**: CPU, memory, Go goroutines, system metrics

### **Grafana Provisioning**
- Automatic dashboard loading via provisioning
- Prometheus datasource auto-configuration
- Organized folder structure (AIWatch, LLM Monitoring, Infrastructure)

### **Monitoring Metrics**
- \`genai_app_llamacpp_tokens_per_second\` - Token generation speed
- \`genai_app_llamacpp_context_size\` - Context window monitoring
- \`genai_app_llamacpp_prompt_eval_seconds\` - Prompt evaluation time
- \`genai_app_llamacpp_memory_per_token_bytes\` - Memory efficiency
- \`genai_app_llamacpp_threads_used\` - Thread utilization
- \`genai_app_llamacpp_batch_size\` - Batch processing metrics

### **Tools & Scripts**
- \`scripts/troubleshoot-dashboards.sh\` - Automated troubleshooting
- \`create-dashboard-files.sh\` - Dashboard setup automation
- Comprehensive setup documentation

## 🔧 Technical Changes

### **File Structure**
\`\`\`
grafana/
├── provisioning/
│   ├── datasources/prometheus.yml     # Auto-configure Prometheus
│   └── dashboards/dashboards.yml      # Dashboard provisioning
└── dashboards/
    ├── llm/llamacpp-performance.json  # LLaMA.cpp metrics
    ├── api-health.json                # API monitoring
    └── infrastructure/resource-utilization.json # System resources
\`\`\`

### **Docker Compose Updates**
- Added proper volume mounts for Grafana provisioning
- Configured Grafana plugins (redis-app, grafana-piechart-panel)
- Set up automatic dashboard loading

## 🚀 How to Test

1. **Apply changes**:
   \`\`\`bash
   docker compose down && docker compose up -d --build
   \`\`\`

2. **Access Grafana**: http://localhost:3001 (admin/admin)

3. **Verify dashboards** appear in organized folders:
   - AIWatch folder: General monitoring
   - LLM Monitoring folder: LLaMA.cpp performance
   - Infrastructure folder: System resources

4. **Check metrics**: Visit Prometheus at http://localhost:9091 and search for \`genai_app_\` metrics

## 🔍 Before/After

**Before**: Empty Grafana with manual datasource setup required
**After**: Pre-configured dashboards with comprehensive LLM and API monitoring

## 📚 Documentation

- Added \`DASHBOARD_SETUP.md\` with complete setup instructions
- Troubleshooting guide for common issues
- Metric definitions and explanations

## ✅ Resolves

- Missing predefined dashboards for Model Runner
- No llama.cpp specific monitoring
- Lack of token metrics visualization
- Manual Grafana configuration requirements

Fixes #[issue-number] (if applicable)" \
        --assignee @me \
        --label "enhancement,monitoring,grafana"
    
    echo "✅ Pull Request created successfully!"
    echo
    echo "🔗 View your PR:"
    gh pr view --web
    
else
    echo "📋 Manual PR Creation Instructions:"
    echo
    echo "1. Go to your GitHub repository"
    echo "2. Click 'Compare & pull request' button"
    echo "3. Use this title:"
    echo "   Add predefined Grafana dashboards for AIWatch monitoring"
    echo
    echo "4. Copy this description:"
    echo
    cat << 'EOF'
## 🎯 Overview

This PR adds the missing predefined Grafana dashboards for AIWatch, providing comprehensive monitoring for Model Runner, llama.cpp, and token metrics.

## 📊 What's Added

### **Predefined Dashboards**
- **LLaMA.cpp Performance**: Token generation speed, context size, memory usage, thread utilization
- **API Health Monitoring**: Request rates, response times, error tracking, memory usage
- **Resource Utilization**: CPU, memory, Go goroutines, system metrics

### **Grafana Provisioning**
- Automatic dashboard loading via provisioning
- Prometheus datasource auto-configuration
- Organized folder structure (AIWatch, LLM Monitoring, Infrastructure)

### **Monitoring Metrics**
- `genai_app_llamacpp_tokens_per_second` - Token generation speed
- `genai_app_llamacpp_context_size` - Context window monitoring
- `genai_app_llamacpp_prompt_eval_seconds` - Prompt evaluation time
- `genai_app_llamacpp_memory_per_token_bytes` - Memory efficiency
- `genai_app_llamacpp_threads_used` - Thread utilization
- `genai_app_llamacpp_batch_size` - Batch processing metrics

## 🚀 How to Test

1. **Apply changes**: `docker compose down && docker compose up -d --build`
2. **Access Grafana**: http://localhost:3001 (admin/admin)
3. **Verify dashboards** appear in organized folders

## ✅ Resolves

- Missing predefined dashboards for Model Runner
- No llama.cpp specific monitoring
- Lack of token metrics visualization
EOF
    echo
    echo "5. Create the pull request"
fi

echo
echo "🎉 Next Steps:"
echo "1. Wait for PR review and approval"
echo "2. Once merged, test the dashboards:"
echo "   docker compose down && docker compose up -d --build"
echo "3. Access Grafana at http://localhost:3001"
echo "4. Check for the new dashboard folders and metrics"
echo
echo "✅ PR creation process complete!"
