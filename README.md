# P3S Blog Platform

A blogging platform written in Go. Build from the ground up with distributed system principles for reliability, scalability, and security. This codebase is a showcase of modern Go practices, including:

- gRPC services behind an Envoy gateway (gRPC-JSON transcoding, CORS, routing)
- Protobuf contracts as the single source of truth, generating Go and TypeScript
- Short-lived ES256 JWTs + refresh-token sessions, verified locally via a JWKS endpoint
- Presigned direct-to-storage uploads + a dedicated, idempotent media-ingest worker
- Transactional outbox (Postgres) feeding RabbitMQ with at-least-once delivery
- Gherkin/Godog integration tests against real PostgreSQL, run in CI
- A React 19 frontend with a drag-and-drop block editor, markdown rendering, and LaTeX rendering
- asynchronous, decoupled services with independent deployments

## Architecture

Protobuf contracts are the single source of truth. Services are independent gRPC apps that validate access tokens locally against the auth JWKS and coordinate asynchronously over an event bus:

```
web ──▶ Envoy ──┬──▶ Auth: identity, JWT, sessions, JWKS
 │               ├──▶ API: post CRUD (Postgres)
 │               └──▶ Upload: presigned URLs ──▶ outbox ──▶ RabbitMQ
 │                                                                │
 │                                                                ▼
 └── presigned upload ─────▶ Silo (S3-compatible) ◀─────────  Media
                            originals + derivatives    validate, resize, DB lease
```

### How an upload works

1. The client asks `upload` for a presigned POST (20 MiB cap, 5 minute expiry, keys namespaced by user) and sends the file straight to Bucket, so file bytes never pass through the services.
2. The client calls `FinalizeUpload`; `upload` checks the object exists and writes its key to a transactional outbox in Postgres.
3. An outbox worker claims in batches with `SKIP LOCKED`, publishes them to RabbitMQ at-least-once, and recovers messages stuck in transit.
4. `media` takes a per-file DB lease, so redelivery is safe: already-processed messages are ACKed, permanent failures are discarded, transient ones are retried (up to 5 tries).
5. It validates the image (jpeg, png, gif or webp; at most 10,000 px per side and 40M pixels) and stores resized derivatives at 320, 640 and 1280 px wide. Derivatives are capped by the source resolution, so images are never upscaled.

## Status

| Component                                           | Status  |
| --------------------------------------------------- | ------- |
| Proto contracts + code-gen pipeline                 | 🟢 done |
| Add Envoy for gRPC and REST transcoding             | 🟢 done |
| Auth service (JWT + refresh tokens, JWKS, sessions) | 🟢 done |
| API service (post CRUD, media-ready)                | 🟢 done |
| API and Auth domain tests with Gherkin              | 🟢 done |
| Upload service + Silo                               | 🟢 done |
| media worker                                        | 🟢 done |
| Wire up worker and upload with RabbitMQ             | 🟢 done |
| Web frontend (auth, home, viewer, block editor)     | 🟢 done |

## Running in production

Requires Docker, Docker Compose and `just`.

```bash
cp .env.production.example .env.production   # then edit the non-secret values
just generate-secrets-prod                   # signing key + random passwords
just prod-build
just prod                                    # database migrations run automatically
```

On the first deploy only, create the upload bucket and allow anonymous reads of its objects (derivatives are served straight from Silo; uploads use presigned URLs and need no public write access):

```bash
set -a; . ./.env.production; set +a
export AWS_ACCESS_KEY_ID="$S3_ROOT_USERNAME" AWS_SECRET_ACCESS_KEY="$S3_ROOT_PASSWORD" AWS_PAGER=""

aws s3api create-bucket --bucket p3s-upload-bucket --endpoint-url "http://$SILO_HOST_IP:9000"

aws s3api put-bucket-policy --bucket p3s-upload-bucket --endpoint-url "http://$SILO_HOST_IP:9000" \
  --policy '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"arn:aws:s3:::p3s-upload-bucket/*"}]}'
```

## Repository layout

```
contracts/       protobuf contracts (auth/, api/, upload/, media/)
packages/go/     shared packages (dbpattern, envutil, idv, usercontext, …)
packages/web/    generated TypeScript types
apps/            auth · api · upload · media (Go services)
apps/frontend/   React frontend
apps/envoy/      Envoy gateway config
infra/compose/   Docker Compose environments (dev and prod)
```

## License

MIT - see [LICENSE](LICENSE).
