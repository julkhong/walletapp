-- SQL Migration to create tables
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT,
    email TEXT,
    created_at TIMESTAMP
);

CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    balance BIGINT,
    created_at TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    wallet_id UUID REFERENCES wallets(id),
    type TEXT,
    amount BIGINT,
    related_user_id UUID,
    created_at TIMESTAMP
);
