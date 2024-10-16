CREATE TABLE IF NOT EXISTS "product" (
    "id" VARCHAR(20) PRIMARY KEY,
    "organization_id" VARCHAR(20) NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "picture" VARCHAR(100) NOT NULL,
    "language" VARCHAR(50) NOT NULL,
    "context" TEXT NOT NULL,
    "categories" VARCHAR(50)[] NOT NULL,
    "release" VARCHAR(50) NOT NULL,
    "settings" JSONB NOT NULL,
    "usage" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "product_organization_id_idx" ON "product" ("organization_id");

CREATE TABLE IF NOT EXISTS "collector" (
    "id" VARCHAR(20) PRIMARY KEY,
    "product_id" VARCHAR(20) NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "settings" JSONB NOT NULL,
    "jobdata" JSONB NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "collector_product_id_idx" ON "collector" ("product_id");

CREATE TABLE IF NOT EXISTS "exporter" (
    "id" VARCHAR(20) PRIMARY KEY,
    "product_id" VARCHAR(20) NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "settings" JSONB NOT NULL,
    "jobdata" JSONB NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "exporter_product_id_idx" ON "exporter" ("product_id");
