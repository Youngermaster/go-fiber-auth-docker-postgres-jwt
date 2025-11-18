# Go Fiber + JWT Auth + Docker + PostgreSQL + PgAdmin Boilerplate

A production-ready Go Fiber boilerplate with secure JWT authentication, refresh token system, session management, and comprehensive security features.

## Features

- **Dual-Token Authentication**: Short-lived access tokens (15 min) and long-lived refresh tokens (7 days)
- **Automatic Token Rotation**: Refresh tokens are automatically rotated on each use to prevent replay attacks
- **Session Management**: Track and manage user sessions across multiple devices with logout and logout-all capabilities
- **Configuration Validation**: Application validates all environment variables and JWT secrets at startup
- **Rate Limiting**: 5 requests per minute on authentication endpoints to prevent brute force attacks
- **Password Security**: Bcrypt hashing with cost factor 14
- **API Versioning**: All routes are versioned under `/api/v1` for future compatibility
- **Health Checks**: Readiness and liveness probes for orchestration platforms
- **Docker Support**: Fully containerized development environment with hot-reloading

## Development Setup

### Prerequisites

- Docker must be installed on your system for an optimal development experience.
- Clone the repository and navigate to the project directory.

### Environment Configuration

Copy the `.env.example` file to a new file named `.env` and configure the environment variables.

**CRITICAL**: Generate strong, unique secrets for production. All three JWT secrets must be different and at least 32 characters long:

```bash
# Generate three different secrets (run this command 3 times)
openssl rand -base64 32
```

Configure your `.env` file with the generated secrets:

```sh
# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_USER=example_user
DB_PASSWORD=example_password
DB_NAME=example_db

# JWT Secrets - MUST be different and 32+ characters
SECRET=<paste-first-generated-secret>
ACCESS_TOKEN_SECRET=<paste-second-generated-secret>
REFRESH_TOKEN_SECRET=<paste-third-generated-secret>

# PgAdmin Configuration
PGADMIN_DEFAULT_EMAIL=user@domain.com
PGADMIN_DEFAULT_PASSWORD=SecurePassword
```

The application will validate your configuration at startup and refuse to start if:

- Any required environment variable is missing
- JWT secrets are less than 32 characters
- JWT secrets contain weak/default values (e.g., "password", "secret", "test")
- All three JWT secrets are not different from each other

Ensure there are no port conflicts or conflicting Docker containers running. If necessary, adjust the ports in the `.env` file and `docker-compose.yml`.

### Starting the Services

Run the following command to start all services defined in the `docker-compose.yml`:

```sh
docker-compose up -d
```

This command will start the API, PostgreSQL database, and PgAdmin.

## Database Management

### Using PgAdmin

PgAdmin is configured to run on port 5050. Access it by navigating to `http://localhost:5050` in your web browser. Login
with the PGADMIN_DEFAULT_EMAIL and PGADMIN_DEFAULT_PASSWORD specified in your `.env` file.

#### Connecting to PostgreSQL through PgAdmin

1. Open PgAdmin and login.
2. Right-click on 'Servers' in the left sidebar and select 'Create' -> 'Server'.
3. Enter a name for the connection in the 'General' tab.
4. Switch to the 'Connection' tab:

- Hostname/address: `db`
- Port: `5432` (or your custom DB_PORT)
- Username: as per `DB_USER`
- Password: as per `DB_PASSWORD`
- Save the password for ease of use.

### Using psql

To connect directly to the database via `psql`, use the script provided:

```sh
./manually_connect_to_db.sh
```

Or use Docker Compose:

```sh
docker-compose exec db psql -U <DB_USER>
```

Replace `<DB_USER>` with the actual database user name from your `.env` file.

## API Usage

The API is versioned and accessible at `http://localhost:3000/api/v1`. All authentication endpoints are rate-limited to 5 requests per minute per IP address.

### Authentication Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/auth/login` | No | Login and receive access + refresh tokens |
| POST | `/auth/refresh` | No | Exchange refresh token for new token pair |
| POST | `/auth/logout` | Yes | Logout from current device |
| POST | `/auth/logout-all` | Yes | Logout from all devices |
| GET | `/auth/sessions` | Yes | Get all active sessions |

### User Management Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/users` | No | Create a new user (registration) |
| GET | `/users/:id` | Yes | Get user by ID |
| PATCH | `/users/:id` | Yes | Update user (ownership enforced) |
| DELETE | `/users/:id` | Yes | Delete user (ownership enforced) |

### Product Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/product` | Yes | Create a product |
| GET | `/product` | No | Get all products (paginated) |
| GET | `/product/:id` | No | Get product by ID |
| PATCH | `/product/:id` | Yes | Update product (ownership enforced) |
| DELETE | `/product/:id` | Yes | Delete product (ownership enforced) |

