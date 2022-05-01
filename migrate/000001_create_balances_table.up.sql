CREATE TABLE IF NOT EXISTS balances
(
    balance_id bigserial primary key,
    user_id    bigint not null UNIQUE ,
    balance    integer not null default 0
);