DROP INDEX CONCURRENTLY IF EXISTS "user_email_idx";
DROP INDEX CONCURRENTLY IF EXISTS "user_organization_id_idx";

DROP TABLE IF EXISTS "user";

DROP INDEX CONCURRENTLY IF EXISTS "invitation_email_idx";
DROP INDEX CONCURRENTLY IF EXISTS "invitation_organization_id_idx";

DROP TABLE IF EXISTS "invitation";

DROP INDEX CONCURRENTLY IF EXISTS "session_user_id_idx";

DROP TABLE IF EXISTS "session";

DROP INDEX CONCURRENTLY IF EXISTS "sign_in_code_email_idx";

DROP TABLE IF EXISTS "sign_in_code";

DROP INDEX CONCURRENTLY IF EXISTS "organization_domain_idx";

DROP TABLE IF EXISTS "organization";
