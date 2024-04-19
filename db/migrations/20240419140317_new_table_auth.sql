-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS auth (
    login varchar NOT NULL,
    password varchar NOT NULL
);

INSERT INTO auth (login, password) VALUES (
    'admin',
    '$2a$14$pUrzBdh29kq4qehX4gu7oe1GQUP/xvRCOimbkv/4Qtc1vS2gAl9RW'
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS auth;