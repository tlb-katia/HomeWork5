package router

import (
	"HomeWork5/internal/middleware"
	"HomeWork5/internal/user"
	"HomeWork5/internal/ws"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
)

func InitRouter(logger *slog.Logger, userHandler *user.Handler, wsHandler *ws.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware(logger))

	r.Post("/signup", userHandler.CreateUser)
	r.Post("/login", userHandler.LoginUser)
	r.Get("/logout", userHandler.LogoutUser)

	r.Post("/ws/CreateRoom", wsHandler.CreateRoom)
	r.Get("/ws/JoinRoom/:roomId", wsHandler.JoinRoom)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
