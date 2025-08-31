-- +goose Up
-- +goose StatementBegin
CREATE TABLE payment (
  id SERIAL PRIMARY KEY,
  order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
  transaction TEXT,
  request_id TEXT,
  currency TEXT,
  provider TEXT,
  amount NUMERIC,
  payment_dt BIGINT,
  bank TEXT,
  delivery_cost NUMERIC,
  goods_total NUMERIC,
  custom_fee NUMERIC
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment CASCADE;
-- +goose StatementEnd
