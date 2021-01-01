CREATE DATABASE chat;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(156) NOT NULL UNIQUE,
    hash TEXT NOT NULL
);