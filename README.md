# Gouril

**A lightweight URL shortener built with Go and SQLite.**

Gouril is a backend project designed to explore Go's standard HTTP library, REST API fundamentals, database integration, input validation, and graceful server shutdown.

## Features

* **URL shortening:** Generate a unique, six-character short code for a URL.
* **URL redirection:** Redirect short URLs to their original destinations using HTTP `302 Found`.
* **Persistent storage:** Store URL mappings in SQLite so they survive server restarts.
* **Input validation:** Validate submitted URLs and handle malformed JSON and missing fields.
* **HTTP status codes:** Return appropriate responses for successful requests, invalid input, unsupported methods, and missing URLs.
* **Graceful shutdown:** Shut down the HTTP server cleanly when interrupted.

## Tech Stack

* **Language:** Go
* **HTTP server:** Go standard library (`net/http`)
* **Database:** SQLite
* **SQLite driver:** `modernc.org/sqlite`

## Project Structure

```text
Gouril/
├── database/
│   └── database.go
├── handler/
│   └── handler.go
├── Gouril.go
├── go.mod
├── go.sum
└── README.md
```

* `Gouril.go` — Initializes the database, configures HTTP routes, and runs the server.
* `database/database.go` — Handles database initialization, URL insertion, short-code lookups, and existence checks.
* `handler/handler.go` — Implements URL shortening, request validation, and redirection.

## Getting Started

### Prerequisites

* [Go](https://go.dev/dl/) compatible with the version specified in `go.mod`.
* Git.

### 1. Clone the repository

```bash
git clone https://github.com/MobinEbrahimkhani/Gouril.git
cd Gouril
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Run the server

```bash
go run .
```

Gouril starts at:

```text
http://localhost:8080
```

On startup, Gouril initializes the SQLite database and creates the `urls` table if it doesn't already exist. The database file, `gouril.db`, is created in the application's working directory when needed.

## API Documentation

### 1. Shorten a URL

`POST /shorten`

**Request body**

```json
{
  "url": "https://go.dev"
}
```

**Example request**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://go.dev"}'
```

**Example response — `201 Created`**

```json
{
  "short_code": "aB12xY",
  "short_url": "http://localhost:8080/aB12xY"
}
```

The short code in this example is illustrative; actual codes are generated dynamically.

### 2. Redirect to the original URL

`GET /{shortCode}`

Open the generated short URL in a browser or use curl:

```bash
curl -i http://localhost:8080/aB12xY
```

If the code exists, Gouril responds with `302 Found` and a `Location` header containing the original URL.

### HTTP Responses

| Status                      | Meaning                                                       |
| --------------------------- | ------------------------------------------------------------- |
| `201 Created`               | URL shortened successfully                                    |
| `302 Found`                 | Redirect to the original URL                                  |
| `400 Bad Request`           | Invalid JSON, missing URL, invalid URL, or missing short code |
| `404 Not Found`             | Short code does not exist                                     |
| `405 Method Not Allowed`    | HTTP method is unsupported                                    |
| `500 Internal Server Error` | Database operation failed                                     |

## Database Schema

Gouril uses a single SQLite table named `urls`.

| Column         | Description                   |
| -------------- | ----------------------------- |
| `id`           | Auto-incrementing primary key |
| `short_code`   | Unique short identifier       |
| `original_url` | Original destination URL      |

The `short_code` column has a uniqueness constraint to prevent duplicate mappings.

## Graceful Shutdown

Press `Ctrl+C` in the terminal to request a graceful shutdown. The server allows up to five seconds for ongoing requests to finish.

## Current Limitations

* The public base URL and server port are currently hardcoded for local development.
* There is no authentication, user management, or rate limiting.
* Automated tests are not yet included.
* Short-code generation and collision handling can be strengthened further.
* No web interface is currently provided; the API can be accessed with tools such as curl.

## Future Improvements

* Add automated unit and integration tests.
* Improve short-code generation and collision handling.
* Add a health-check endpoint.
* Configure the server using environment variables.
* Add structured logging and more robust error handling.
* Add deployment instructions and containerization.

## Learning Goals

Gouril is a hands-on project for developing practical Go backend skills, including HTTP handlers, JSON processing, SQL queries, modular code organization, and server lifecycle management.

---

**Built with Go.**
