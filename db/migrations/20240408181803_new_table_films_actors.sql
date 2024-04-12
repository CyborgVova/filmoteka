-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS films_actors (
    film_id int,
    actor_id int
);

INSERT INTO films_actors (film_id, actor_id) VALUES
( 1, 1 ),
( 1, 2 ),
( 2, 1 ),
( 2, 2 ),
( 3, 1 ),
( 3, 2 ),
( 3, 3 ),;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS films_actors;
-- +goose StatementEnd
