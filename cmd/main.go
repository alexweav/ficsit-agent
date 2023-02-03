package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/alexweav/ficsit-agent/pkg/agent"
)

func main() {
	l := log.New(os.Stdout, "", 0)

	agent := agent.New(l)

	inter := make(chan os.Signal, 1)
	signal.Notify(inter, os.Interrupt)

	errCh := make(chan error)
	go func() {
		errCh <- agent.Run(context.Background())
	}()

	select {
	case err := <-errCh:
		if err != nil {
			l.Fatalf("exited with fatal error: %s", err.Error())
		}
	// Not great, doesn't wait on anything to close
	case <-inter:
		l.Println("Exiting...")
		os.Exit(0)
	}

}
