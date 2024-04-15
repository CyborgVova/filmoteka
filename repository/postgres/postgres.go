package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/filmoteka/entities"
)

type Film struct {
	ID          int
	Title       string
	Description string
	Release     int
	Rating      int
}

type Actor struct {
	ID         int
	FullName   string
	Sex        string
	DayOfBirth time.Time
}

type FilmsActors struct {
	FilmID  int
	ActorID int
}

type Repository struct {
	Conn *pgx.Conn
}

const (
	dbstring = "postgres://docker:docker@localhost:5432/docker"
)

func NewRepository() *Repository {
	conn, err := pgx.Connect(context.Background(), dbstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	return &Repository{Conn: conn}
}

func (r *Repository) GetFilmInfo(ctx context.Context, title string) ([]entities.Film, error) {
	rows, err := r.Conn.Query(ctx,
		fmt.Sprintf("SELECT * FROM films "+
			"WHERE lower(title) SIMILAR TO '%%%s%%'",
			strings.ToLower(title)))
	if err != nil {
		log.Fatalf("error gettitng films info: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	films := []entities.Film{}
	for rows.Next() {
		film := entities.Film{}
		rows.Scan(&film.ID, &film.Title, &film.Description, &film.Release, &film.Rating)

		films = append(films, film)
	}

	for i, film := range films {
		films[i].Actors = r.getAllFilmActors(ctx, film.ID)
	}

	return films, nil
}

func (r *Repository) getAllFilmActors(ctx context.Context, filmid int) []entities.Actor {
	rows, err := r.Conn.Query(ctx, "SELECT actors.* FROM actors "+
		"JOIN films_actors ON films_actors.actor_id = actors.id "+
		"JOIN films ON films.id = films_actors.film_id "+
		"WHERE films.id = $1", filmid)
	if err != nil {
		log.Fatalf("error getting all film actors: %v\n", err)
		return nil
	}
	defer rows.Close()

	actors := []entities.Actor{}
	for rows.Next() {
		actor := entities.Actor{}
		rows.Scan(&actor.ID, &actor.FullName, &actor.Sex, &actor.DateOfBirth)
		actors = append(actors, actor)
	}
	return actors
}

func (r *Repository) GetActorInfo(ctx context.Context, fullname string) ([]entities.Actor, error) {
	rows, err := r.Conn.Query(ctx,
		fmt.Sprintf("SELECT * FROM actors "+
			"WHERE lower(fullname) SIMILAR TO '%%%s%%'",
			strings.ToLower(fullname)))
	if err != nil {
		log.Fatal("error getting actors:", err)
	}
	defer rows.Close()

	actors := []entities.Actor{}
	for rows.Next() {
		actor := entities.Actor{}
		rows.Scan(&actor.ID, &actor.FullName, &actor.Sex, &actor.DateOfBirth)
		actors = append(actors, actor)
	}

	for i, actor := range actors {
		actors[i].Films = r.getAllFilmsActor(ctx, actor.ID)
	}
	return actors, nil
}

func (r *Repository) getAllFilmsActor(ctx context.Context, actorid int) []entities.Film {
	rows, err := r.Conn.Query(ctx, "SELECT films.* from films "+
		"JOIN films_actors ON films_actors.film_id = films.id "+
		"JOIN actors ON actors.id = films_actors.actor_id "+
		"WHERE actors.id = $1", actorid)
	if err != nil {
		log.Fatal("error getting all films actor:", err)
	}

	films := []entities.Film{}
	for rows.Next() {
		film := entities.Film{}
		rows.Scan(&film.ID, &film.Title, &film.Description, &film.Release, &film.Rating)
		films = append(films, film)
	}
	return films
}

func (r *Repository) AddFilm(ctx context.Context, film entities.Film) (int, error) {
	row := r.Conn.QueryRow(ctx, fmt.Sprintf("INSERT INTO films (title, description, release, rating) "+
		"VALUES ('%s', '%s', '%d', '%d') RETURNING ID",
		film.Title, film.Description, film.Release, film.Rating))
	var id int
	row.Scan(&id)
	if id == 0 {
		return 0, errors.New("error insert a film")
	}
	return id, nil
}

func (r *Repository) AddActor(ctx context.Context, actor entities.Actor) (int, error) {
	row := r.Conn.QueryRow(ctx, fmt.Sprintf("INSERT INTO actors (fullname, sex, dateofbirth) "+
		"VALUES ('%s', '%s', '%s') RETURNING id",
		actor.FullName, actor.Sex, actor.DateOfBirth))
	var id int
	row.Scan(&id)
	if id == 0 {
		return 0, errors.New("error insert an actor")
	}
	return id, nil
}

func (r *Repository) SetFilmInfo(ctx context.Context, film entities.Film) bool {
	return true
}

func (r *Repository) SetActorInfo(ctx context.Context, actor entities.Actor) bool {
	return false
}
