#!/bin/bash

if [ "$#" -lt 2 ]; then
  echo "Usage: $0 <connection-string> <command>"
  echo "Example: $0 \"user=postgres password=Thang@240803 dbname=go-ecommerce host=localhost port=5432 sslmode=disable\" up"
  exit 1
fi

CONNECTION_STRING=$1
COMMAND=$2

MIGRATIONS_DIR="./migrations"

goose -dir "$MIGRATIONS_DIR" postgres "$CONNECTION_STRING" "$COMMAND"
