package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/filmoteka/repository"
	"github.com/filmoteka/usecase"
)

type Server struct {
	Conn repository.DBHandler
	Serv *http.Server
}

func NewServer(conn repository.DBHandler) *Server {
	mux := http.NewServeMux()

	serv := &Server{Conn: conn,
		Serv: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}

	mux.HandleFunc("/get_film", serv.GetFilmInfo)
	return serv
}

func (s *Server) Run(ctx context.Context) {
	op := "server.Run"
	fmt.Printf("Start server on port %s ...\n", s.Serv.Addr)
	log.Fatalf("foo=%s resume=%v", op, s.Serv.ListenAndServe())
}

func (s *Server) GetFilmInfo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("film")
	if title == "" {
		return
	}

	uc := usecase.UseCase{Repo: s.Conn}

	films, err := uc.GetFilmInfo(context.Background(), title)
	if err != nil {
		log.Fatal("inner server error: ", err)
	}

	b, err := json.MarshalIndent(films, "", "    ")
	if err != nil {
		log.Fatal("error serialized:", err)
	}
	fmt.Fprint(w, string(b))
}
