#!/bin/bash

# 检查是否提供了环境参数
if [ -z "$1" ]; then
  echo "Usage: $0 <environment>"
  echo "Available environments: dev, prod"
  exit 1
fi

# 设置环境变量
export APP_ENV=$1

echo "Starting application in $APP_ENV environment..."

# 如果是生产环境，可能需要设置数据库密码
if [ "$APP_ENV" == "prod" ]; then
  # 检查是否设置了数据库密码环境变量
  if [ -z "$DB_PASSWORD" ]; then
    echo "Warning: DB_PASSWORD not set. Please set it for production environment."
    echo "Example: DB_PASSWORD=your_secure_password $0 prod"
    exit 1
  fi
fi

# 进入项目根目录
cd "$(dirname "$0")/.."

# 构建应用
echo "Building application..."
go build -o bin/textile-admin cmd/api/main.go

# 检查是否构建成功
if [ $? -ne 0 ]; then
  echo "Build failed!"
  exit 1
fi

# 运行应用
echo "Running application..."
./bin/textile-admin 