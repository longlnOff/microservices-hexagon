-- +goose Up
-- +goose StatementBegin
-- First, create the UUID extension if it doesn't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Then create your table
CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "code" varchar(64) UNIQUE NOT NULL,
  "name" json NOT NULL,
  "description" varchar(64) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "modified_when" timestamptz NOT NULL DEFAULT (now()),
  "created_by" varchar(64) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
  "modified_by" varchar(64) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "categories";
-- We don't drop the UUID extension in down migration as other tables might be using it
-- +goose StatementEnd
