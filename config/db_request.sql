-- Фильмы по имени актера
SELECT films.* FROM actors
    JOIN films_actors ON films_actors.actor_id = actors.id
    JOIN films ON films.id = films_actors.film_id
    WHERE actors.fullname = 'Anita Tsoy'
ORDER BY 1;

-- Актеры по названию фильма
SELECT actors.* FROM films
    JOIN films_actors ON films_actors.film_id = films.id
    JOIN actors ON actors.id = films_actors.actor_id
    WHERE films.title = 'Robocop 2'
ORDER BY 1;


-- Добавить фильм
INSERT INTO films (title, description, release, rating) VALUES
(
	'Robocop 4',
    'Dead police officer became android',
    1996,
    5
) RETURNING id;

-- Добавить актера
INSERT INTO actors (fullname, sex, dateofbirth) VALUES 
(
    'Alice Seleznyova',
    'female',
    '22-01-1986'
) RETURNING id;


-- Связать информацию о новом фильме с актерами этого фильма 
INSERT INTO films_actors (film_id, actor_id) VALUES
(4,2),
(4,4);

-- Фильмы по имени актера
SELECT films.* FROM actors
    JOIN films_actors ON films_actors.actor_id = actors.id
    JOIN films ON films.id = films_actors.film_id
    WHERE actors.fullname = 'Ivan Ivanov'
ORDER BY 1;

-- Актеры по названию фильма
SELECT actors.* FROM films
    JOIN films_actors ON films_actors.film_id = films.id
    JOIN actors ON actors.id = films_actors.actor_id
    WHERE films.title = 'Robocop 3'
ORDER BY 1;

-- Таблица сопоставления фильмов с актерами
SELECT * FROM films_actors;