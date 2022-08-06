-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE transaction
(
    id           UUID                  NOT NULL,
    user_id      UUID                  NOT NULL,
    state        character varying(50) NOT NULL,
    amount       numeric(8, 2)         NOT NULL,
    processed_at timestamp             NOT NULL,
    is_canceled  bool,
    CONSTRAINT "PK_transaction" PRIMARY KEY (id)
);

CREATE TABLE users
(
    id      UUID          NOT NULL,
    balance numeric(8, 2) NOT NULL,
    CONSTRAINT "PK_user" PRIMARY KEY (id)
);

INSERT INTO users(id, balance)
VALUES ('0268c107-c4cf-45a5-8547-c16f86616d61', 100.00);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS transaction;
DROP TABLE IF EXISTS users;