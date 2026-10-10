mod auth "apps/auth"
mod api "apps/api"
mod upload "apps/upload"
mod media "apps/media"
mod frontend "apps/frontend"

set shell := ["sh", "-cu"]

@default:
    just --list

@prod:
    docker compose \
        -f ./infra/compose/docker-compose.prod.yml \
        --env-file .env.production up -d --build --remove-orphans --wait

@prod-down:
    docker compose \
        -f ./infra/compose/docker-compose.prod.yml \
        down

@generate-secrets-prod:
    sudo mkdir -p /etc/p3s
    sudo openssl genpkey \
    -algorithm EC \
    -pkeyopt ec_paramgen_curve:P-256 \
    -out /etc/p3s/jwt_private_key.pem

    sudo chown 65532:65532 /etc/p3s/jwt_private_key.pem
    sudo chmod 400 /etc/p3s/jwt_private_key.pem
    sudo chmod 700 /etc/p3s

    sed -i.bak \
        -e "s|PG_ADMIN_PASSWORD=.*|PG_ADMIN_PASSWORD='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        -e "s|AUTH_PG_PASSWORD=.*|AUTH_PG_PASSWORD='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        -e "s|API_PG_PASSWORD=.*|API_PG_PASSWORD='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        -e "s|RMQ_PASSWORD=.*|RMQ_PASSWORD='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        -e "s|S3_ROOT_USERNAME=.*|S3_ROOT_USERNAME='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        -e "s|S3_ROOT_PASSWORD=.*|S3_ROOT_PASSWORD='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        -e "s|PGADMIN_EMAIL=.*|PGADMIN_EMAIL=admin@admin.com|" \
        -e "s|PGADMIN_PASSWORD=.*|PGADMIN_PASSWORD='$(openssl rand -base64 32 | tr -d "=+/")'|" \
        .env.production

@start: skip-prod
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --profile apps up -d --remove-orphans --wait

@build: skip-prod
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --profile apps up -d --remove-orphans --wait --build

@stop: skip-prod
    docker compose \
        -f ./infra/compose/docker-compose.dev.yml \
        --profile apps down

@dev: skip-prod
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

@db: skip-prod
    just auth db-reset
    just api db-reset
    just upload db-reset
    just media db-reset

@run: skip-prod
    just auth run & \
    sleep 0.5 && just api run & \
    sleep 0.5 && just upload run & \
    sleep 0.5 && just media run & \
    wait

@test: skip-prod
    go test ./packages/go/...
    just --dotenv-filename .env.test auth test
    just --dotenv-filename .env.test api test
    just --dotenv-filename .env.test upload test
    just --dotenv-filename .env.test media test

@generate:
    buf generate
    buf build contracts \
        -o apps/envoy/descriptor.pb

@generate-secrets-ci: skip-prod
    openssl genpkey \
    -algorithm EC \
    -pkeyopt ec_paramgen_curve:P-256 \
    -out apps/auth/jwt_private_key.pem

@skip-prod:
    if [ -f .env.production ]; then \
        echo ".env.production exists, refusing to continue"; \
        exit 1; \
    fi
