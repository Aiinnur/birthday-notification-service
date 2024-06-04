package main

import (
	"birthday-notification-service/internal/config"
	"birthday-notification-service/internal/http-server/handlers"
	"birthday-notification-service/internal/repository"
	"birthday-notification-service/internal/sendEmail"
	"birthday-notification-service/pkg/postgres"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	postgresClient, err := postgres.NewPostgresClient(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to crate new PostgreSQL client")
	}

	repo := repository.NewRepository(postgresClient)

	router := chi.NewRouter()

	router.Route("/user", func(r chi.Router) {
		r.Use(middleware.BasicAuth("birthday-notification", map[string]string{
			cfg.Server.User: cfg.Server.Password,
		}))
		r.Post("/", handlers.New(repo))
	})

	router.Route("/subscribe", func(r chi.Router) {
		r.Use(middleware.BasicAuth("birthday-notification", map[string]string{
			cfg.Server.User: cfg.Server.Password,
		}))
		r.Post("/", handlers.NewSubscribe(repo))
	})

	router.Route("/unsubscribe", func(r chi.Router) {
		r.Use(middleware.BasicAuth("birthday-notification", map[string]string{
			cfg.Server.User: cfg.Server.Password,
		}))
		r.Post("/", handlers.NewUnsubscribe(repo))
	})

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.TimeOut,
		WriteTimeout: cfg.Server.TimeOut,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	c := cron.New()
	_, err = c.AddFunc("@daily", func() { sendEmail.SendBirthdayNotifications(cfg, repo) })
	if err != nil {
		log.Fatalf("Unable to schedule the cron job: %v", err)
	}
	c.Start()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stopChan
	log.Println("Shutdown signal received, exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown the server: %v", err)
	}

	c.Stop()
	log.Println("Cron scheduler stopped")

}
