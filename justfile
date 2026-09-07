mod auth "apps/auth"
mod api "apps/api"
mod upload "apps/upload"
mod media "apps/media"

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
        --env-file apps/media/.env \
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
        --env-file apps/media/.env \
        down

@dev:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        --env-file apps/media/.env \
        up -d --build --remove-orphans

@compose-exec *ARGS:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        --env-file apps/api/.env \
        --env-file apps/upload/.env \
        --env-file apps/media/.env \
        exec {{ ARGS }}

@db:
    just auth db-reset
    just api db-reset
    just upload db-reset

@run:
    just auth run & \
    just api run & \
    just upload run & \
    just media run & \
    wait

@test:
    just --dotenv-filename .env.test auth test
    just --dotenv-filename .env.test api test
    just --dotenv-filename .env.test upload test
    just --dotenv-filename .env.test media test

@generate:
    buf generate
    buf build contracts \
        -o apps/envoy/descriptor.pb

@generate-secrets:
    openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out apps/auth/jwt_private_key.pem
