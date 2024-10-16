DROP INDEX CONCURRENTLY IF EXISTS "feedback_processed_at_id_idx";
DROP INDEX CONCURRENTLY IF EXISTS "feedback_translated_at_id_idx";
DROP INDEX CONCURRENTLY IF EXISTS "feedback_collected_at_id_idx";
DROP INDEX CONCURRENTLY IF EXISTS "feedback_posted_at_id_idx";
DROP INDEX CONCURRENTLY IF EXISTS "feedback_hash_idx";
DROP INDEX CONCURRENTLY IF EXISTS "feedback_product_id_idx";

DROP TABLE IF EXISTS "feedback";

DROP INDEX CONCURRENTLY IF EXISTS "issue_relevance_idx";
DROP INDEX CONCURRENTLY IF EXISTS "issue_assignee_id_idx";
-- DROP INDEX CONCURRENTLY IF EXISTS "issue_embedding_idx";
DROP INDEX CONCURRENTLY IF EXISTS "issue_product_id_idx";

DROP TABLE IF EXISTS "issue";

DROP INDEX CONCURRENTLY IF EXISTS "issue_feedback_feedback_id_idx";
DROP INDEX CONCURRENTLY IF EXISTS "issue_feedback_issue_id_idx";

DROP TABLE IF EXISTS "issue_feedback";

DROP TABLE IF EXISTS "partial_issue";

DROP INDEX CONCURRENTLY IF EXISTS "suggestion_relevance_idx";
DROP INDEX CONCURRENTLY IF EXISTS "suggestion_assignee_id_idx";
-- DROP INDEX CONCURRENTLY IF EXISTS "suggestion_embedding_idx";
DROP INDEX CONCURRENTLY IF EXISTS "suggestion_product_id_idx";

DROP TABLE IF EXISTS "suggestion";

DROP INDEX CONCURRENTLY IF EXISTS "suggestion_feedback_feedback_id_idx";
DROP INDEX CONCURRENTLY IF EXISTS "suggestion_feedback_suggestion_id_idx";

DROP TABLE IF EXISTS "suggestion_feedback";

DROP TABLE IF EXISTS "partial_suggestion";

DROP INDEX CONCURRENTLY IF EXISTS "review_product_id_idx";

DROP TABLE IF EXISTS "review";
