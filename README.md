# ABG Backend

A Go-based backend service with authentication and user management.

## Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Git

## Setup

1. Clone the repository:
```bash
git clone https://github.com/How-to-get-ABG/backend.git
cd backend
```

2. Create a `.env` file in the root directory with the following variables:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_db

# JWT Configuration
JWT_SECRET=bfa3e0190f2c8e7d98e9a3ecb0c3e2b59f4f25d6c6a0be8eb25a0bfa3e9eae13

# Server Configuration (optional, defaults to 8080)
PORT=8080
```

3. Install dependencies:
```bash
go mod tidy
```

## Building

To build the project, run:
```bash
go build -o main.exe
```

This will create an executable file named `main.exe` in the backend directory.

## Running

To run the server:
```bash
./main.exe
```

The server will start on port 8080 by default (or the port specified in your .env file).

## API Endpoints

### Public Routes
- `POST /api/register` - Register a new user
- `POST /api/login` - Login and get JWT token

### Protected Routes
- `GET /api/me` - Get user profile (requires JWT token)

## Development

The project structure is organized as follows:
```
backend/
├── internal/
│   ├── config/     # Configuration management
│   ├── database/   # Database connection and queries
│   ├── handlers/   # HTTP request handlers
│   └── middleware/ # HTTP middleware
├── main.go         # Application entry point
└── .env           # Environment variables
```

## Testing

To run tests:
```bash
go test ./...
```

## Contributing

1. Pull the repository
2. Create your feature branch (`git checkout -b username/feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin username/feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 