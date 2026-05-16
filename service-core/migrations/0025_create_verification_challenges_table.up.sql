CREATE TYPE otp_type AS ENUM (
    'numeric',
    'magic_link'
);

CREATE TYPE otp_channel AS ENUM (
    'email',
    'sms'
);

CREATE TYPE otp_purpose AS ENUM (
    'register',
    'login',
    'password_reset',
    'email_verification'
);

CREATE TABLE verification_challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID,

    type otp_type NOT NULL,
    channel otp_channel NOT NULL,
    purpose otp_purpose NOT NULL,

    target TEXT NOT NULL,
    code_hash TEXT NOT NULL,

    attempt_count INT NOT NULL DEFAULT 0,

    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    consumed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_verification_challenges_user_id
    ON verification_challenges(user_id);

CREATE INDEX idx_verification_challenges_target
    ON verification_challenges(target);

CREATE INDEX idx_verification_challenges_purpose
    ON verification_challenges(purpose);

CREATE INDEX idx_verification_challenges_expires_at
    ON verification_challenges(expires_at);