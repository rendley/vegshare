BEGIN;

CREATE TABLE users (
                       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email         VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
                                id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token_hash VARCHAR(255) UNIQUE NOT NULL,
                                expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

COMMIT;