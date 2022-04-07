BEGIN;

CREATE DATABASE IF NOT EXISTS avito;

CREATE TABLE IF NOT EXISTS avito.balances
(
    balance_id bigint primary key auto_increment,
    user_id    bigint not null,
    balance    integer not null default 0
);

CREATE TABLE IF NOT EXISTS avito.transactions
(
    transaction_id bigint primary key auto_increment,
    balance_id     bigint   not null,
    from_id        bigint,
    action         integer   not null,
    date           timestamp not null,
    description    tinytext,
    FOREIGN KEY (balance_id) REFERENCES balances (balance_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    FOREIGN KEY (from_id) REFERENCES balances (balance_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

COMMIT;