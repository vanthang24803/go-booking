# Booking System

This is a booking system API written in Go and Fiber. It uses Postgres as the database and Redis as the cache layer.

## Features

- User authentication using JWT
- User profile management
- Booking management
- Payment management
- Cache layer using Redis

## Installation

- Clone the repository
- Run `go build` to build the binary
- Run `go run ./cmd/main.go` to start the server or start with air `air`

## API Endpoints

- `POST /login`: login and get JWT token
- `POST /register`: register a new user
- `GET /me`: get the current user's profile
- `PUT /me`: update the current user's profile
- `GET /me/avatar`: upload a new avatar

## Environment Variables

- `PORT`: the port to listen on (default is 8080)
- `DB_HOST`: the host of the Postgres database
- `DB_USER`: the username to use when connecting to the Postgres database
- `DB_PASSWORD`: the password to use when connecting to the Postgres database
- `DB_NAME`: the name of the Postgres database
- `REDIS_ADDR`: the address of the Redis server
- `REDIS_PASSWORD`: the password to use when connecting to the Redis server
- `REDIS_DB`: the database number to use when connecting to the Redis server
- `JWT_SECRET`: the secret to use when generating JWT tokens
