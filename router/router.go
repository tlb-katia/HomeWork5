package router

import (
	"HomeWork5/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(userHandler *user.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Post("/signup", userHandler.CreateUser)
	r.Post("/login", userHandler.LoginUser)
	r.Get("/logout", userHandler.LogoutUser)

	return r
}
