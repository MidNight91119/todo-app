-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2026-08-26T10:07:49.565Z

CREATE TABLE "tasks" (
  "id" bigserial PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "completed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);
