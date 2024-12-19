-- +goose Up
-- +goose StatementBegin
-- SET statement_timeout = 0;
-- SET lock_timeout = 0;
-- SET idle_in_transaction_session_timeout = 0;
-- SET client_encoding = 'UTF8';
-- SET standard_conforming_strings = on;
-- SET check_function_bodies = false;
-- SET xmloption = content;
-- SET client_min_messages = warning;
-- SET row_security = off;

-- --
-- -- Name: keys; Type: TABLE; Schema: public; Owner: zvault_user
-- --

-- CREATE TABLE keys (
--     id SERIAL PRIMARY KEY,
--     user_id TEXT,
--     client_id TEXT NOT NULL,
--     client_key TEXT NOT NULL,
--     mnemonic TEXT,
--     private_key TEXT NOT NULL UNIQUE,
--     created_at BIGINT NOT NULL
-- );

-- ALTER TABLE keys OWNER TO zvault_user;

-- +goose StatementEnd