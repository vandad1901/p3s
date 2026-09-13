# P3S Blog Platform

A blogging platform written in Go. Build from the ground up with distributed system principles for reliability, scalability, and security. This codebase is a showcase of modern Go practices, including:

- gRPC + Envoy REST gateway
- Protobuf contracts + code-gen pipeline
- JWT + refresh tokens, JWKS
- Bucketed media uploads + dedicated media-ingest worker
- asynchronous, decoupled services with independent deployments

## Architecture

Protobuf contracts are the single source of truth. Services are independent gRPC apps that validate access tokens locally against the auth JWKS and coordinate asynchronously over an event bus:

```
web ──▶ Envoy ──▶ auth (identity, JWT, sessions)
            └────▶ api (post CRUD, postgres)
            └────▶ upload ──▶ Silo
                          │
                          │ publish to RabbitMQ
                          │
                          └──▶ media-ingest (thumbnails, normalize)
```

## Status

| Component                                           | Status     |
| --------------------------------------------------- | ---------- |
| Proto contracts + code-gen pipeline                 | 🟢 done    |
| Add Envoy for gRPC and REST transcoding             | 🟢 done    |
| Auth service (JWT + refresh tokens, JWKS, sessions) | 🟢 done    |
| API service (post CRUD, media-ready )               | 🟢 done    |
| API and Auth domain tests with Gherkin              | 🟢 done    |
| Upload service + Silo                               | 🟢 done    |
| media-ingest worker                                 | 🟢 done    |
| wire up worker and upload with rabbitMQ             | 🟢 done    |
| Web editor (pending-upload gating)                  | ⬜ planned |

## Getting Started

0. requirements: `just`, `docker`, `docker-compose` and `mise`. You can use the scripts in `scripts/` to install them on Linux.
1. Clone the repo

```bash
git clone github.com/vandad1901/p3s
```

2. The only real dependencies are `docker`, `docker compose` and `mise`. `mise` will take care of the rest. You can use the scripts in `scripts/` to install them on Linux.

```bash
./scripts/dev/00_install_essentials.sh
./scripts/dev/01_install_mise.sh
./scripts/dev/02_install_docker.sh
```

3. Start the service dependencies (Postgres, RabbitMQ, Silo, etc.)

```bash
just dev
```

4. Initialize the database and buckets and run database migrations

```bash
just db-up
```

5. Run the services

```bash
just run
```

the API will be accessible at `http://localhost:8080` , and the frontend at `http://localhost:3000` (upcoming). You can use the `just` commands to run the services individually, or use `docker compose` directly.

## Repository layout

```
contracts/       protobuf contracts (auth/, api/)
packages/go/     shared packages (dbpattern, envutil, idv, usercontext, …)
packages/web/    generated TypeScript types
apps/            auth · api · upload · media
infra/compose/   Docker Compose environment
```

## Contributing

Contributions are welcome! Open an issue or pull request on GitHub.

## License

MIT - see [LICENSE](LICENSE).
