# Go Chat Application

Author: Thiago Menezes - thg.mnzs@gmail.com

This is a simple chat application built using Go, Echo framework, and WebSockets. Users can log in, join chatrooms, send messages, and fetch stock quotes using a chatbot command. The application uses RabbitMQ as a message broker for communication between the chat server and the stock quote bot, and PostgreSQL for user and chat data storage.

## Features

- User registration and authentication
- Fully featured RESTful endpoints for authentication, changing password, and CRUD operations on the user entity
- JWT authentication and session
- Application configuration via environment variables
- RBAC (role-based access control)
- Structured logging
- Great performance
- Request marshaling and data validation
- API Docs using SwaggerUI
- Mocking using stdlib
- Complete test coverage
- Containerized database query tests
- Real-time messaging using WebSockets
- Stock quote fetching using Stooq API
- Message ordering by timestamp and limiting to the last 50 messages
- (Optional) Multiple chatrooms
- (Optional) Error handling for the bot

## Libraries

1. Echo - HTTP 'framework'.
2. Go-Pg - PostgreSQL ORM
3. JWT-Go - JWT Authentication
4. Zerolog - Structured logging
5. Bcrypt - Password hashing
6. Yaml - Unmarshalling YAML config file
7. Validator - Request validation.
8. lib/pq - PostgreSQL driver
9. zxcvbn-go - Password strength checker
10. DockerTest - Testing database queries
11. Testify/Assert - Asserting test results

## Prerequisites

- [Go](https://golang.org/doc/install) (1.20 or later)
- [RabbitMQ](https://www.rabbitmq.com/download.html)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

1. Clone the repository:

```sh
git clone https://github.com/yourusername/go-chat-application.git
cd go-chat-application
```

2. Install dependencies:

```sh
go mod download
```

3. Configure the application:

Set the following environment variables in your system:

```
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
STOOQ_API_URL=https://stooq.com/q/l/
DATABASE_URL=postgres://postgres:postgres@localhost:5432/chat?sslmode=disable
JWT_SECRET=8d1b46e1f34a2b2c599ae78d7f53c0b98d1b46e1f34a2b2c599ae78d7f53c0b9
```

Replace the placeholder values with your RabbitMQ connection URL, Stooq API URL, PostgreSQL connection URL, and JWT secret if needed.

4. Run the application:

```sh
go run cmd/api/main.go
```

The chat server will start, and you can access the chatroom at `http://localhost:8080`.

## Usage

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `POST /login`: accepts username/passwords and returns jwt token and refresh token
* `GET /refresh/:token`: refreshes sessions and returns jwt token
* `GET /me`: returns info about currently logged in user
* `GET /swaggerui/` (with trailing slash): launches swaggerui in browser
* `GET /v1/users`: returns list of users
* `GET /v1/users/:id`: returns single user
* `POST /v1/users`: creates a new user
* `PATCH /v1/password/:id`: changes password for a user
* `DELETE /v1/users/:id`: deletes a user

To use the chat application:

1. Register a new user or log in with an existing user.
2. Join a chatroom and send messages.
3. To fetch a stock quote, use the following command format:

```
/stock=stock_code
```

For example:

```
/stock=aapl.us
```

The stock quote bot will fetch the stock quote and display it in the chatroom.

## Project Structure

1. Root directory contains things not related to code directly, e.g. docker-compose, CI/CD, readme, bash scripts etc. It should also contain vendor folder, Gopkg.toml and Gopkg.lock if dep is being used.

2. Cmd package contains code for starting applications (main packages). The directory name for each application should match the name of the executable you want to have. Jobsity is structured as a monolith application but can be easily restructured to contain multiple microservices. An application may produce multiple binaries, therefore Jobsity uses the Go convention of placing main package as a subdirectory of the cmd package. As an example, scheduler application's binary would be located under cmd/cron. It also loads the necessery configuration and passes it to the service initializers.

3. Rest of the code is located under /pkg. The pkg directory contains `utl` and 'microservice' directories.

4. Microservice directories, like api (naming corresponds to `cmd/` folder naming) contains multiple folders for each domain it interacts with, for example: user, car, appointment etc.

5. Domain directories, like user, contain all application/business logic and two additional directories: platform and transport.

6. Platform folder contains various packages that provide support for things like databases, authentication or even marshaling. Most of the packages located under platform are decoupled by using interfaces. Every platform has its own package, for example, postgres, elastic, redis, memcache etc.

7. Transport package contains HTTP handlers. The package receives the requests, marshals, validates then passes it to the corresponding service.

8. Utl directory contains helper packages and models. Packages such as mock, middleware, configuration, server are located here.

