-- +goose Up
-- +goose StatementBegin
CREATE TABLE "order_items" (
  "order_id" int NOT NULL,
  "product_code" string,
  "quantity" int,
  "unit_price" float,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY NOT NULL,
  "customer_id" int,
  "status" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE TABLE IF EXISTS "order_items";
DELETE TABLE IF EXISTS "orders";
-- +goose StatementEnd
