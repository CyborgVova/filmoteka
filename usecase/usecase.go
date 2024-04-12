package usecase

import (
	"context"

	"github.com/filmoteka/entities"
	"github.com/filmoteka/repository"
	"github.com/filmoteka/repository/postgres"
)

type UseCase struct {
	Repo repository.DBHandler
}

func (u *UseCase) GetFilmInfo(ctx context.Context, title string) ([]entities.Film, error) {
	return u.Repo.(*postgres.Repository).GetFilmInfo(ctx, title)
}

func (u *UseCase) GetActorInfo(ctx context.Context, fullname string) (entities.Actor, error) {
	return u.Repo.GetActorInfo(ctx, fullname)
}

func (u *UseCase) AddFilm(ctx context.Context, film entities.Film) error {
	return u.Repo.AddFilm(ctx, film)
}

func (u *UseCase) SetFilmInfo(ctx context.Context, film entities.Film) bool {
	return u.Repo.SetFilmInfo(ctx, film)
}

func (u *UseCase) SetActorInfo(ctx context.Context, actor entities.Actor) bool {
	return u.Repo.SetActorInfo(ctx, actor)
}
