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
    started BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (issuer) REFERENCES users(id)
);

--
-- Name: generations; Type: TABLE; Schema: public;
--

CREATE TABLE generations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    instance TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    active BOOLEAN NOT NULL,
    position_x REAL NOT NULL,
    position_y REAL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

--
-- Name: associations; Type: TABLE; Schema: public;
--

CREATE TABLE associations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    generation_id INTEGER NOT NULL,
    instance TEXT NOT NULL,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
    FOREIGN KEY (generation_id) REFERENCES generations(id) ON DELETE CASCADE
);

--
-- Name: lobbies; Type: TABLE; Schema: public; 
--

CREATE TABLE lobbies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_id INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    host BOOLEAN NOT NULL,
    skin INTEGER NOT NULL,
    health INTEGER NOT NULL DEFAULT 100,
    eliminated BOOLEAN NOT NULL,
    position_x REAL NOT NULL,
    position_y REAL NOT NULL,
    position_static BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

--
-- Name: idx_lobbies_user_id_session_id; Type: INDEX; Schema: public; 
--

CREATE UNIQUE INDEX idx_lobbies_user_id_session_id
ON lobbies (user_id, session_id);

--
-- Name: idx_lobbies_user_id_session_id_skin; Type: INDEX; Schema: public; 
--

CREATE UNIQUE INDEX idx_lobbies_user_id_session_id_skin
ON lobbies (user_id, session_id, skin);

--
-- Name: inventory; Type: TABLE; Schema: public; 
--

CREATE TABLE inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lobby_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    session_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
    FOREIGN KEY (lobby_id) REFERENCES lobbies(id) ON DELETE CASCADE
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