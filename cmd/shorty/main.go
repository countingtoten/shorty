package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/countingtoten/shorty"
	"github.com/countingtoten/shorty/handler"
	memory "github.com/countingtoten/shorty/in-memory"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg := shorty.Config{}

	err := env.Parse(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to parse env config")
	}

	datastore := memory.NewDatastore(cfg)

	mux := handler.New(datastore, logger)

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.BindPort)

	server := &http.Server{
		Handler:        mux,
		Addr:           addr,
		WriteTimeout:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info().Str("bind_address", addr).Msg("Shorty started")
		err = server.ListenAndServe()
		if err != nil {
			logger.Fatal().Err(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	logger.Info().Msg("Shorty shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatal().Err(err)
	}

	logger.Info().Msg("Shorty shut down complete")
}
