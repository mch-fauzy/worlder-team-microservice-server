# Worlder Team Microservice Server

A comprehensive microservice backend system built with Go, Echo, MySQL, GORM, and gRPC following clean architecture principles.

## Architecture Overview

This system consists of two main microservices:

### Microservice A (Data Generator Service)
- Generates sensor data streams with configurable frequency
- Multiple instances can run with different sensor types
- REST API for frequency control
- gRPC client to send data to Microservice B

### Microservice B (Data Storage Service)
- Receives sensor data via gRPC
- Stores data in MySQL database using GORM
- Comprehensive REST API for data management
- Authentication & authorization
- Pagination support
- Scalable to handle multiple data sources

## System Architecture Diagram

```mermaid
graph TB
    %% External clients
    Client[🌐 HTTP Clients<br/>Postman/Browser/API]
    
    %% Load balancer
    LB[🔄 Nginx Load Balancer<br/>Port: 80]
    
    %% Microservice A instances
    subgraph "Microservice A (Data Generators)"
        A1[🌡️ Temperature Generator<br/>Port: 8080<br/>REST API + gRPC Client]
        A2[💧 Humidity Generator<br/>Port: 8082<br/>REST API + gRPC Client]  
        A3[📊 Pressure Generator<br/>Port: 8083<br/>REST API + gRPC Client]
        A4[💡 Light Generator<br/>Port: 8084<br/>REST API + gRPC Client]
        A5[🚶 Motion Generator<br/>Port: 8085<br/>REST API + gRPC Client]
    end
    
    %% Microservice B
    subgraph "Microservice B (Storage Service)"
        B[🗄️ Storage Service<br/>Port: 8081<br/>REST API + gRPC Server]
        
        subgraph "Clean Architecture Layers"
            BH[📝 Handlers Layer<br/>HTTP/gRPC Controllers]
            BS[⚙️ Services Layer<br/>Business Logic]
            BR[📦 Repository Layer<br/>Data Access]
        end
        
        subgraph "Modules"
            AUTH[🔐 Auth Module<br/>JWT Authentication]
            SENSOR[📊 Sensor Data Module<br/>CRUD Operations]
            HEALTH[❤️ Health Module<br/>Status Monitoring]
        end
    end
    
    %% Database
    DB[(🗃️ MySQL Database<br/>Port: 3306<br/>GORM ORM)]
    
    %% Redis (optional caching)
    REDIS[(🔴 Redis Cache<br/>Port: 6379<br/>Session Storage)]
    
    %% External monitoring
    SWAGGER[📚 Swagger UI<br/>API Documentation<br/>localhost:8081/swagger]
    
    %% Connections
    Client -.-> LB
    LB --> B
    Client --> A1
    Client --> A2  
    Client --> A3
    Client --> A4
    Client --> A5
    
    %% gRPC connections from generators to storage
    A1 -.->|gRPC<br/>SendSensorData| B
    A2 -.->|gRPC<br/>SendSensorData| B
    A3 -.->|gRPC<br/>SendSensorData| B
    A4 -.->|gRPC<br/>SendSensorData| B
    A5 -.->|gRPC<br/>SendSensorData| B
    
    %% Internal architecture
    B --> BH
    BH --> BS
    BS --> BR
    BR --> DB
    
    B --> AUTH
    B --> SENSOR
    B --> HEALTH
    
    %% Database connections
    AUTH -.-> DB
    SENSOR -.-> DB
    B -.-> REDIS
    
    %% Documentation
    B -.-> SWAGGER
    
    %% Styling
    classDef microserviceA fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef microserviceB fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef database fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px
    classDef client fill:#fff3e0,stroke:#e65100,stroke-width:2px
    classDef infrastructure fill:#fce4ec,stroke:#880e4f,stroke-width:2px
    
    class A1,A2,A3,A4,A5 microserviceA
    class B,BH,BS,BR,AUTH,SENSOR,HEALTH microserviceB
    class DB,REDIS database
    class Client,SWAGGER client
    class LB infrastructure
```

