CREATE TABLE IF NOT EXISTS refresh_tokens(
    id BIGSERIAL PRIMARY KEY,
    user_code VARCHAR(20) NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_code)
        REFERENCES users(user_code)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_code
ON refresh_tokens(user_code);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token
ON refresh_tokens(token);