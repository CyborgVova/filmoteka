package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/filmoteka/repository/postgres"
	"github.com/filmoteka/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repository := postgres.NewRepository()
	defer repository.Conn.Close(ctx)
	s := server.NewServer(repository)
	go s.Run(ctx)
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Serv.Shutdown(ctx)
}
