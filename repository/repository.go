package repository

import (
	"context"

	"github.com/filmoteka/entities"
)

type DBHandler interface {
	GetFilmInfo(ctx context.Context, title string) ([]entities.Film, error)
	GetActorInfo(ctx context.Context, fullname string) ([]entities.Actor, error)
	AddFilm(ctx context.Context, film entities.Film) (int, error)
	AddActor(ctx context.Context, actor entities.Actor) (int, error)
	DeleteFilm(ctx context.Context, film entities.Film) bool
	DeleteActor(ctx context.Context, actor entities.Actor) bool
	SetFilmInfo(ctx context.Context, film entities.Film) bool
	SetActorInfo(ctx context.Context, actor entities.Actor) bool
}
