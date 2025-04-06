-- USERS table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- WALLETS table
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    balance DECIMAL(18, 4) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL
);

-- TRANSACTIONS table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    type TEXT NOT NULL,
    amount DECIMAL(18, 4) NOT NULL,
    related_user_id UUID NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    key TEXT PRIMARY KEY,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    response TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Indexes
CREATE INDEX IF NOT EXISTS idx_tx_wallet_id ON transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_tx_related_user_id ON transactions(related_user_id);
CREATE INDEX IF NOT EXISTS idx_tx_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_tx_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_idempotency_composite ON idempotency_keys (key, method, path);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_balance ON wallets(balance);
-- Index on email for lookups (e.g., login or search)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);