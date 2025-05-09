-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(30) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- packs table (with category)
CREATE TABLE packs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(30) UNIQUE NOT NULL,
    category TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- cards table
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question VARCHAR(255) NOT NULL,
    answer TEXT NOT NULL,
    pack_id UUID NOT NULL
        REFERENCES packs(id)
        ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- add flag to cards
ALTER TABLE cards
  ADD COLUMN IF NOT EXISTS last_wrong BOOLEAN DEFAULT FALSE;

-- subscriptions table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,
    pack_id UUID NOT NULL
        REFERENCES packs(id)
        ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- logs table
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,
    pack_id UUID NOT NULL
        REFERENCES packs(id)
        ON DELETE CASCADE,
    rating_improved INTEGER DEFAULT 0 CHECK (rating_improved >= 0),
    rating_worsen   INTEGER DEFAULT 0 CHECK (rating_worsen   >= 0),
    cards_learned   INTEGER DEFAULT 0 CHECK (cards_learned   >= 0),
    cards_mastered  INTEGER DEFAULT 0 CHECK (cards_mastered  >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- минимальная табличка статистики пользователя
CREATE TABLE IF NOT EXISTS user_stats (
    user_id UUID PRIMARY KEY
        REFERENCES users(id)
        ON DELETE CASCADE,
    rating        INT DEFAULT 0,
    packs_created INT DEFAULT 0,
    packs_mastered INT DEFAULT 0
);

-- Indexes
CREATE INDEX idx_users_username        ON users(username);
CREATE INDEX idx_packs_name            ON packs(name);
CREATE INDEX idx_cards_pack_id         ON cards(pack_id);
CREATE INDEX idx_cards_rating          ON cards(rating);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_pack_id ON subscriptions(pack_id);
CREATE INDEX idx_logs_user_id          ON logs(user_id);
CREATE INDEX idx_logs_pack_id          ON logs(pack_id);

-- Trigger function to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER set_updated_at_users
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_packs
BEFORE UPDATE ON packs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_cards
BEFORE UPDATE ON cards
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_subscriptions
BEFORE UPDATE ON subscriptions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_updated_at_logs
BEFORE UPDATE ON logs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();