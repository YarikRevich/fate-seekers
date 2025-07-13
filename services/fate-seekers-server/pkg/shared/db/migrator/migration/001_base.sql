-- +goose Up
-- +goose StatementBegin

--
-- Name: sessions; Type: TABLE; Schema: public; 
--

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    issuer INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (issuer) REFERENCES users(id)
);

--
-- Name: messages; Type: TABLE; Schema: public; 
--

CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    issuer INTEGER NOT NULL,
    FOREIGN KEY (issuer) REFERENCES users(id)
);

--
-- Name: users; Type: TABLE; Schema: public; 
--

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);

-- +goose StatementEnd