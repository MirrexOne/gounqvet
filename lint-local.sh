#!/bin/bash

# Script for running golangci-lint locally with v2.x configuration

# Create temporary config with version field for local v2.x
cat > .golangci.local.yml << 'EOF'
version: "2"

EOF

# Append the rest of the config
cat .golangci.yml >> .golangci.local.yml

# Run golangci-lint with local config
golangci-lint run --config=.golangci.local.yml "$@"

# Clean up
rm -f .golangci.local.yml
