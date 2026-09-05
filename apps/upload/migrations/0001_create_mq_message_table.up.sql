CREATE TABLE IF NOT EXISTS mq_message(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    claim_id UUID NOT NULL,
    ---
    exchange_key varchar(255) NOT NULL,
    routing_key varchar(255) NOT NULL,
    message_body bytea NOT NULL,
    ---
    queued_at timestamptz NOT NULL,
    last_tried_at timestamptz NOT NULL,
    attempts int NOT NULL,
    in_transit boolean NOT NULL,
);

