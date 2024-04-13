package repository

import (
	"context"

	"github.com/filmoteka/entities"
)

type DBHandler interface {
	GetFilmInfo(ctx context.Context, title string) ([]entities.Film, error)
	GetActorInfo(ctx context.Context, fullname string) ([]entities.Actor, error)
	AddFilm(ctx context.Context, film entities.Film) error
	SetFilmInfo(ctx context.Context, film entities.Film) bool
	SetActorInfo(ctx context.Context, actor entities.Actor) bool
}
