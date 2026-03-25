-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  uptdated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  email VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
