CREATE TABLE IF NOT EXISTS tokens (
    jti TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users (
user_id PRIMARY KEY,
username TEXT UNIQUE,
password_hash TEXT,S
created_at TIMESTAMP NOT NULL,
last_login TIMESTAMP NOT NULL,

FOREIGN KEY(user_id) REFERENCES tokens(user_id)
);