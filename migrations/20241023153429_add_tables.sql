-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE transaction_type AS ENUM ('DEBIT', 'CREDIT');
CREATE TYPE payment_status AS ENUM ('DRAFT', 'PENDING', 'COMPLETED', 'FAILED');
CREATE TYPE account_type AS ENUM ('CLIENT', 'BANK');
CREATE TYPE currency AS ENUM (
    'RUB' -- Российский рубль
);
CREATE TYPE send_transaction_status AS ENUM ('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED');

-- Таблица счетов
CREATE TABLE IF NOT EXISTS accounts
(
    id           UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    user_id      UUID,
    balance      BIGINT                   NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency     currency                 NOT NULL DEFAULT 'RUB',
    account_type account_type             NOT NULL DEFAULT 'CLIENT',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at   TIMESTAMP WITH TIME ZONE,
    CONSTRAINT user_id_not_null_for_client_accounts
        CHECK (NOT (account_type = 'CLIENT' AND user_id IS NULL))
);

CREATE UNIQUE INDEX unique_bank_account ON accounts (account_type) WHERE account_type = 'BANK';

INSERT INTO accounts (user_id, balance, currency, account_type, created_at)
VALUES ( NULL, 0, 'RUB', 'BANK', now());

-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments
(
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_account_id   UUID                     NOT NULL,
    receiver_account_id UUID                     NOT NULL,
    amount              BIGINT                   NOT NULL CHECK (amount > 0),
    currency            currency                 NOT NULL,
    description         TEXT,
    status              payment_status           NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (sender_account_id) REFERENCES accounts (id),
    FOREIGN KEY (receiver_account_id) REFERENCES accounts (id)
);

-- Таблица транзакций
CREATE TABLE IF NOT EXISTS transactions
(
    id          UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    payment_id  UUID                     NOT NULL,
    account_id  UUID                     NOT NULL,
    type        transaction_type         NOT NULL CHECK (type IN ('DEBIT', 'CREDIT')),
    amount      BIGINT                   NOT NULL CHECK (amount > 0),
    currency    currency                 NOT NULL,
    description TEXT,
    status      payment_status           NOT NULL,
    send_status send_transaction_status  NOT NULL DEFAULT 'PENDING',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (payment_id) REFERENCES payments (id),
    FOREIGN KEY (account_id) REFERENCES accounts (id)
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS payment_status;