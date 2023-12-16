package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	repo, err := newRepository()
	if err != nil {
		panic(err)
	}

	server := newServer(repo)
	go func() {
		err := server.run(8080)
		if errors.Is(err, http.ErrServerClosed) {
			return
		} else if err != nil {
			panic(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	if err := server.shutdown(5 * time.Second); err != nil {
		panic(err)
	}
}
