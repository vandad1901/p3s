mod auth "apps/auth"

set shell := ["sh", "-cu"]

@default:
    just --list

@start:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        --env-file apps/auth/.env \
        up -d --build --remove-orphans

@stop:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        -f ./infra/compose/docker-compose.yml \
        down

@compose-dev:
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --env-file apps/auth/.env \
        up -d --build --remove-orphans

@db:
    just auth db-reset

@generate:
    buf generate

@generate-secrets:
    openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out apps/auth/jwt_private_key.pem
