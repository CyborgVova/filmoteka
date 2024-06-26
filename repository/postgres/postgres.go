package postgres

import (
	"context"
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
	localhost = "postgres://docker:docker@localhost:5432/docker"
)

func NewRepository() *Repository {
	dbstring := os.Getenv("GOOSE_DBSTRING")
	if dbstring == "" {
		dbstring = localhost
	}
	conn, err := pgx.Connect(context.Background(), dbstring)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &Repository{Conn: conn}
}

func (r *Repository) GetFilmInfo(ctx context.Context, title, order string) ([]entities.Film, error) {
	ascDesc := "ASC"
	if order == "rating" {
		ascDesc = "DESC"
	}
	rows, err := r.Conn.Query(ctx,
		fmt.Sprintf("SELECT * FROM films "+
			"WHERE lower(title) SIMILAR TO '%%%s%%' ORDER BY %s %s",
			strings.ToLower(title), order, ascDesc))
	if err != nil {
		log.Printf("error getting films info: %v\n", err)
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

func (r *Repository) GetActorInfo(ctx context.Context, fullname string) ([]entities.Actor, error) {
	rows, err := r.Conn.Query(ctx,
		fmt.Sprintf("SELECT * FROM actors "+
			"WHERE lower(fullname) SIMILAR TO '%%%s%%'",
			strings.ToLower(fullname)))
	if err != nil {
		log.Printf("error getting actors info: %s", err)
		return nil, err
	}
	defer rows.Close()

	actors := []entities.Actor{}
	for rows.Next() {
		actor := entities.Actor{}
		t := time.Time{}
		rows.Scan(&actor.ID, &actor.FullName, &actor.Sex, &t)
		actor.DateOfBirth = t.Format("02/01/2006")
		actors = append(actors, actor)
	}

	for i, actor := range actors {
		actors[i].Films = r.getAllFilmsActor(ctx, actor.ID)
	}
	return actors, nil
}

func (r *Repository) AddFilm(ctx context.Context, film entities.Film) (int, error) {
	id := r.findFilm(ctx, film)
	if id == 0 {
		err := r.Conn.QueryRow(ctx, fmt.Sprintf("INSERT INTO films (title, description, release, rating) "+
			"VALUES ('%s', '%s', '%d', '%d') RETURNING ID",
			film.Title, film.Description, film.Release, film.Rating)).Scan(&id)
		if err != nil {
			return id, err
		}
	} else {
		film.ID = id
	}

	r.fillActorsByFilm(ctx, film)
	return id, nil
}

func (r *Repository) AddActor(ctx context.Context, actor entities.Actor) (int, error) {
	id := r.findActor(ctx, actor)
	if id == 0 {
		err := r.Conn.QueryRow(ctx, fmt.Sprintf("INSERT INTO actors (fullname, sex, dateofbirth) "+
			"VALUES ('%s', '%s', '%s') RETURNING id",
			actor.FullName, actor.Sex, actor.DateOfBirth)).Scan(&id)
		if err != nil {
			return id, err
		}
	} else {
		actor.ID = id
	}

	r.fillFilmsByActor(ctx, actor)
	return id, nil
}

var filmFields = []string{"title", "description", "release", "rating"}

func (r *Repository) SetFilmInfo(ctx context.Context, m map[string]interface{}) bool {
	if !r.hasFilm(m["id"].(string)) {
		log.Printf("film with id = %s does not exist\n", m["id"])
		return false
	}

	for _, field := range filmFields {
		if value, ok := m[field]; ok {
			_, err := r.Conn.Exec(ctx, fmt.Sprintf("UPDATE films SET %s = '%v' WHERE id = %s", field, value, m["id"]))
			if err != nil {
				log.Println("error to update film info:", err)
				return false
			}
		}
	}
	return true
}

var actorFields = []string{"fullname", "sex", "dateofbirth"}

func (r *Repository) SetActorInfo(ctx context.Context, m map[string]interface{}) bool {
	if !r.hasActor(m["id"].(string)) {
		log.Printf("actor with id = %s does not exist\n", m["id"])
		return false
	}

	for _, field := range actorFields {
		if value, ok := m[field]; ok {
			_, err := r.Conn.Exec(ctx, fmt.Sprintf("UPDATE actors SET %s = '%v' WHERE actors.id = %s", field, value, m["id"]))
			if err != nil {
				log.Println("error to update actor info:", err)
				return false
			}
		}
	}
	return true
}

func (r *Repository) DeleteActor(ctx context.Context, actor entities.Actor) bool {
	var id = 0
	r.Conn.QueryRow(ctx, "SELECT id FROM actors WHERE fullname = $1 "+
		"AND dateofbirth = $2", actor.FullName, actor.DateOfBirth).Scan(&id)
	if id == 0 {
		log.Println("such actor not exist")
		return false
	}
	var row interface{}
	r.Conn.QueryRow(ctx, "DELETE FROM films_actors WHERE actor_id = $1", id).Scan(&row)
	r.Conn.QueryRow(ctx, "DELETE FROM actors WHERE id = $1", id).Scan(&row)
	return true
}

func (r *Repository) DeleteFilm(ctx context.Context, film entities.Film) bool {
	var id = 0
	r.Conn.QueryRow(ctx, "SELECT id FROM films WHERE title = $1 "+
		"AND release = $2", film.Title, film.Release).Scan(&id)
	if id == 0 {
		log.Println("such film not exist")
		return false
	}

	var row interface{}
	r.Conn.QueryRow(ctx, "DELETE FROM films_actors WHERE film_id = $1", id).Scan(&row)
	r.Conn.QueryRow(ctx, "DELETE FROM films WHERE id = $1", id).Scan(&row)
	return true
}

func (r *Repository) hasFilm(id string) bool {
	exist := 0
	r.Conn.QueryRow(context.Background(), "SELECT count(*) FROM films WHERE id = $1", id).Scan(&exist)
	return exist != 0
}

func (r *Repository) hasActor(id string) bool {
	exist := 0
	r.Conn.QueryRow(context.Background(), "SELECT count(*) FROM actors WHERE id = $1", id).Scan(&exist)
	return exist != 0
}

func (r *Repository) findFilm(ctx context.Context, film entities.Film) (id int) {
	r.Conn.QueryRow(ctx, "SELECT id FROM films "+
		"WHERE films.title = $1 AND films.release = $2",
		film.Title, film.Release).Scan(&id)
	return
}

func (r *Repository) findActor(ctx context.Context, actor entities.Actor) (id int) {
	r.Conn.QueryRow(ctx, "SELECT id FROM actors "+
		"WHERE actors.fullname = $1 AND actors.dateofbirth = $2",
		actor.FullName, actor.DateOfBirth).Scan(&id)
	return
}

func (r *Repository) fillFilmsByActor(ctx context.Context, actor entities.Actor) {
	for _, film := range actor.Films {
		id, _ := r.AddFilm(ctx, film)
		r.addFilmsActors(ctx, id, actor.ID)
	}
}

func (r *Repository) fillActorsByFilm(ctx context.Context, film entities.Film) {
	for _, actor := range film.Actors {
		id, _ := r.AddActor(ctx, actor)
		r.addFilmsActors(ctx, film.ID, id)
	}
}

func (r *Repository) addFilmsActors(ctx context.Context, filmID, actorID int) {
	var count int
	r.Conn.QueryRow(ctx, "SELECT count(*) FROM films_actors "+
		"WHERE film_id = $1 AND actor_id = $2",
		filmID, actorID).Scan(&count)
	if count == 0 {
		r.Conn.QueryRow(ctx, fmt.Sprintf("INSERT INTO films_actors(film_id, actor_id) "+
			"VALUES ('%d', '%d')", filmID, actorID))
	}
}

func (r *Repository) getAllFilmActors(ctx context.Context, filmid int) []entities.Actor {
	rows, _ := r.Conn.Query(ctx, "SELECT actors.* FROM actors "+
		"JOIN films_actors ON films_actors.actor_id = actors.id "+
		"JOIN films ON films.id = films_actors.film_id "+
		"WHERE films.id = $1", filmid)
	defer rows.Close()

	actors := []entities.Actor{}
	for rows.Next() {
		actor := entities.Actor{}
		t := time.Time{}
		rows.Scan(&actor.ID, &actor.FullName, &actor.Sex, &t)
		actor.DateOfBirth = t.Format("02/01/2006")
		actors = append(actors, actor)
	}
	return actors
}

func (r *Repository) getAllFilmsActor(ctx context.Context, actorid int) []entities.Film {
	rows, _ := r.Conn.Query(ctx, "SELECT films.* from films "+
		"JOIN films_actors ON films_actors.film_id = films.id "+
		"JOIN actors ON actors.id = films_actors.actor_id "+
		"WHERE actors.id = $1", actorid)
	defer rows.Close()

	films := []entities.Film{}
	for rows.Next() {
		film := entities.Film{}
		rows.Scan(&film.ID, &film.Title, &film.Description, &film.Release, &film.Rating)
		films = append(films, film)
	}
	return films
}
