#!/bin/bash

set -e

source .env

echo "Checking DB role..."
ROLE_EXISTS=$(psql -U $DB_NAME -tc "SELECT 1 FROM pg_roles WHERE rolname = '$DB_USER';" | grep -q 1 && echo yes || echo no)
if [ "$ROLE_EXISTS" = "no" ]; then
  echo "Creating DB role..."
  psql -U $DB_NAME -c "CREATE ROLE $DB_USER WITH LOGIN PASSWORD '$DB_PASSWORD' SUPERUSER;"
fi

echo "Checking database..."
DB_EXISTS=$(psql -U $DB_NAME -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME';" | grep -q 1 && echo yes || echo no)
if [ "$DB_EXISTS" = "no" ]; then
  echo "Creating and loading DB..."
  dbmate create
  dbmate load
  psql -U $DB_NAME -c "ALTER ROLE $DB_USER NOSUPERUSER;"
fi

echo "Checking admin user..."
ADMIN_EXISTS=$(psql "$DATABASE_URL" -tc "SELECT 1 FROM users WHERE username = '$DB_USER';" | grep -q 1 && echo yes || echo no)
if [ "$ADMIN_EXISTS" = "no" ]; then
  echo "Inserting admin user..."
  psql "$DATABASE_URL" -c "INSERT INTO users (username, hashed_password, created_at) VALUES ('$ADMIN_USER', '$ADMIN_PASSWORD', NOW());"
fi

echo "Database ready."
