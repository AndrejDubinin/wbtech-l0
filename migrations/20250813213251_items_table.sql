-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  id SERIAL PRIMARY KEY,
  order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id BIGINT,
  track_number TEXT,
  price NUMERIC,
  rid TEXT,
  name TEXT,
  sale NUMERIC,
  size TEXT,
  total_price NUMERIC,
  nm_id BIGINT,
  brand TEXT,
  status INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items CASCADE;
-- +goose StatementEnd
