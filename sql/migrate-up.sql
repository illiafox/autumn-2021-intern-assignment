BEGIN;

CREATE TABLE IF NOT EXISTS balances
(
    balance_id bigserial primary key,
    user_id    bigint not null,
    balance    integer not null default 0
);

CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id bigserial primary key,
    balance_id     bigint   not null,
    from_id        bigint,
    action         integer   not null,
    date           timestamp not null,
    description    text
);

COMMIT;