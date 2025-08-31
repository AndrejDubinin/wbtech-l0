-- +goose Up
-- +goose StatementBegin
CREATE TABLE delivery (
  id SERIAL PRIMARY KEY,
  order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
  name TEXT,
  phone TEXT,
  zip TEXT,
  city TEXT,
  address TEXT,
  region TEXT,
  email TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS delivery CASCADE;
-- +goose StatementEnd
