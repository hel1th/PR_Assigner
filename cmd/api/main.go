package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/hel1th/PR_Assigner/internal/config"
	"github.com/hel1th/PR_Assigner/internal/handler"
	"github.com/hel1th/PR_Assigner/internal/repository"
	"github.com/hel1th/PR_Assigner/internal/server"
	"github.com/hel1th/PR_Assigner/internal/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg := config.Load()

	db, err := repository.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}

	log.Info().Msg("successfully connected to database")

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPRRepo(db)
	statsRepo := repository.NewStatsRepository(db)

	userService := service.NewUserService(userRepo)
	teamService := service.NewTeamService(teamRepo)
	prService := service.NewPullReqSvc(prRepo, userRepo, teamRepo)
	statsService := service.NewStatsService(statsRepo)

	teamHandler := handler.NewTeamHandler(teamService)
	userHandler := handler.NewUserHandler(userService)
	prHandler := handler.NewPullRequestHandler(prService)
	healthHandler := handler.NewHealthHandler()
	statsHandler := handler.NewStatsHandler(statsService)

	r := chi.NewRouter()
	server.RegisterMiddleware(r)

	teamHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	prHandler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)
	statsHandler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Info().Msg("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server stopped gracefully")
}
