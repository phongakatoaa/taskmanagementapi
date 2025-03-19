\connect task_management;

CREATE SCHEMA IF NOT EXISTS auth;
SET search_path TO auth,public;

CREATE TYPE auth.role AS ENUM ('EMPLOYER', 'EMPLOYEE');

CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    role       auth.role    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);
