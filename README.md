# Textile Admin - Reading Tasks API

A Golang web service for managing reading tasks and file uploads.

## Features

- Upload files and create reading tasks
- Query task details by ID
- List all reading tasks for a user
- Download files associated with reading tasks

## Project Structure

```
textile-admin/
├── cmd/
│   └── api/
│       └── main.go         # Main application entry point
├── internal/
│   ├── config/             # Application configuration
│   ├── domain/entity/      # Domain entities
│   ├── repository/         # Database access layer
│   ├── service/            # Business logic layer
│   └── handler/            # HTTP request handlers
├── pkg/
│   ├── db/                 # Database utilities
│   └── response/           # API response utilities
├── scripts/
│   └── schema.sql          # Database schema
├── uploads/                # File storage directory
├── .gitignore              # Git ignore file
├── CONTRIBUTING.md         # Contribution guidelines
├── LICENSE                 # MIT License
├── README.md               # This file
└── setup.sh                # Setup script
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- MySQL 5.7 or higher
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/textile-admin.git
cd textile-admin
./setup.sh
```

### Setup Database

1. Create a database in MySQL:

```sql
CREATE DATABASE textile_admin;
```

2. Run the schema script to create the tables:

```bash
mysql -u username -p textile_admin < scripts/schema.sql
```

### Configuration

The application can be configured using environment variables:

- `SERVER_ADDRESS`: Server address and port (default: ":8080")
- `UPLOAD_DIR`: Directory for storing uploaded files (default: "uploads")
- `DB_HOST`: Database host (default: "localhost")
- `DB_PORT`: Database port (default: 3306)
- `DB_USER`: Database username (default: "root")
- `DB_PASSWORD`: Database password (default: "password")
- `DB_NAME`: Database name (default: "textile_admin")
- `FILE_URL_PREFIX`: URL prefix for file downloads (default: "http://localhost:8080/files")

### Running the Application

```bash
go run cmd/api/main.go
```

## API Endpoints

### Upload File

```
POST /api/reading/upload
Content-Type: multipart/form-data

Parameters:
- file: The file to upload
- user_id: The user ID
```

### Get Task by ID

```
GET /api/reading/task/:task_id
```

### Get Tasks by User ID

```
GET /api/reading/tasks/user/:user_id
```

### Update Task Status

```
PUT /api/reading/task/:task_id/status
Content-Type: application/json

Body:
{
  "status": "pending" | "processing" | "completed" | "failed"
}
```

### Download File

```
GET /files/:file_name
```

## Technical Implementation

- The application uses GORM as an Object-Relational Mapper for database operations
- Database migrations are handled automatically using GORM AutoMigrate
- All initialization steps are grouped into a single function for better organization

## License

[MIT](LICENSE)

## 配置系统

本项目使用基于 YAML 的配置系统，支持不同的运行环境（开发环境和生产环境）。

### 配置文件

配置文件位于 `configs/` 目录下：

- `config.dev.yaml` - 开发环境配置
- `config.prod.yaml` - 生产环境配置

### 配置结构

配置文件包含以下主要部分：

```yaml
server:
  address: ":8080"              # 服务器监听地址
  upload_dir: "uploads"         # 文件上传目录
  file_url_prefix: "..."        # 文件URL前缀

database:
  host: "localhost"             # 数据库主机
  port: 3306                    # 数据库端口
  user: "root"                  # 数据库用户名
  password: ""                  # 数据库密码
  dbname: "textile_admin"       # 数据库名称

log:
  level: "debug"                # 日志级别 (debug, info, warn, error, fatal)
  format: "text"                # 日志格式 (text, json)
```

### 环境变量

系统支持通过环境变量覆盖配置文件中的设置：

- `APP_ENV` - 运行环境，可选值：`dev`（默认）或 `prod`
- `SERVER_ADDRESS` - 服务器监听地址
- `UPLOAD_DIR` - 文件上传目录
- `FILE_URL_PREFIX` - 文件URL前缀
- `DB_HOST` - 数据库主机
- `DB_PORT` - 数据库端口
- `DB_USER` - 数据库用户名
- `DB_PASSWORD` - 数据库密码
- `DB_NAME` - 数据库名称
- `LOG_LEVEL` - 日志级别
- `LOG_FORMAT` - 日志格式

此外，配置文件中可以使用 `${ENV_VAR}` 语法引用环境变量，例如：

```yaml
database:
  password: "${DB_PASSWORD}"
```

### 运行应用

使用提供的脚本在特定环境中运行应用：

```bash
# 开发环境
./scripts/run.sh dev

# 生产环境
DB_PASSWORD=your_secure_password ./scripts/run.sh prod
```

## 开发

### 依赖

- Go 1.16+
- MySQL 5.7+

### 安装依赖

```bash
go mod download
```

### 构建

```bash
go build -o bin/textile-admin cmd/api/main.go
```

### 运行测试

```bash
go test ./...
``` 