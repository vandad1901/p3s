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
        --env-file apps/media/.env \
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

@db-up:
    just auth db-up
    just api db-up
    just upload db-up
    just media db-up

@db:
    just auth db-reset
    just api db-reset
    just upload db-reset
    just media db-reset

@run:
    APP_ENV=development just auth run & \
    APP_ENV=development just api run & \
    APP_ENV=development just upload run & \
    APP_ENV=development just media run & \
    wait

@test:
    go test ./packages/go/...
    APP_ENV=test just --dotenv-filename .env.test auth test
    APP_ENV=test just --dotenv-filename .env.test api test
    APP_ENV=test just --dotenv-filename .env.test upload test
    APP_ENV=test just --dotenv-filename .env.test media test

@generate:
    buf generate
    buf build contracts \
        -o apps/envoy/descriptor.pb

@generate-secrets:
    openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out apps/auth/jwt_private_key.pem
