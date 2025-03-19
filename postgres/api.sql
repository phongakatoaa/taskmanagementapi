\connect task_management;

CREATE SCHEMA IF NOT EXISTS api;
SET search_path TO api,public;

CREATE TYPE api.task_status AS ENUM ('PENDING', 'IN_PROGRESS', 'COMPLETED');

CREATE TABLE IF NOT EXISTS api.tasks
(
    id               SERIAL PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    description      TEXT,
    created_at       TIMESTAMPTZ     DEFAULT NOW(),
    due_date         TIMESTAMPTZ,
    status           api.task_status DEFAULT 'PENDING',
    assigned_user_id INT REFERENCES auth.users (id)
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON api.tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_user_id ON api.tasks (assigned_user_id);