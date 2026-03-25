CREATE TABLE IF NOT EXISTS users (
    id          VARCHAR(100) NOT NULL PRIMARY KEY,
    password    VARCHAR(255) NOT NULL,
    username    VARCHAR(100) NOT NULL,
    email       VARCHAR(100),
    token       VARCHAR(500),
    created_at  BIGINT       NOT NULL,
    updated_at  BIGINT       NOT NULL
);
