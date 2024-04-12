-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS actors (
    id bigserial,
    fullname varchar(30),
    sex varchar (6),
    dateofbirth date
);

INSERT INTO actors (fullname, sex, dateofbirth) VALUES 
(
    'Anita Tsoy',
    'female',
    '07-07-1982'
),(
    'Ivan Ivanov',
    'male',
    '05-01-1980'
),(
    'Sergey Burunov',
    'male',
    '11-03-1976'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS actors;
-- +goose StatementEnd
