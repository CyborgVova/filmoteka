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

	mux.HandleFunc("/get_film", serv.GetFilmInfo)
	mux.HandleFunc("/get_actor", serv.GetActorInfo)
	mux.HandleFunc("/add_actor", serv.AddActor)
	mux.HandleFunc("/add_film", serv.AddFilm)
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
func (s *Server) GetFilmInfo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("film")
	if title == "" {
		return
	}

	films, err := s.Service.UseCase.Repo.GetFilmInfo(context.Background(), title)
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
