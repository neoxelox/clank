DROP INDEX CONCURRENTLY IF EXISTS "product_organization_id_idx";

DROP TABLE IF EXISTS "product";

DROP INDEX CONCURRENTLY IF EXISTS "collector_product_id_idx";

DROP TABLE IF EXISTS "collector";

DROP INDEX CONCURRENTLY IF EXISTS "exporter_product_id_idx";

DROP TABLE IF EXISTS "exporter";
