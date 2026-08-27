mod auth "apps/auth"
mod api "apps/api"

set shell := ["sh", "-cu"]

@default:
    just --list

@start:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        up -d --remove-orphans

@build:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        up -d --build --remove-orphans

@stop:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        down

@dev:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        up -d --build --remove-orphans

@compose-exec *ARGS:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        exec {{ ARGS }}

@db:
    just auth db-reset
    just api db-reset

@generate:
    buf generate
    buf build contracts \
        -o apps/envoy/descriptor.pb

@generate-secrets:
    openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out apps/auth/jwt_private_key.pem
