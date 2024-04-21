# Go Fiber + JWT Auth + Docker + PostgreSQL Boilerplate

Boilerplate for Go Fiber that makes use of Docker, PostgreSQL and JWT.

## Development

### In order to have an optimal development experience you need to have Docker installed

Set the environment variables in a `.env` file:

- DB_PORT=5432
- DB_USER=example_user
- DB_PASSWORD=example_password
- DB_NAME=example_db
- SECRET=example_secret

Be sure you don't have any conflicting containers.
Then run the commands:

```sh
docker-compose up -d
```

This should start the API and the database.

## Database Management

You can manage the database via `psql` with the command:

```sh
docker-compose exec db psql -U <DB_USER (from the environment variables)>
```
