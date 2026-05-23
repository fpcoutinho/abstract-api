-- initial schema
CREATE TABLE IF NOT EXISTS users (
    uid TEXT PRIMARY KEY,
    name TEXT,
    gender TEXT,
    shell_balance BIGINT DEFAULT 0,
    avatar_idx INT,
    active_frame TEXT,
    active_accessory TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trails (
    id SERIAL PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    title TEXT,
    description TEXT,
    order_index INT,
    published BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    trail_id INT REFERENCES trails(id),
    slug TEXT,
    title TEXT,
    content_md TEXT,
    reward_shells INT DEFAULT 0,
    order_index INT,
    published BOOLEAN DEFAULT false,
    answer_key_hash TEXT
);

CREATE TABLE IF NOT EXISTS user_submissions (
    id SERIAL PRIMARY KEY,
    uid TEXT REFERENCES users(uid),
    mission_id INT REFERENCES missions(id),
    answers_json JSONB,
    is_correct BOOLEAN,
    earned_shells INT DEFAULT 0,
    idempotency_key TEXT,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS shell_ledger (
    id SERIAL PRIMARY KEY,
    uid TEXT REFERENCES users(uid),
    submission_id INT REFERENCES user_submissions(id),
    delta BIGINT,
    reason TEXT,
    balance_before BIGINT,
    balance_after BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_unlocks (
    id SERIAL PRIMARY KEY,
    uid TEXT REFERENCES users(uid),
    item_type TEXT,
    item_code TEXT,
    unlocked_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);