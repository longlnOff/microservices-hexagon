-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE "product_status" AS ENUM (
  'ENABLED',
  'DISABLED'
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "status" product_status NOT NULL DEFAULT 'ENABLED',
  "category_id" uuid NOT NULL,
  "name" varchar(128) NOT NULL,
  "provider" varchar(64) NOT NULL,
  "serving_endpoint" varchar(256) NOT NULL,
  "api_reference" varchar(256) NOT NULL,
  "short_introduction" json NOT NULL,
  "long_introduction" json NOT NULL,
  "usecases" json NOT NULL,
  "is_free_trial" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_when" timestamptz NOT NULL DEFAULT (now()),
  "created_by" varchar(64) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
  "modified_by" varchar(64) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "products";
DROP TYPE IF EXISTS "product_status";
-- +goose StatementEnd
