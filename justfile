set shell := ["sh", "-cu"]
set dotenv-load := true

MIGRATE := "~/.asdf/installs/golang/1.26.3/bin/migrate"
AUTHPREFIX := "AUTH_"
AUTH_PGUSER := env_var_or_default(AUTHPREFIX + "PG_USER", "postgres")
AUTH_PGPASSWORD := env_var_or_default(AUTHPREFIX + "PG_PASSWORD", "postgres")
AUTH_PGHOST := env_var_or_default(AUTHPREFIX + "PG_HOST", "localhost")
AUTH_PGPORT := env_var_or_default(AUTHPREFIX + "PG_PORT", "5432")
AUTH_PGDATABASE := env_var_or_default(AUTHPREFIX + "PG_DATABASE", "purpl3shadow")
AUTH_DATABASE_URL := "postgres://" + AUTH_PGUSER + ":" + AUTH_PGPASSWORD + "@" + AUTH_PGHOST + ":" + AUTH_PGPORT + "/" + AUTH_PGDATABASE + "?sslmode=disable"

@default:
    just --list

@db-reset: && db-up
    echo "[INFO] dropping database {{ AUTH_PGDATABASE }} at {{ AUTH_PGHOST }}:{{ AUTH_PGPORT }} with user {{ AUTH_PGUSER }}"
    PGPASSWORD="{{ AUTH_PGPASSWORD }}" dropdb --if-exists -h "{{ AUTH_PGHOST }}" -p "{{ AUTH_PGPORT }}" -U "{{ AUTH_PGUSER }}" "{{ AUTH_PGDATABASE }}"
    echo "[INFO] creating database {{ AUTH_PGDATABASE }} at {{ AUTH_PGHOST }}:{{ AUTH_PGPORT }} with user {{ AUTH_PGUSER }}"
    PGPASSWORD="{{ AUTH_PGPASSWORD }}" createdb -h "{{ AUTH_PGHOST }}" -p "{{ AUTH_PGPORT }}" -U "{{ AUTH_PGUSER }}" "{{ AUTH_PGDATABASE }}"

@db-up:
    {{ MIGRATE }} -database "{{ AUTH_DATABASE_URL }}" -path apps/api/migrations up

@db-down:
    {{ MIGRATE }} -database "{{ AUTH_DATABASE_URL }}" -path apps/api/migrations down

alias compose := compose-up

@compose-up:
    docker compose -f ./infra/compose/docker-compose.yml up -d

@compose-down:
    docker compose -f ./infra/compose/docker-compose.yml down

@generate:
    buf generate
