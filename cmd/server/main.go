package main

import (
	"HomeWork5/internal/storage"
	"HomeWork5/internal/user"
	"HomeWork5/router"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// TODO configure logger, user authentication (jwt)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("%w", err)
	}

	log := configureLogger(os.Getenv("ENV"))
	log.Info("starting url-shortener")

	db, err := storage.NewDB()
	if err != nil {
		log.Error("Failed to connect to database:", err)
	}
	defer db.Close()

	userRep := user.NewRepository(db)
	userService := user.NewService(userRep)
	userHandler := user.NewHandler(userService)

	r := router.InitRouter(userHandler)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Error("Failed to start server:", err)
	}

}

func configureLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}
