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
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	actor := entities.Actor{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("method: %s, url: %s, error reading body add actor: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(b, &actor)
	id, err := s.Service.UseCase.AddActor(context.Background(), actor)
	if err != nil {
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		log.Printf("method: %s, url: %s, error adding actor: %s", r.Method, r.URL, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "added actor with id: %d", id)
}

func (s *Server) AddFilm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	film := entities.Film{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("method: %s, url: %s, error reading body add film: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(b, &film)
	id, err := s.Service.UseCase.AddFilm(context.Background(), film)
	if err != nil {
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		log.Printf("method: %s, url: %s, error adding film: %s", r.Method, r.URL, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "added film with id: %d", id)
}

func (s *Server) SetActorInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	id := r.PathValue("id")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("method: %s, url: %s, error reading body set actor info: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}

	m := map[string]interface{}{"id": id}
	json.Unmarshal(b, &m)
	ok := s.Service.SetActorInfo(context.Background(), m)
	if !ok {
		fmt.Fprintf(w, "%d bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "set complete:\n%s", string(b))
}

func (s *Server) SetFilmInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	id := r.PathValue("id")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("method: %s, url: %s, error reading body set film info: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}
	m := map[string]interface{}{"id": id}
	json.Unmarshal(b, &m)
	ok := s.Service.SetFilmInfo(context.Background(), m)
	if !ok {
		fmt.Fprintf(w, "%d bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "set complete:\n%s", string(b))
}

func (s *Server) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	film := entities.Film{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("method: %s, url: %s, error reading body delete film: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}

	json.Unmarshal(b, &film)
	ok := s.Service.DeleteFilm(context.Background(), film)
	if !ok {
		fmt.Fprintf(w, "%d bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "deleted:\n%s", string(b))
}

func (s *Server) DeleteActor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Header.Set("Content-Type", "application/json")
	actor := entities.Actor{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("method: %s, url: %s, error reading body delete actor: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}

	json.Unmarshal(b, &actor)
	ok := s.Service.DeleteActor(context.Background(), actor)
	if !ok {
		fmt.Fprintf(w, "%d bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "deleted:\n%s", string(b))
}

func (s *Server) GetFilmInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Header.Set("Content-Type", "application/json")
	title := r.URL.Query().Get("film")
	if title == "" {
		fmt.Fprintf(w, "%d bad request", http.StatusBadRequest)
		return
	}

	sortBy := map[string]struct{}{
		"title":   {},
		"rating":  {},
		"release": {}}
	order := r.URL.Query().Get("order")
	_, ok := sortBy[order]
	if order == "" || !ok {
		order = "rating"
	}

	films, err := s.Service.UseCase.Repo.GetFilmInfo(context.Background(), title, order)
	if err != nil {
		log.Printf("method: %s, url: %s, error get film info: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}

	if len(films) == 0 {
		fmt.Fprintf(w, "%d not found", http.StatusNotFound)
		return
	}

	b, err := json.MarshalIndent(films, "", "    ")
	if err != nil {
		log.Printf("method: %s, url: %s, error serialized get film info: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}

func (s *Server) GetActorInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "%d method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Header.Set("Content-Type", "application/json")
	fullname := r.URL.Query().Get("actor")
	if fullname == "" {
		fmt.Fprintf(w, "%d bad request", http.StatusBadRequest)
		return
	}

	actors, err := s.Service.UseCase.Repo.GetActorInfo(context.Background(), fullname)
	if err != nil {
		log.Printf("method: %s, url: %s, error get actor info: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}

	if len(actors) == 0 {
		fmt.Fprintf(w, "%d not found", http.StatusNotFound)
		return
	}

	b, err := json.MarshalIndent(actors, "", "    ")
	if err != nil {
		log.Printf("method: %s, url: %s, error serialized get actor info: %s", r.Method, r.URL, err)
		fmt.Fprintf(w, "%d internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}
