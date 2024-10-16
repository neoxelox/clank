CREATE TABLE IF NOT EXISTS "feedback" (
    "id" VARCHAR(20) PRIMARY KEY,
    "product_id" VARCHAR(20) NOT NULL,
    "hash" VARCHAR(40) UNIQUE NOT NULL,
    "source" VARCHAR(50) NOT NULL,
    "customer" JSONB NOT NULL,
    "content" TEXT NOT NULL,
    "language" VARCHAR(50) NOT NULL,
    "translation" TEXT NOT NULL,
    "release" VARCHAR(50) NOT NULL,
    "metadata" JSONB NOT NULL,
    "tokens" BIGINT NOT NULL,
    "posted_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "collected_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "translated_at" TIMESTAMP WITH TIME ZONE NULL,
    "processed_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "feedback_product_id_idx" ON "feedback" ("product_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "feedback_hash_idx" ON "feedback" ("hash");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "feedback_posted_at_id_idx" ON "feedback" ("posted_at", "id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "feedback_collected_at_id_idx" ON "feedback" ("collected_at", "id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "feedback_translated_at_id_idx" ON "feedback" ("translated_at", "id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "feedback_processed_at_id_idx" ON "feedback" ("processed_at", "id");

CREATE TABLE IF NOT EXISTS "issue" (
    "id" VARCHAR(20) PRIMARY KEY,
    "product_id" VARCHAR(20) NOT NULL,
    "embedding" VECTOR(1536) NOT NULL,
    "sources" JSONB NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "description" TEXT NOT NULL,
    "steps" TEXT[] NOT NULL,
    "severities" JSONB NOT NULL,
    "priority" BIGINT NOT NULL,
    "categories" JSONB NOT NULL,
    "releases" JSONB NOT NULL,
    "customers" BIGINT NOT NULL,
    "assignee_id" VARCHAR(20) NULL,
    "quality" BIGINT NULL,
    "first_seen_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "last_seen_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE NULL,
    "last_aggregated_at" TIMESTAMP WITH TIME ZONE NULL,
    "exported_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "issue_product_id_idx" ON "issue" ("product_id");
-- Don't create an index for the embeddings until it can be composited with the product id: https://github.com/pgvector/pgvector/issues/259
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS "issue_embedding_idx" ON "issue" USING hnsw ("embedding" vector_cosine_ops)
--     WITH (m = 16, ef_construction = 64);
CREATE INDEX CONCURRENTLY IF NOT EXISTS "issue_assignee_id_idx" ON "issue" ("assignee_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "issue_relevance_idx" ON "issue" ("priority", "last_seen_at", "id");

CREATE TABLE IF NOT EXISTS "issue_feedback" (
    "issue_id" VARCHAR(20) NOT NULL REFERENCES "issue" ("id") ON DELETE CASCADE,
    "feedback_id" VARCHAR(20) NOT NULL REFERENCES "feedback" ("id") ON DELETE CASCADE,
    PRIMARY KEY ("issue_id", "feedback_id")
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "issue_feedback_issue_id_idx" ON "issue_feedback" ("issue_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "issue_feedback_feedback_id_idx" ON "issue_feedback" ("feedback_id");

CREATE TABLE IF NOT EXISTS "partial_issue" (
    "id" VARCHAR(20) PRIMARY KEY,
    "feedback_id" VARCHAR(20) NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "description" TEXT NOT NULL,
    "steps" TEXT[] NOT NULL,
    "severity" VARCHAR(50) NOT NULL,
    "category" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "partial_issue_created_at_id_idx" ON "partial_issue" ("created_at", "id");

CREATE TABLE IF NOT EXISTS "suggestion" (
    "id" VARCHAR(20) PRIMARY KEY,
    "product_id" VARCHAR(20) NOT NULL,
    "embedding" VECTOR(1536) NOT NULL,
    "sources" JSONB NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "description" TEXT NOT NULL,
    "reason" TEXT NOT NULL,
    "importances" JSONB NOT NULL,
    "priority" BIGINT NOT NULL,
    "categories" JSONB NOT NULL,
    "releases" JSONB NOT NULL,
    "customers" BIGINT NOT NULL,
    "assignee_id" VARCHAR(20) NULL,
    "quality" BIGINT NULL,
    "first_seen_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "last_seen_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE NULL,
    "last_aggregated_at" TIMESTAMP WITH TIME ZONE NULL,
    "exported_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "suggestion_product_id_idx" ON "suggestion" ("product_id");
-- Don't create an index for the embeddings until it can be composited with the product id: https://github.com/pgvector/pgvector/issues/259
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS "suggestion_embedding_idx" ON "suggestion" USING hnsw ("embedding" vector_cosine_ops)
--     WITH (m = 16, ef_construction = 64);
CREATE INDEX CONCURRENTLY IF NOT EXISTS "suggestion_assignee_id_idx" ON "suggestion" ("assignee_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "suggestion_relevance_idx" ON "suggestion" ("priority", "last_seen_at", "id");

CREATE TABLE IF NOT EXISTS "suggestion_feedback" (
    "suggestion_id" VARCHAR(20) NOT NULL REFERENCES "suggestion" ("id") ON DELETE CASCADE,
    "feedback_id" VARCHAR(20) NOT NULL REFERENCES "feedback" ("id") ON DELETE CASCADE,
    PRIMARY KEY ("suggestion_id", "feedback_id")
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "suggestion_feedback_suggestion_id_idx" ON "suggestion_feedback" ("suggestion_id");
CREATE INDEX CONCURRENTLY IF NOT EXISTS "suggestion_feedback_feedback_id_idx" ON "suggestion_feedback" ("feedback_id");

CREATE TABLE IF NOT EXISTS "partial_suggestion" (
    "id" VARCHAR(20) PRIMARY KEY,
    "feedback_id" VARCHAR(20) NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "description" TEXT NOT NULL,
    "reason" TEXT NOT NULL,
    "importance" VARCHAR(50) NOT NULL,
    "category" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "partial_suggestion_created_at_id_idx" ON "partial_suggestion" ("created_at", "id");

CREATE TABLE IF NOT EXISTS "review" (
    "id" VARCHAR(20) PRIMARY KEY,
    "product_id" VARCHAR(20) NOT NULL,
    "feedback_id" VARCHAR(20) NOT NULL,
    "keywords" VARCHAR(100)[] NOT NULL,
    "sentiment" VARCHAR(50) NOT NULL,
    "emotions" VARCHAR(50)[] NOT NULL,
    "intention" VARCHAR(50) NOT NULL,
    "category" VARCHAR(50) NOT NULL,
    "quality" BIGINT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "exported_at" TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "review_product_id_idx" ON "review" ("product_id");
