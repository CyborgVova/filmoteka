package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/filmoteka/entities"
	"github.com/filmoteka/repository"
	"github.com/filmoteka/server/middleware"
	"github.com/filmoteka/service"
	"github.com/filmoteka/usecase"
)

type Server struct {
	Conn    repository.DBHandler
	Serv    *http.Server
	Service *service.Service
}

func NewServer(conn repository.DBHandler) *Server {
	mux := http.NewServeMux()

	serv := &Server{Conn: conn,
		Serv: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		Service: &service.Service{
			UseCase: &usecase.UseCase{
				Repo: conn,
			},
		},
	}

	mux.Handle("/get_film", middleware.Logging(http.HandlerFunc(serv.GetFilmInfo)))
	mux.Handle("/get_actor", middleware.Logging(http.HandlerFunc(serv.GetActorInfo)))
	mux.Handle("/add_actor", middleware.Logging(middleware.Authorization(http.HandlerFunc(serv.AddActor))))
	mux.Handle("/add_film", middleware.Logging(middleware.Authorization(http.HandlerFunc(serv.AddFilm))))
	mux.Handle("/delete_film", middleware.Logging(middleware.Authorization(http.HandlerFunc(serv.DeleteFilm))))
	mux.Handle("/delete_actor", middleware.Logging(middleware.Authorization(http.HandlerFunc(serv.DeleteActor))))
	mux.Handle("/set_actor/{id}", middleware.Logging(middleware.Authorization(http.HandlerFunc(serv.SetActorInfo))))
	mux.Handle("/set_film/{id}", middleware.Logging(middleware.Authorization(http.HandlerFunc(serv.SetFilmInfo))))
	return serv
}

func (s *Server) Run(ctx context.Context) {
	op := "server.Run"
	fmt.Printf("Start server on port %s ...\n", s.Serv.Addr)
	log.Fatalf("foo=%s resume=%v", op, s.Serv.ListenAndServe())
}

func (s *Server) AddActor(w http.ResponseWriter, r *http.Request) {
	actor := entities.Actor{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading body:", err)
	}
	json.Unmarshal(b, &actor)
	s.Service.UseCase.AddActor(context.Background(), actor)
}

func (s *Server) AddFilm(w http.ResponseWriter, r *http.Request) {
	film := entities.Film{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading body:", err)
	}
	json.Unmarshal(b, &film)
	s.Service.UseCase.AddFilm(context.Background(), film)
}

func (s *Server) SetActorInfo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error set actor:", err)
	}
	m := map[string]interface{}{"id": id}
	json.Unmarshal(b, &m)
	s.Service.SetActorInfo(context.Background(), m)
}

func (s *Server) SetFilmInfo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error set film:", err)
	}
	m := map[string]interface{}{"id": id}
	json.Unmarshal(b, &m)
	s.Service.SetFilmInfo(context.Background(), m)
}

func (s *Server) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	film := entities.Film{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading body:", err)
	}

	json.Unmarshal(b, &film)
	s.Service.DeleteFilm(context.Background(), film)
}

func (s *Server) DeleteActor(w http.ResponseWriter, r *http.Request) {
	actor := entities.Actor{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading body:", err)
	}

	json.Unmarshal(b, &actor)
	s.Service.DeleteActor(context.Background(), actor)
}

func (s *Server) GetFilmInfo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("film")
	if title == "" {
		return
	}

	sortBy := map[string]struct{}{
		"title":       {},
		"description": {},
		"rating":      {},
		"release":     {}}
	order := r.URL.Query().Get("order")
	_, ok := sortBy[order]
	if order == "" || !ok {
		order = "rating"
	}

	films, err := s.Service.UseCase.Repo.GetFilmInfo(context.Background(), title, order)
	if err != nil {
		log.Fatal("inner server error: ", err)
	}

	b, err := json.MarshalIndent(films, "", "    ")
	if err != nil {
		log.Fatal("error serialized:", err)
	}
	fmt.Fprint(w, string(b))
}

func (s *Server) GetActorInfo(w http.ResponseWriter, r *http.Request) {
	fullname := r.URL.Query().Get("actor")
	if fullname == "" {
		return
	}

	actors, err := s.Service.UseCase.Repo.GetActorInfo(context.Background(), fullname)
	if err != nil {
		log.Fatal("inner server error: ", err)
	}

	b, err := json.MarshalIndent(actors, "", "    ")
	if err != nil {
		log.Fatal("error serialized:", err)
	}
	fmt.Fprint(w, string(b))
}
