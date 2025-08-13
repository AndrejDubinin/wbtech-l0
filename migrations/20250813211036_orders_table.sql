-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
  order_uid TEXT PRIMARY KEY,
  track_number TEXT,
  entry TEXT,
  locale TEXT,
  internal_signature TEXT,
  customer_id TEXT,
  delivery_service TEXT,
  shardkey TEXT,
  sm_id INT,
  date_created TIMESTAMPTZ,
  oof_shard TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders CASCADE;
-- +goose StatementEnd
