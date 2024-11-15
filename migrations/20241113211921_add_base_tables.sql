-- +goose Up
create table TRANSACTION
(
    ID   serial primary key,
    DATE date not null,
    AMOUNT decimal not null,
    FILE_PATH varchar(255) not null
);

create table ACCOUNT_SUMMARY
(
    ID   serial primary key,
    EMAIL varchar(255) not null,
    TOTAL_BALANCE decimal not null,
    AVERAGE_DEBIT_AMOUNT decimal not null,
    AVERAGE_CREDIT_AMOUNT decimal not null,
    TRANSACTIONS_PER_MONTH jsonb not null,
    FILE_PATH varchar(255) not null
);

-- +goose Down
DROP TABLE IF EXISTS ACCOUNT_SUMMARY;
DROP TABLE IF EXISTS TRANSACTION;
