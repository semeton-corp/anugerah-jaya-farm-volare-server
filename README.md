# Anugerah Jaya Farm Volare Server

A backend server for managing poultry farm operations, built with Go and Fiber framework.

## Features

- **Chicken Management**: Track chicken cages, placements, and inventory
- **Feed Management**: Monitor feed stock and distribution per cage
- **Health Monitoring**: Track chicken health metrics and performance
- **Sales Management**: Handle afkir chicken sales and customer transactions
- **Work Tracking**: Log additional work and user activities
- **Cashflow Management**: Track financial transactions and payment history

## Tech Stack

- **Framework**: [Fiber v2](https://github.com/gofiber/fiber) - Fast HTTP web framework
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Cache**: Redis
- **Authentication**: JWT (JSON Web Tokens)
- **Storage**: S3 Biznet Gio
- **Logging**: Uber Zap with Lumberjack for log rotation
- **Configuration**: Viper
- **Excel Export**: Excelize
- **Email**: GoMail

## Prerequisites

- Go 1.24.0 or higher
- Docker and Docker Compose
- PostgreSQL
- Redis
- S3

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/semeton-corp/anugerah-jaya-farm-volare.git
cd anugerah-jaya-farm-volare-server
```

### 2. Configure environment

Copy the example environment files and configure them:

```bash
cp env.example.yaml env.yaml
cp compose.env.example compose.env
```

Edit `env.yaml` with your configuration:

- Database connection details
- JWT secret and settings
- AWS S3 credentials
- Server settings
- CORS configuration

### 3. Run with Docker

Start all services (PostgreSQL, Redis, and the application):

```bash
make compose-up
```

Or for development mode with dependencies only:

```bash
make compose-dev
```

### 4. Run locally (development)

Install dependencies:

```bash
go mod download
```

Run with hot reload using Air:

```bash
make air
```

Or run directly:

```bash
make run
```

## Available Make Commands

- `make run` - Run the application directly
- `make build` - Build the application binary
- `make air` - Run with Air hot reload
- `make compose-up` - Start all Docker containers (production)
- `make compose-down` - Stop all Docker containers
- `make compose-dev` - Start development containers (DB & Redis only)
- `make compose-dev-down` - Stop development containers
- `make clean-go-cache` - Clean Go build cache
- `make clean-docker-cache` - Clean Docker system cache

## Project Structure

```
.
├── cmd/app/            # Application entry point
├── internal/           # Internal application code
│   ├── dto/           # Data Transfer Objects
│   ├── entity/        # Database entities/models
│   ├── handler/       # HTTP request handlers
│   ├── listener/      # Event listeners
│   ├── mapper/        # Data mapping utilities
│   ├── middleware/    # HTTP middlewares
│   ├── repository/    # Database repositories
│   └── service/       # Business logic layer
├── pkg/               # Public packages
├── bootstrap/         # Application bootstrap
├── infra/            # Infrastructure setup
├── seed/             # Database seeders
├── templates/        # Email/document templates
└── log/              # Application logs
```

## Documentation

- [Postman](https://documenter.getpostman.com/view/25537573/2sB2cYdgJk)

## Configuration

The application uses a YAML configuration file (`env.yaml`) with the following sections:

- **log**: Logging level and environment
- **database**: PostgreSQL connection settings
- **jwt**: Authentication configuration
- **server**: HTTP server settings and CORS
- **app**: Application metadata

## Development

The project uses Air for hot reloading during development. Configuration is in `.air.toml`.

## License

Copyright © 2025 Semeton Corp

## Authors

- **Indra** - [@indrabrata](https://github.com/indrabrata)
