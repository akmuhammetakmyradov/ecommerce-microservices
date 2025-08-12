#!/bin/bash
set -e

# Defaults (can be overridden via environment variables)
DB_USER="${DB_USER:-user_stocks}"
DB_PASS="${DB_PASS:-stocks1234}"
DB_NAME="${DB_NAME:-stocks}"

echo "🔧 Initializing PostgreSQL with database: $DB_NAME and user: $DB_USER"

# Drop role if exists, create role
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  DROP ROLE IF EXISTS $DB_USER;
  CREATE ROLE $DB_USER LOGIN PASSWORD '$DB_PASS';
EOSQL

# Drop database if exists
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  DROP DATABASE IF EXISTS $DB_NAME;
EOSQL

# Create database (must be in its own call, no other commands)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  CREATE DATABASE $DB_NAME
    WITH OWNER = $DB_USER
         ENCODING = 'UTF8'
         CONNECTION LIMIT = -1
         IS_TEMPLATE = FALSE;
EOSQL

# Grant privileges inside the new database
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB_NAME" <<-EOSQL
  GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
  GRANT USAGE ON SCHEMA public TO $DB_USER;
  GRANT CREATE ON SCHEMA public TO $DB_USER;
EOSQL

echo "✅ PostgreSQL init script completed."
