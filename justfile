mod auth "apps/auth"

set shell := ["sh", "-cu"]

@default:
    just --list

@start:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        up -d --remove-orphans

@build:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        up -d --build --remove-orphans

@stop:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        down

@dev:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --env-file apps/auth/.env \
        up -d --build --remove-orphans

@compose-exec *ARGS:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        exec {{ ARGS }}

@db:
    just auth db-reset

@generate:
    buf generate
    buf build contracts \
        -o apps/envoy/descriptor.pb

@generate-secrets:
    openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out apps/auth/jwt_private_key.pem
