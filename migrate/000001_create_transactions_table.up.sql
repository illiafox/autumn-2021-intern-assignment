CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id bigserial primary key,
    balance_id     bigint    not null,
    from_id        bigint,
    action         integer   not null,
    date           timestamp not null,
    description    text
);