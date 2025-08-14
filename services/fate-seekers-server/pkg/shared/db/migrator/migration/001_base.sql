-- +goose Up
-- +goose StatementBegin

--
-- Name: sessions; Type: TABLE; Schema: public; 
--

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    seed INTEGER NOT NULL,
    issuer INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (issuer) REFERENCES users(id)
);

--
-- Name: lobbies; Type: TABLE; Schema: public; 
--

CREATE TABLE lobbies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    session_id INTEGER NOT NULL,
    host BOOLEAN NOT NULL,
    skin INTEGER NOT NULL,
    health INTEGER NOT NULL,
    eliminated BOOLEAN NOT NULL,
    position_x REAL NOT NULL,
    position_y REAL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
    FOREIGN KEY (session_id) REFERENCES sessions(id)
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