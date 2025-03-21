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
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- MySQL 5.7 or higher

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