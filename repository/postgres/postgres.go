package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
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

type User struct {
	id   int
	name string
	role string
}

func (r *Repository) GetFilmInfo(ctx context.Context, title string) ([]entities.Film, error) {
	rows, err := r.Conn.Query(ctx, fmt.Sprintf("select * from films where title similar to '%%%s%%'", title))
	if err != nil {
		log.Fatalf("error_1: %v\n", err)
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
		films[i].Actors = r.GetAllFilmActors(ctx, film.ID)
	}

	return films, nil
}

func (r *Repository) GetAllFilmActors(ctx context.Context, filmid int) []entities.Actor {
	rows, err := r.Conn.Query(ctx, "select actors.* from actors "+
		"JOIN films_actors ON films_actors.actor_id = actors.id "+
		"JOIN films ON films.id = films_actors.film_id "+
		"WHERE films.id = $1", filmid)
	if err != nil {
		log.Fatalf("error_2: %v\n", err)
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

func (r *Repository) GetActorInfo(ctx context.Context, fullname string) (entities.Actor, error) {
	actor := entities.Actor{}
	return actor, nil
}

func (r *Repository) AddFilm(ctx context.Context, film entities.Film) error {
	return nil
}

func (r *Repository) SetFilmInfo(ctx context.Context, film entities.Film) bool {
	return true
}

func (r *Repository) SetActorInfo(ctx context.Context, actor entities.Actor) bool {
	return false
}
