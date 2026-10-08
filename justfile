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
        --profile apps up -d --remove-orphans --wait

@build:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --profile apps up -d --remove-orphans --wait --build

@stop:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --profile apps down

@dev:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        up -d --remove-orphans --wait --build

@compose-exec *ARGS:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
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
    just auth run & \
    sleep 0.5 && just api run & \
    sleep 0.5 && just upload run & \
    sleep 0.5 && just media run & \
    wait

@test:
    go test ./packages/go/...
    just --dotenv-filename .env.test auth test
    just --dotenv-filename .env.test api test
    just --dotenv-filename .env.test upload test
    just --dotenv-filename .env.test media test

@generate:
    buf generate
    buf build contracts \
        -o apps/envoy/descriptor.pb

@generate-secrets:
    sudo mkdir -p /etc/p3s
    sudo openssl genpkey \
    -algorithm EC \
    -pkeyopt ec_paramgen_curve:P-256 \
    -out /etc/p3s/jwt_private_key.pem

    sudo chown 65532:65532 /etc/p3s/jwt_private_key.pem
    sudo chmod 400 /etc/p3s/jwt_private_key.pem
    sudo chmod 700 /etc/p3s
