CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE provider AS ENUM (
    'PHONE',
    'GOOGLE'
);

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(20),
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(20) UNIQUE,
    provider provider NOT NULL,
    profile_picture VARCHAR(500), -- In case of google sign in frontend may provider image
    meta JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);