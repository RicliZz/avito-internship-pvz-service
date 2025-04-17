CREATE TYPE role AS ENUM ('moderator', 'employee');

CREATE EXTENSION pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    "ID" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text NOT NULL,
    password text NOT NULL,
    role role NOT NULL,
    constraint unique_email UNIQUE (email)
)