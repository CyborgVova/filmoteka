-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS films (
    id bigserial,
    title varchar(30),
    description varchar(150),
    release int,
    rating smallint check (rating >= 0 and rating <= 10)
);

CREATE INDEX ON films USING BTREE (lower(title));

INSERT INTO films (title, description, release, rating) VALUES
(
    'Robocop',
    'Dead police officer became android',
    1990,
    7
),(
    'Robocop 2',
    'Dead police officer became android',
    1992,
    8
),(
	'Robocop 3',
    'Dead police officer became android',
    1994,
    6
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS films;
-- +goose StatementEnd
