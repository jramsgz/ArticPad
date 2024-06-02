-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    verified_at DATETIME,
    verification_token VARCHAR(255) NOT NULL,
    password_reset_token VARCHAR(255),
    password_reset_expires_at DATETIME,
    is_admin BOOLEAN NOT NULL,
    lang VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME
);

CREATE UNIQUE INDEX idx_users_username ON users (username);
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_verification_token ON users (verification_token);
CREATE UNIQUE INDEX idx_users_password_reset_token ON users (password_reset_token);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE TABLE sessions (
    id UUID NOT NULL PRIMARY KEY,
    refresh_token VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    type VARCHAR(255) NOT NULL,
    ip VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    last_used_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    CONSTRAINT fk_sessions_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- +goose Down
DROP TABLE users;
DROP TABLE sessions;
