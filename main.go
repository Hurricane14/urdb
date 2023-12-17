package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urdb/server"
)

func main() {
	if err := initRepositories(); err != nil {
		panic(err)
	}

	server := server.New(users, movies)
	go func() {
		err := server.Run(8080)
		if errors.Is(err, http.ErrServerClosed) {
			return
		} else if err != nil {
			panic(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	if err := server.Shutdown(5 * time.Second); err != nil {
		panic(err)
	}
}
