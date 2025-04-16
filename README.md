# Go Expense Tracker

A simple expense tracking application written in Go.

## Configuration

### Database Configuration

The application can use either an in-memory storage or a PostgreSQL database for storing expenses.

To use the PostgreSQL database, you need to:

1. Create a `.env` file in the root directory of the project with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=expense_tracker
```

2. Run the application with the `-db` flag:

```
go run main.go -db
```

The `-db` flag enables the use of PostgreSQL database instead of in-memory storage.

### Environment Variables

The following environment variables can be set in the `.env` file:

- `DB_HOST`: PostgreSQL host (default: "localhost")
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: PostgreSQL user (default: "postgres")
- `DB_PASSWORD`: PostgreSQL password (default: "postgres")
- `DB_NAME`: PostgreSQL database name (default: "expense_tracker")

## Running with Docker

You can run the PostgreSQL database using Docker Compose:

```
docker-compose up -d
```

This will start a PostgreSQL container with the configuration specified in the `docker-compose.yml` file.