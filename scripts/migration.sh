#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <migration_name>"
  echo "Example: $0 fix_name_test"
  exit 1
fi

cd ..

MIGRATION_NAME=$1

MIGRATIONS_DIR="./migrations"

mkdir -p "$MIGRATIONS_DIR"

goose create -dir "$MIGRATIONS_DIR" "$MIGRATION_NAME" sql
