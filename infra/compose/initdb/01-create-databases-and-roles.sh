#!/bin/sh
set -eu

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres \
  -v auth_user="$AUTH_PG_USER" -v auth_pw="$AUTH_PG_PASSWORD" -v auth_db="$AUTH_PG_DATABASE" \
  -v api_user="$API_PG_USER" -v api_pw="$API_PG_PASSWORD" -v api_db="$API_PG_DATABASE" <<'SQL'
CREATE ROLE :"auth_user" LOGIN PASSWORD :'auth_pw';
CREATE DATABASE :"auth_db" OWNER :"auth_user";
REVOKE CONNECT ON DATABASE :"auth_db" FROM PUBLIC;

CREATE ROLE :"api_user" LOGIN PASSWORD :'api_pw';
CREATE DATABASE :"api_db" OWNER :"api_user";
REVOKE CONNECT ON DATABASE :"api_db" FROM PUBLIC;
SQL
