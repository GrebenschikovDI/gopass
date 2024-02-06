package data

import (
	"GoPass/internal/server/records"
	"GoPass/internal/server/users"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type PgStorage struct {
	DB         *sql.DB
	UserRepo   users.Repo
	RecordRepo records.Repo
}

func InitDB(_ context.Context, dsn, migrations string) (*PgStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	userRepo := NewUserRepo(db)
	recordRepo := NewRecordRepo(db)
	storage := &PgStorage{
		DB:         db,
		UserRepo:   userRepo,
		RecordRepo: recordRepo,
	}
	return storage, nil
}

func (s *PgStorage) runMigrations(dsn, migrations string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrations), dsn)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
