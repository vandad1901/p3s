set shell := ["sh", "-cu"]
set dotenv-load := true

PGHOST := env_var_or_default("PGHOST", "localhost")
PGPORT := env_var_or_default("PGPORT", "5432")
PGUSER := env_var_or_default("PGUSER", "postgres")
PGPASSWORD := env_var_or_default("PGPASSWORD", "postgres")
PGDATABASE := env_var_or_default("PGDATABASE", "purpl3shadow")
MIGRATE := "~/.asdf/installs/golang/1.26.1/bin/migrate"
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

@compose-up:
    podman compose -f ./infra/compose/docker-compose.yml up -d 

@compose-down:
    podman compose -f ./infra/compose/docker-compose.yml down

@generate:
    protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto \
        --proto_path=./protobuf \
        --ts_proto_opt=onlyTypes=true \
        --ts_proto_out=./apps/lib/types \
        --go_out=./apps/api/gen \
        protobuf/**/*.proto
