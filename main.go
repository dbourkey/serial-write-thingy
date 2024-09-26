package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dbourkey/serial-write-thingy/status"
	"golang.org/x/sync/errgroup"
)

func main() {
	mockDBClient := status.NewMockDBClient(0.5)
	serialiser := status.NewSerialiser(mockDBClient, 3*time.Second)
	updateHandler := status.NewStatusHandler(serialiser)
	server := status.NewServer(updateHandler)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		if err := serialiser.Run(ctx); err != nil {
			return fmt.Errorf("serialiser runtime failure: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		if err := server.ListenAndServe(); err != nil {
			return fmt.Errorf("http server runtime failure: %w", err)
		}
		return http.ErrServerClosed
	})

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	eg.Go(func() error {
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown timeout: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("runtime failure: %w", err)
	}
}
