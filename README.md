# Retot API

Houses the backend for the Retot application.

## Requirements

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/dl/) 1.18 or higher

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/gofiber/recipes.git
   cd recipes/auth-docker-postgres-jwt
   ```

2. Set the environment variables in a `.env` file:

   ```env
   DB_PORT=5432
   DB_USER=example_user
   DB_PASSWORD=example_password
   DB_NAME=example_db
   SECRET=example_secret
   ```

3. Build and start the Docker containers:
   ```bash
   docker-compose build
   docker-compose up
   ```

The API and the database should now be running.

## Database Management

You can manage the database via `psql` with the following command:

```bash
docker-compose exec db psql -U <DB_USER>
```

Replace `<DB_USER>` with the value from your `.env` file.
