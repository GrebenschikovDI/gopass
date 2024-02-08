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
)

const Dsn = ""
const migr = "migrations"
const defaultRun = "localhost:8000"

func main() {
	log := logger.Initialize("debug")
	db, err := data.InitDB(context.Background(), Dsn, migr)
	if err != nil {
		log.WithField("error", err).Error("init DB failed")
	}
	userUseCase := users.NewUseCase(db.UserRepo)
	recordUseCase := records.NewUseCase(db.RecordRepo)

	_, err = userUseCase.RegisterUser(context.Background(), "Lol", "password")

	//rc := records.Record{
	//	Name:      "sdfsdf",
	//	Site:      "sdfsdf",
	//	Login:     "dsfsdf",
	//	Password:  "dsfsdf",
	//	Info:      "sdfsdf",
	//	CreatedAt: time.Time{},
	//}

	//create, err := db.RecordRepo.Create(context.Background(), &rc)
	//if err != nil {
	//	fmt.Printf("%e", err)
	//}
	//fmt.Println(create)

	s := &http.Server{
		Addr:    defaultRun,
		Handler: transport.ServerRouter(*userUseCase, *recordUseCase, log),
	}

	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.WithField("error", err).Fatal("Could not start server")
	}
}
