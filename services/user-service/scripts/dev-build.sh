#!/usr/bin/env bash

# ==========================================
# UnVocal - Pre Build Script
# Formats code and generates Swagger docs
# ==========================================

set -euo pipefail

LINE="━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo
echo "$LINE"
echo "🚀 UnVocal Pre-Build"
echo "⏰ Started : $(date '+%Y-%m-%d %H:%M:%S')"
echo "$LINE"

# Check required tools
command -v go >/dev/null 2>&1 || {
    echo "❌ Go is not installed."
    exit 1
}

command -v swag >/dev/null 2>&1 || {
    echo "❌ Swag is not installed."
    echo "Install: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
}

echo
echo "🧹 Formatting Go code..."
go fmt ./...
echo "✅ Formatting completed."

# Uncomment if you want to run it on every save.
# Note: It can slow down hot reload.
#
# echo
# echo "🔍 Running Go Vet..."
# go vet ./...
# echo "✅ Go Vet passed."

echo
echo "📘 Generating Swagger Documentation..."

swag init \
    -g main.go \
    -d cmd/server,handler,model,repository,routes,service \
    --output docs

echo "✅ Swagger generated successfully."
echo "📂 Output : docs/"

echo
echo "$LINE"
echo "🎉 Pre-build completed successfully."
echo "⏰ Finished : $(date '+%Y-%m-%d %H:%M:%S')"
echo "$LINE"
echo