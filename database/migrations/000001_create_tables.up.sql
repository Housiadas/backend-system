-- Description: Create table users
CREATE TABLE users
(
    user_id       UUID        NOT NULL,
    name          TEXT        NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    roles         TEXT[]      NOT NULL,
    password_hash TEXT        NOT NULL,
    department    TEXT NULL,
    enabled       BOOLEAN     NOT NULL,
    date_created  TIMESTAMP   NOT NULL,
    date_updated  TIMESTAMP   NOT NULL,

    PRIMARY KEY (user_id)
);

-- Description: Create table products
CREATE TABLE products
(
    product_id   UUID           NOT NULL,
    user_id      UUID           NOT NULL,
    name         TEXT           NOT NULL,
    cost         NUMERIC(10, 2) NOT NULL,
    quantity     INT            NOT NULL,
    date_created TIMESTAMP      NOT NULL,
    date_updated TIMESTAMP      NOT NULL,

    PRIMARY KEY (product_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

-- Description: Create table audit
CREATE TABLE audit
(
    id         UUID      NOT NULL,
    obj_id     UUID      NOT NULL,
    obj_entity TEXT      NOT NULL,
    obj_name   TEXT      NOT NULL,
    actor_id   UUID      NOT NULL,
    action     TEXT      NOT NULL,
    data       JSONB NULL,
    message    TEXT NULL,
    timestamp  TIMESTAMP NOT NULL,

    PRIMARY KEY (id)
);