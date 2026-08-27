# Blog Platform

Go microservices + Angular frontend. A blogging platform where posts are composed from blocks of text and media. Media is handled by a dedicated pipeline: uploads land in MinIO (S3-compatible), and RabbitMQ-driven workers process them asynchronously (thumbnails, normalization).

## Architecture

Protobuf contracts in `contracts/` are the single source of truth (generated into `packages/go/gen` + `packages/web/gen`). Services are independent gRPC apps that validate access tokens locally against the auth JWKS and coordinate asynchronously over an event bus:

```
web ──▶ Envoy ──▶ auth (identity, JWT, sessions)
            └────▶ api (post CRUD, postgres)
            └────▶ upload-service ──▶ MinIO
                          │
                          │ publish
                          ▼
                     RabbitMQ ──▶ media-ingest (thumbnails, normalize)
```

## Status

| Component                                                                  | Status        |
| -------------------------------------------------------------------------- | ------------- |
| Auth service (JWT + refresh tokens, JWKS, sessions)                        | 🟢 done       |
| API service (post CRUD, gRPC + Envoy REST, media-ready `post_block` model) | 🟢 done       |
| Proto contracts + codegen pipeline                                         | 🟢 done       |
| Compose environment (`auth`, `api`, `pg`, `envoy`)                         | 🟢 done       |
| Upload service + MinIO                                                     | 🔨 scaffolded |
| RabbitMQ bus + media-ingest worker                                         | ⬜ planned    |
| Web editor (pending-upload gating)                                         | ⬜ planned    |
| Public post routes + cross-posting                                         | ⬜ planned    |

## Roadmap

1. **Upload service** — authenticated multipart upload to MinIO; client-defined object keys (`userID/uuidv7`) prevent cross-user collisions. Add MinIO to Compose.
2. **Event bus** — RabbitMQ in Compose; upload service publishes upload events after object + record are durable.
3. **Media-ingest worker** — consumes upload events, generates thumbnails (optionally strips EXIF / normalizes formats), stores artifacts back in MinIO; DLQ + ack-after-durable for at-least-once delivery.
4. **API/web integration** — wire media blocks into post create/read; public post pages; gate draft/publish until all uploads succeed.
5. **Extras** — cross-posting consumers (e.g. Telegram), scaling, observability.

## Getting Started

1. Clone, then generate secrets once: `just generate-secrets`
2. Regenerate contracts after proto changes: `just generate`
3. Run the environment: `just build` (or `just dev`); API at `http://localhost:8080` (Envoy), web at `http://localhost:4200`.
4. Per-service commands: `just auth dev`, `just api dev`, etc. (each app has its own `justfile`).

## Repository layout

```
contracts/       protobuf contracts (auth/, api/)
packages/go/     shared packages (dbpattern, envutil, idv, usercontext, …)
packages/web/    generated TypeScript types
apps/            auth · api · upload-service · envoy
infra/compose/   Docker Compose environment
```

## Contributing

Contributions are welcome! Open an issue or pull request on GitHub.

## License

MIT - see [LICENSE](LICENSE).
