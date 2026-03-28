-- +goose Up
SELECT 'up SQL query';

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (email, password_hash, role, created_at) VALUES ('admin', '$2a$12$4jbGkizBM7QONI22WTzYzOkM6BljUB.3NzV2oUnuzymgnCLvr.2HS', 'admin', NOW());

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;