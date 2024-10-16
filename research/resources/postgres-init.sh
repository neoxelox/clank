#!/bin/bash

psql -U postgres <<-END

CREATE DATABASE clank;

END

psql -U postgres -d clank <<-END

CREATE EXTENSION vector;

CREATE TABLE issue (
    id VARCHAR(20) PRIMARY KEY,
    text TEXT,
    embedding VECTOR(1536)
);

CREATE INDEX CONCURRENTLY issue_embedding_idx ON issue USING hnsw (embedding vector_cosine_ops)
    WITH (m = 24, ef_construction = 100);

CREATE TABLE suggestion (
    id VARCHAR(20) PRIMARY KEY,
    text TEXT,
    embedding VECTOR(1536)
);

CREATE INDEX CONCURRENTLY suggestion_embedding_idx ON suggestion USING hnsw (embedding vector_cosine_ops)
    WITH (m = 24, ef_construction = 100);

END
