#!/bin/bash

set -e

source .env

PSQL_BASE=(psql -h "localhost" -p "5432" -d postgres -qtAX)

echo "Checking if DB role '$DB_USER' exists..."
if ! "${PSQL_BASE[@]}" -c "SELECT 1 FROM pg_roles WHERE rolname = '$DB_USER';" | grep -q 1; then
  echo "Creating DB role '$DB_USER'..."
  "${PSQL_BASE[@]}" -c "CREATE ROLE $DB_USER WITH LOGIN PASSWORD '$DB_PASSWORD' SUPERUSER;"
fi

echo "Checking if database '$DB_NAME' exists..."
if ! "${PSQL_BASE[@]}" -c "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME';" | grep -q 1; then
  echo "Creating and loading DB..."
  dbmate create
  dbmate load
  "${PSQL_BASE[@]}" -c "ALTER ROLE $DB_USER NOSUPERUSER;"
fi

echo "Checking if admin user '$ADMIN_USER' exists in DB..."
if ! psql "$DATABASE_URL" -qtAX -c "SELECT 1 FROM users WHERE username = '$ADMIN_USER';" | grep -q 1; then
  echo "Inserting admin user '$ADMIN_USER'..."
  psql "$DATABASE_URL" -c "INSERT INTO users (username, hashed_password, created_at) VALUES ('$ADMIN_USER', '$ADMIN_PASSWORD', NOW());"
fi

echo "Database ready."
