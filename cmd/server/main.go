package main

import (
	"GoPass/internal/server/data"
	"GoPass/internal/server/logger"
	"GoPass/internal/server/records"
	"GoPass/internal/server/transport"
	"GoPass/internal/server/users"
	"context"
	"errors"
	"net/http"
	"time"
)

const Dsn = "postgres://gopher:1234@localhost:5432/gopass"
const migr = "migrations"
const defaultRun = "localhost:8000"

func main() {
	log := logger.Initialize("info")
	db, err := data.InitDB(context.Background(), Dsn, migr)
	if err != nil {
		log.WithField("error", err).Error("init DB failed")
	}
	userUseCase := users.NewUseCase(db.UserRepo)
	recordUseCase := records.NewUseCase(db.RecordRepo)

	rc := records.Record{
		UserID:    123,
		Name:      "sdfsdf",
		Site:      "sdfsdf",
		Login:     "dsfsdf",
		Password:  "dsfsdf",
		Info:      "sdfsdf",
		CreatedAt: time.Time{},
	}

	db.RecordRepo.Create(context.Background(), &rc)

	s := &http.Server{
		Addr:    defaultRun,
		Handler: transport.ServerRouter(*userUseCase, *recordUseCase, log),
	}

	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.WithField("error", err).Fatal("Could not start server")
	}
}
