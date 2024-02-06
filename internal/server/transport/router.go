package transport

import (
	"GoPass/internal/server/records"
	"GoPass/internal/server/transport/handlers"
	mw "GoPass/internal/server/transport/middleware"
	"GoPass/internal/server/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func ServerRouter(
	userUseCase users.UseCase,
	recordUseCase records.UseCase,
	log *logrus.Logger,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(mw.LoggerMiddleware(log))
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		orderHandler := handlers.NewRecordHandler(recordUseCase)
		r.Get("/api/records/id", orderHandler.GetById)
		r.Get("/api/records", orderHandler.List)
		r.Post("/api/records", orderHandler.Create)
		r.Patch("/api/records", orderHandler.Update)
		r.Delete("/api/records", orderHandler.Delete)
	})

	userHandler := handlers.NewUserHandler(userUseCase)
	r.Post("/api/user/register", userHandler.RegisterUser)
	r.Post("/api/user/login", userHandler.LoginUser)

	return r
}
