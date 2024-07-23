package main

import (
	_ "HomeWork5/docs"
	"HomeWork5/internal/storage"
	"HomeWork5/internal/user"
	"HomeWork5/internal/ws"
	"HomeWork5/router"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// @title           RESTful Chat Web Server
// @version         1.0
// @description     This is a server for instant messaging
// @termsOfService  http://swagger.io/terms/

// @contact.name   Katia
// @contact.url    github.com/tlb-katia
// @contact.email  tlb-kei7@yandex.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          http://localhost:8080/swagger/index.html

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("%w", err)
	}

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error starting http server", slog.String("error", err.Error()))
		return
	}
	defer logFile.Close()

	log := configureLogger(os.Getenv("ENV"), logFile)
	log.Info("starting web server")

	db, err := storage.NewDB()
	if err != nil {
		log.Error("Failed to connect to database:", err)
	}
	defer db.Close()

	userRep := user.NewRepository(db)
	userService := user.NewService(userRep)
	userHandler := user.NewHandler(log, userService)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(log, hub)

	r := router.InitRouter(log, userHandler, wsHandler)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Error("Failed to start server:", err)
	}

}

func configureLogger(env string, fileName *os.File) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "local":
		log = slog.New(
			slog.NewTextHandler(fileName, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}
