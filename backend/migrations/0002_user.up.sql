CREATE TABLE IF NOT EXISTS "user" (
    "id" VARCHAR(20) PRIMARY KEY,
	"organization_id" VARCHAR(20) NOT NULL,
	"name" VARCHAR(100) NOT NULL,
	"picture" VARCHAR(100) NOT NULL,
	"email" VARCHAR(100) UNIQUE NOT NULL,
	"role" VARCHAR(50) NOT NULL,
	"settings" JSONB NOT NULL,
	"created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
	"deleted_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "user_organization_id_idx" ON "user" ("organization_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "user_email_idx" ON "user" ("email");

CREATE TABLE IF NOT EXISTS "invitation" (
    "id" VARCHAR(20) PRIMARY KEY,
	"organization_id" VARCHAR(20) NOT NULL,
	"email" VARCHAR(100) UNIQUE NOT NULL,
	"role" VARCHAR(50) NOT NULL,
	"expires_at" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "invitation_organization_id_idx" ON "invitation" ("organization_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "invitation_email_idx" ON "invitation" ("email");

CREATE TABLE IF NOT EXISTS "session" (
    "id" VARCHAR(20) PRIMARY KEY,
    "user_id" VARCHAR(20) NOT NULL,
    "provider" VARCHAR(50) NOT NULL,
    "metadata" JSONB NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "last_seen_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "expired_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "session_user_id_idx" ON "session" ("user_id");

CREATE TABLE IF NOT EXISTS "sign_in_code" (
    "id" VARCHAR(20) PRIMARY KEY,
    "email" VARCHAR(100) UNIQUE NOT NULL,
    "code" VARCHAR(6) NOT NULL,
    "attempts" BIGINT NOT NULL,
    "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "sign_in_code_email_idx" ON "sign_in_code" ("email");

CREATE TABLE IF NOT EXISTS "organization" (
    "id" VARCHAR(20) PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "picture" VARCHAR(100) NOT NULL,
    "domain" VARCHAR(100) UNIQUE NOT NULL,
    "settings" JSONB NOT NULL,
    "plan" VARCHAR(50) NOT NULL,
    "trial_ends_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "capacity" JSONB NOT NULL,
    "usage" JSONB NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "organization_domain_idx" ON "organization" ("domain");
