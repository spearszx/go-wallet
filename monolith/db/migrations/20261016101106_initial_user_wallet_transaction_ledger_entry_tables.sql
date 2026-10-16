-- +goose Up

CREATE TABLE users (
    id              VARCHAR(36) PRIMARY KEY,
    full_name       VARCHAR(100) NOT NULL,
    email           VARCHAR(150) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP NULL
);
CREATE INDEX idx_users_email       ON users(email);
CREATE INDEX idx_users_is_active   ON users(is_active);
CREATE INDEX idx_users_deleted_at  ON users(deleted_at);

CREATE TABLE wallets (
    id          VARCHAR(36) PRIMARY KEY,
    user_id     VARCHAR(36) NOT NULL,
    balance     DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency    CHAR(3) NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    version     INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    CONSTRAINT fk_wallets_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_wallets_status
        CHECK (status IN ('active', 'frozen'))
);
CREATE INDEX idx_wallets_user_id       ON wallets(user_id);
CREATE INDEX idx_wallets_status        ON wallets(status);
CREATE INDEX idx_wallets_user_currency ON wallets(user_id, currency);
CREATE INDEX idx_wallets_deleted_at    ON wallets(deleted_at);

CREATE TABLE transactions (
    id                  VARCHAR(36) PRIMARY KEY,
    sender_wallet_id    VARCHAR(36) NULL,
    receiver_wallet_id  VARCHAR(36) NULL,
    amount              DECIMAL(15, 2) NOT NULL,
    description         TEXT NULL,
    idempotency_key     VARCHAR(100) NOT NULL UNIQUE,
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP NULL,
    CONSTRAINT fk_tx_sender
        FOREIGN KEY (sender_wallet_id) REFERENCES wallets(id)
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_tx_receiver
        FOREIGN KEY (receiver_wallet_id) REFERENCES wallets(id)
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT chk_tx_status
        CHECK (status IN ('PENDING', 'SUCCESS', 'FAILED')),
    CONSTRAINT chk_tx_amount
        CHECK (amount > 0)
);
CREATE INDEX idx_tx_sender_wallet     ON transactions(sender_wallet_id);
CREATE INDEX idx_tx_receiver_wallet   ON transactions(receiver_wallet_id);
CREATE INDEX idx_tx_status            ON transactions(status);
CREATE INDEX idx_tx_idempotency_key   ON transactions(idempotency_key);
CREATE INDEX idx_tx_created_at        ON transactions(created_at);
CREATE INDEX idx_tx_deleted_at        ON transactions(deleted_at);

CREATE TABLE ledger_entries (
    id              VARCHAR(36) PRIMARY KEY,
    wallet_id       VARCHAR(36) NOT NULL,
    transaction_id  VARCHAR(36) NOT NULL,
    entry_type      VARCHAR(10) NOT NULL,
    amount          DECIMAL(15, 2) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ledger_wallet
        FOREIGN KEY (wallet_id) REFERENCES wallets(id)
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_ledger_transaction
        FOREIGN KEY (transaction_id) REFERENCES transactions(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_ledger_entry_type
        CHECK (entry_type IN ('CREDIT', 'DEBIT')),
    CONSTRAINT chk_ledger_amount
        CHECK (amount > 0)
);
CREATE INDEX idx_ledger_wallet_id      ON ledger_entries(wallet_id);
CREATE INDEX idx_ledger_transaction_id ON ledger_entries(transaction_id);
CREATE INDEX idx_ledger_entry_type     ON ledger_entries(entry_type);
CREATE INDEX idx_ledger_created_at     ON ledger_entries(created_at);

-- +goose Down
DROP TABLE IF EXISTS ledger_entries CASCADE;
DROP TABLE IF EXISTS transactions   CASCADE;
DROP TABLE IF EXISTS wallets        CASCADE;
DROP TABLE IF EXISTS users          CASCADE;