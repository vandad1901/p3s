set shell := ["sh", "-cu"]
set dotenv-load := true

PGHOST := env_var_or_default("PG_HOST", "localhost")
PGPORT := env_var_or_default("PG_PORT", "5432")
PGUSER := env_var_or_default("PG_USER", "postgres")
PGPASSWORD := env_var_or_default("PG_PASSWORD", "postgres")
PGDATABASE := env_var_or_default("PG_DATABASE", "purpl3shadow")
MIGRATE := "~/.asdf/installs/golang/1.26.3/bin/migrate"
DATABASE_URL := "postgres://" + PGUSER + ":" + PGPASSWORD + "@" + PGHOST + ":" + PGPORT + "/" + PGDATABASE + "?sslmode=disable"

@default:
    just --list

@db-reset: && db-up
    echo "[INFO] dropping database {{ PGDATABASE }} at {{ PGHOST }}:{{ PGPORT }} with user {{ PGUSER }}"
    PGPASSWORD="{{ PGPASSWORD }}" dropdb --if-exists -h "{{ PGHOST }}" -p "{{ PGPORT }}" -U "{{ PGUSER }}" "{{ PGDATABASE }}"
    echo "[INFO] creating database {{ PGDATABASE }} at {{ PGHOST }}:{{ PGPORT }} with user {{ PGUSER }}"
    PGPASSWORD="{{ PGPASSWORD }}" createdb -h "{{ PGHOST }}" -p "{{ PGPORT }}" -U "{{ PGUSER }}" "{{ PGDATABASE }}"

@db-up:
    {{ MIGRATE }} -database "{{ DATABASE_URL }}" -path apps/api/migrations up

@db-down:
    {{ MIGRATE }} -database "{{ DATABASE_URL }}" -path apps/api/migrations down

alias compose := compose-up

@compose-up:
    docker compose -f ./infra/compose/docker-compose.yml up -d

@compose-down:
    docker compose -f ./infra/compose/docker-compose.yml down

@generate:
    protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto \
        --proto_path=./protobuf \
        --ts_proto_opt=onlyTypes=true \
        --ts_proto_out=./apps/lib/types \
        --go_out=./apps/api/gen \
        protobuf/**/*.proto
