-- +goose Up
-- +goose StatementBegin

--
-- Name: collections; Type: TABLE; Schema: public; 
--

CREATE TABLE collections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    path TEXT NOT NULL UNIQUE,
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

-- +goose StatementEnd