### Architecture Principles

- **Clean Architecture**: Each microservice follows clean architecture with clear separation of concerns
- **Modularity**: Feature-based module organization for maintainability
- **Scalability**: Horizontal scaling through multiple generator instances
- **Communication**: gRPC for inter-service communication, REST for client APIs
- **Data Flow**: Unidirectional data flow from generators to storage
- **Authentication**: JWT-based authentication with role-based access control
- **Monitoring**: Health checks and comprehensive logging

### Data Flow

1. **Data Generation**: Microservice A instances generate sensor data at configurable frequencies
2. **gRPC Transmission**: Generated data is sent to Microservice B via gRPC for efficiency
3. **Data Storage**: Microservice B stores received data in MySQL using GORM
4. **API Access**: Clients can query stored data through REST APIs with pagination and filtering
5. **Authentication**: All API requests are authenticated using JWT tokens
6. **Load Balancing**: Nginx distributes HTTP traffic across service instances

## Project Structure

```
worlder-team-microservice-server/
├── microservice-a/              # Data Generator Service ✅ MIGRATED
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Application entry point
│   ├── modules/               # Feature-based organization (MIGRATED!)
│   │   ├── generator/         # Data generation & sensor management
│   │   │   ├── entities/      # SensorData & GeneratorStatus entities
│   │   │   ├── handlers/      # Generator HTTP handlers
│   │   │   ├── services/      # Data generation business logic
│   │   │   ├── interfaces/    # Generator & client interfaces
│   │   │   ├── dtos/          # Frequency request/response DTOs
│   │   │   └── grpc/          # gRPC client implementation
│   │   └── health/            # Health check endpoints
│   │       └── handlers/      # Health handlers
│   ├── configs/               # Configuration management
│   │   └── config.go         # Config struct and loading
│   ├── shared/               # Local shared types (APIResponse)
│   └── Dockerfile            # Container build configuration
├── microservice-b/              # Data Storage Service ✅ MIGRATED
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Application entry point
│   ├── modules/               # Feature-based organization (MIGRATED!)
│   │   ├── auth/              # Authentication & user management
│   │   │   ├── entities/      # User entity with GORM tags
│   │   │   ├── handlers/      # Auth HTTP handlers (login, etc.)
│   │   │   ├── services/      # Auth & JWT token services
│   │   │   ├── interfaces/    # Auth service interfaces
│   │   │   └── dtos/          # Auth request/response DTOs
│   │   ├── sensor-data/       # Sensor data management
│   │   │   ├── entities/      # SensorData entity with GORM tags
│   │   │   ├── handlers/      # Sensor CRUD HTTP handlers
│   │   │   ├── services/      # Sensor business logic
│   │   │   ├── repositories/  # Data access layer (GORM)
│   │   │   ├── interfaces/    # Service & repository interfaces
│   │   │   ├── dtos/          # DTOs, filters & pagination
│   │   │   └── grpc/          # gRPC server implementation
│   │   └── health/            # Health check endpoints
│   │       └── handlers/      # Health check handlers
│   ├── configs/               # Configuration management
│   │   ├── config.go         # Config struct and loading
│   │   └── logger.go         # Logger configuration
│   ├── shared/               # Local shared types (APIResponse)
│   ├── migrations/            # Database schema migrations
│   ├── docs/                 # Auto-generated Swagger docs
│   └── Dockerfile            # Container build configuration
├── shared/                      # Shared components across services
│   ├── proto/                 # Protocol buffer definitions
│   │   └── sensor/           # Sensor service protobuf files
│   ├── constants/             # Shared constants (status codes, etc.)
│   ├── utils/                 # Utility functions (logging, parsing)
│   ├── middleware/            # Shared HTTP middleware
│   └── response.go           # Standard API response structure
├── infrastructures/            # Infrastructure configurations
│   └── nginx/                # Load balancer configuration
│       └── nginx.conf        # Nginx proxy settings
├── docs/                       # Project documentation
│   ├── architecture/          # Architecture diagrams
│   ├── api/                   # API documentation & Postman collections
│   └── erd/                   # Database entity relationship diagrams
├── .env                       # Environment variables (copy from .env.example)
├── .env.example              # Environment variables template
├── docker-compose.yml        # Multi-service orchestration
├── Makefile                  # Build and deployment automation
└── README.md                 # This file
```

