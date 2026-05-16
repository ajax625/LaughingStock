CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    telegram_chat_id TEXT,
    alpaca_key TEXT,
    alpaca_secret TEXT,
    webhook_token TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_symbols (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    symbol TEXT NOT NULL,
    trader_type TEXT NOT NULL DEFAULT 'day',
    min_roi DOUBLE PRECISION NOT NULL DEFAULT 0.5,
    min_rr DOUBLE PRECISION NOT NULL DEFAULT 2.0,
    min_confidence DOUBLE PRECISION NOT NULL DEFAULT 0.6,
    tradex_subscription_id TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, symbol)
);
