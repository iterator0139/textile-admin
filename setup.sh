#!/bin/bash

# Textile Admin Setup Script

echo "Setting up Textile Admin project..."

# Create uploads directory if it doesn't exist
if [ ! -d "uploads" ]; then
  mkdir -p uploads
  echo "Created uploads directory"
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
  echo "Error: Go is not installed or not in PATH"
  echo "Please install Go from https://golang.org/dl/"
  exit 1
fi

# Install dependencies
echo "Installing dependencies..."
go mod download
go mod tidy

# Build the application
echo "Building application..."
go build -o textile-admin ./cmd/api

echo "Setup complete! You can now run the application with ./textile-admin"
echo "Make sure to configure your database settings in environment variables or use the defaults."
echo "For more information, see the README.md file." 