## Features

### Microservice A Features
- **Data Generation**: Continuous sensor data generation with timestamps
- **Configurable Frequency**: REST API to control data generation frequency
- **Multiple Sensor Types**: Support for different sensor types per instance
- **gRPC Communication**: Efficient data transmission to Microservice B
- **Health Monitoring**: Health check endpoints

### Microservice B Features
- **Data Storage**: Persistent storage in MySQL with optimized queries
- **REST API Endpoints**:
  - Retrieve by ID combination
  - Retrieve by duration/time range
  - Retrieve by IDs + timestamps
  - Delete sensor values with filters
  - Edit/Update sensor values with filters
  - Pagination support
- **Authentication & Authorization**: JWT-based auth with role-based access
- **Data Validation**: Comprehensive input validation
- **Rate Limiting**: API rate limiting to prevent abuse
- **Monitoring**: Metrics and health checks

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+

### Using Makefile (Recommended)

1. **Clone and navigate to project**:
```bash
git clone https://github.com/mch-fauzy/worlder-team-microservice-server.git
cd worlder-team-microservice-server
```

2. **Set up environment variables**:
```bash
# Copy the example environment file
cp .env.example .env

# Edit the .env file with your preferred settings
nano .env
# or
vim .env
```

> **Important**: The `.env` file contains all configuration including database credentials, JWT secrets, and service ports. Review and modify the values as needed for your environment.

3. **Start all services with one command**:
```bash
make quick-start
```

### After Starting Services

Once services are running, you can access:

- **Swagger API Documentation**: http://localhost:8081/swagger/index.html
- **Temperature Sensor Service**: http://localhost:8080
- **Humidity Sensor Service**: http://localhost:8082  
- **Pressure Sensor Service**: http://localhost:8083
- **Storage Service (Microservice B)**: http://localhost:8081
- **Load Balancer**: http://localhost:80

### Available Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make quick-start` | Start all services quickly |
| `make run` | Start all services with Docker Compose |
| `make stop` | Stop all services |
| `make restart` | Restart all services |
| `make rebuild` | Rebuild and restart all services |
| `make rebuild-generators` | Rebuild only microservice-a generators |
| `make rebuild-storage` | Rebuild only microservice-b storage service |
| `make logs` | View logs from all services |
| `make clean` | Clean up containers and volumes |
| `make proto` | Generate protobuf files |
| `make swagger` | Generate Swagger documentation |

### Troubleshooting

**Services not starting?**
```bash
# 1. Check if .env file exists and is configured
ls -la .env
cat .env

# 2. Verify environment variables are loaded
docker-compose config

# 3. Clean and restart services
make clean
make quick-start
```

**Need to rebuild after changes?**
```bash
make rebuild
```

**Database connection issues?**
```bash
# Wait for MySQL to be ready, then restart
make logs
# Wait until you see "MySQL ready for connections"
make restart
```

## API Documentation

### Microservice A Endpoints
- `GET /health` - Health check
- `POST /frequency` - Set data generation frequency
- `GET /frequency` - Get current frequency
- `GET /status` - Get service status

### Microservice B Endpoints
- `POST /auth/login` - Authentication
- `GET /sensors` - List sensor data (with pagination)
- `GET /sensors/{id1}/{id2}` - Get by ID combination
- `GET /sensors/duration` - Get by time range
- `PUT /sensors/{id}` - Update sensor data
- `DELETE /sensors` - Delete with filters

Full API documentation is available at `/swagger` when running the services.
