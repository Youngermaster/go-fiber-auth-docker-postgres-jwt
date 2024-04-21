# Go Fiber + JWT Auth + Docker + PostgreSQL + PgAdmin Boilerplate

This boilerplate provides a starting point for Go Fiber that utilizes Docker, PostgreSQL, JWT for authentication, and
PgAdmin for database management.

## Development Setup

### Prerequisites

- Docker must be installed on your system for an optimal development experience.
- Clone the repository and navigate to the project directory.

### Environment Configuration

Copy the `.env.example` file to a new file named `.env` and adjust the environment variables:

```sh
DB_PORT=5432
DB_USER=example_user
DB_PASSWORD=example_password
DB_NAME=example_db
SECRET=example_secret
PGADMIN_DEFAULT_EMAIL=user@domain.com
PGADMIN_DEFAULT_PASSWORD=SecurePassword
```

Ensure there are no port conflicts or conflicting Docker containers running. If necessary, adjust the ports in the
`.env` file and `docker-compose.yml`.

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

    ### Creating a User

    To create a user via the API, send a POST request to `http://localhost:3000/api/user/` with the following JSON
    payload:

    ```json
    {
    "username": "johndoe",
    "email": "johndoe@test.com",
    "password": "1234567890"
    }
    ```

    You can use tools like `curl`, Postman, or any HTTP client in your programming language of choice.

    ## Troubleshooting

    Ensure all environment variables are set correctly in your `.env` file, as incorrect settings may prevent the
    services from starting properly.

    Check the Docker logs if any service fails to start:

    ```sh
    docker-compose logs <service-name>
      ```

      Replace `<service-name>` with `web`, `db`, or `pgadmin` to view logs for a specific service.

        ---

        For further details, refer to the Go Fiber, Docker, and PostgreSQL documentation. This setup is ideal for
        development environments and should be adapted for production use with security best practices.
        ```