### Health Check Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Basic health check |
| GET | `/health/ready` | Readiness probe (checks database) |
| GET | `/health/live` | Liveness probe for Kubernetes |

### Example: Authentication Flow

**1. Login**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identity": "user@example.com",
    "password": "your_password"
  }'
```

Response includes `access_token` (15-minute lifespan) and `refresh_token` (7-day lifespan).

**2. Use Access Token**
```bash
curl -X GET http://localhost:3000/api/v1/auth/sessions \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**3. Refresh Tokens**
```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

This returns a new token pair and automatically revokes the old refresh token (token rotation).

For comprehensive API documentation and examples, see [AUTH_GUIDE.md](AUTH_GUIDE.md).

## Project Structure

```plaintext
.
├── cmd/
│   └── main.go                 # Application entry point with config validation
├── config/
│   └── validation.go           # Environment and JWT secret validation
├── database/
│   └── connect.go              # Database connection and migration
├── handler/
│   ├── auth.go                 # Authentication handlers (login, logout, refresh)
│   ├── token.go                # Token generation, validation, and rotation
│   ├── user.go                 # User management handlers
│   ├── product.go              # Product handlers
│   ├── password.go             # Password hashing and validation
│   ├── validation.go           # Input validation helpers
│   └── response.go             # Standardized response utilities
├── middleware/
│   └── auth.go                 # JWT authentication middleware
├── model/
│   ├── user.go                 # User model
│   ├── product.go              # Product model
│   └── session.go              # Session/refresh token model
├── router/
│   ├── router.go               # Main router setup
│   ├── auth.go                 # Authentication routes
│   ├── user.go                 # User routes
│   ├── product.go              # Product routes
│   └── health.go               # Health check routes
├── docker-compose.yml          # Docker services configuration
├── Dockerfile                  # Application container definition
├── .env.example                # Example environment configuration
└── go.mod                      # Go module dependencies
```

## Documentation

- **[AUTH_GUIDE.md](AUTH_GUIDE.md)** - Comprehensive authentication guide with API examples
- **[SECURITY.md](SECURITY.md)** - Security features and best practices
- **[SECURITY_AUDIT_AND_IMPROVEMENTS.md](SECURITY_AUDIT_AND_IMPROVEMENTS.md)** - Complete security audit and improvements log

## Troubleshooting

### Configuration Validation Failed

If the application refuses to start with configuration errors:

1. Generate proper secrets: `openssl rand -base64 32`
2. Ensure all three JWT secrets are different
3. Verify all required environment variables are set in `.env`
4. Check that secrets are at least 32 characters long

### Invalid or Expired Refresh Token

If refresh endpoint returns 401:

- Token may have expired (older than 7 days)
- Token was already used and rotated to a new one
- User logged out or logged out from all devices
- Solution: User must login again

### Rate Limit Exceeded

If receiving 429 status code:

- Wait 1 minute before retrying
- Implement exponential backoff in your client
- Authentication endpoints are limited to 5 requests per minute per IP

### Docker Issues

Check the Docker logs if any service fails to start:

```sh
docker-compose logs <service-name>
```

Replace `<service-name>` with `web`, `db`, or `pgadmin` to view logs for a specific service.

Ensure all environment variables are set correctly in your `.env` file, as incorrect settings may prevent the services from starting properly.

## Technology Stack

- **Go 1.24** - Programming language
- **Fiber v2.52.9** - Web framework
- **GORM v1.31.1** - ORM for database operations
- **PostgreSQL 18** - Primary database
- **Docker & Docker Compose** - Containerization
- **JWT** - Token-based authentication
- **Bcrypt** - Password hashing

## Security Features

- Dual-token authentication with automatic rotation
- Session tracking and management
- Rate limiting on authentication endpoints
- Configuration validation at startup
- Password hashing with bcrypt (cost 14)
- Input validation and sanitization
- Ownership enforcement on resource operations
- CORS with secure defaults

## Production Deployment

Before deploying to production:

1. **Use HTTPS**: Configure TLS/SSL via reverse proxy (Nginx, Caddy)
2. **Generate Strong Secrets**: Use `openssl rand -base64 32` for all JWT secrets
3. **Database SSL**: Enable SSL for PostgreSQL connections
4. **Environment Variables**: Use secret managers (AWS Secrets Manager, HashiCorp Vault)
5. **Monitor Sessions**: Implement session cleanup and suspicious activity alerts
6. **Update Dependencies**: Keep all dependencies up to date

For detailed production deployment guidance, see [SECURITY.md](SECURITY.md).

## License

MIT
