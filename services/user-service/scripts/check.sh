#!/usr/bin/env bash

set -e

echo "🧹 Formatting..."
go fmt ./...

echo "🔍 Running Go Vet..."
go vet ./...

echo "🧪 Running Tests..."
go test ./...

echo "✅ All checks passed!"