-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
    username varchar(32) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(255) NOT NULL,
    verified_at timestamp with time zone,
    verification_token varchar(255) NOT NULL,
    password_reset_token varchar(255),
    password_reset_expires_at timestamp with time zone,
    is_admin boolean NOT NULL,
    lang varchar(255) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE UNIQUE INDEX idx_users_username ON users USING btree (username);
CREATE UNIQUE INDEX idx_users_email ON users USING btree (email);
CREATE UNIQUE INDEX idx_users_verification_token ON users USING btree (verification_token);
CREATE UNIQUE INDEX idx_users_password_reset_token ON users USING btree (password_reset_token);
CREATE INDEX idx_users_deleted_at ON users USING btree (deleted_at);

CREATE TABLE sessions (
    id UUID NOT NULL PRIMARY KEY,
    refresh_token varchar(255) NOT NULL,
    user_id UUID NOT NULL,
    type varchar(255) NOT NULL,
    ip varchar(255) NOT NULL,
    user_agent varchar(255) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    last_used_at timestamp with time zone NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    CONSTRAINT fk_sessions_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_sessions_refresh_token ON sessions USING btree (refresh_token);
CREATE INDEX idx_sessions_user_id ON sessions USING btree (user_id);
CREATE INDEX idx_sessions_expires_at ON sessions USING btree (expires_at);

-- +goose Down
DROP TABLE users;
DROP TABLE sessions;
