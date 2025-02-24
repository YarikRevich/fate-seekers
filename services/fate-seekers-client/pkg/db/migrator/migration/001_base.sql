+goose Up
+goose StatementBegin
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: collections; Type: TABLE; Schema: public; 
--

CREATE TABLE collections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);

--
-- Name: flags; Type: TABLE; Schema: public; 
--

CREATE TABLE flags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL UNIQUE DEFAULT "",
    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

+goose StatementEnd