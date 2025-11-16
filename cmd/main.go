package main

import (
	"avito-backend-trainee-assignment-autumn-2025/config"
	"avito-backend-trainee-assignment-autumn-2025/internal/app/handlers"
	"avito-backend-trainee-assignment-autumn-2025/internal/app/routing"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/service"
	"avito-backend-trainee-assignment-autumn-2025/internal/infrastructure/persistance/postgres"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	cfg := config.GetConfig()
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = config.NewConnection(cfg)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migrationsDirectory := filepath.Join(projectRoot, "migrations")
	if err := postgres.Migrate(db, migrationsDirectory); err != nil {
		log.Fatal("migration failed", err)
	}
	userRepository := postgres.NewUserRepositoryPostgres(db)
	teamRepository := postgres.NewTeamRepositoryPostgres(db)
	prRepository := postgres.NewPullRequestRepositoryPostgres(db)
	statsRepository := postgres.NewStatsPostgresRepository(db)
	userService := service.NewUserService(userRepository)
	teamService := service.NewTeamService(teamRepository)
	prService := service.NewPullRequestService(prRepository, userRepository)
	statsService := service.NewStatsService(statsRepository)
	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	prHandler := handlers.NewPullRequestHandler(prService)
	statsHandler := handlers.NewStatsHandler(statsService)
	router := routing.CreateRouter(prHandler, userHandler, teamHandler, statsHandler)

	server := &http.Server{
		Addr:         config.GetPort(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to listen and serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}
