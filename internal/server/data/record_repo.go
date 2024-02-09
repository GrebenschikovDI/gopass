package data

import (
	"GoPass/internal/server/records"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type recordRepo struct {
	db *sql.DB
}

func NewRecordRepo(db *sql.DB) records.Repo {
	return &recordRepo{
		db: db,
	}
}

func (r *recordRepo) Create(ctx context.Context, record *records.Record) (*records.Record, error) {
	query := `
		INSERT INTO records (user_id, name, site, login, password_hash, info)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		record.UserID,
		record.Name,
		record.Site,
		record.Login,
		record.Password,
		record.Info,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert record: %v", err)
	}

	record.ID = id
	return record, nil
}

func (r *recordRepo) Update(ctx context.Context, id, userID int, name, site, login, password, info string) (*records.Record, error) {
	query := `
		UPDATE records
		SET name=$1, site=$2, login=$3, password_hash=$4, info=$5
		WHERE id=$6 AND user_id=$7
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		name,
		site,
		login,
		password,
		info,
		id,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update record: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return nil, errors.New("record not found or not owned by user")
	}

	return r.GetById(ctx, id)
}

func (r *recordRepo) Delete(ctx context.Context, id, userID int) error {
	query := `
		DELETE FROM records
		WHERE id=$1 AND user_id=$2
	`

	res, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("record not found or not owned by user")
	}

	return nil
}

func (r *recordRepo) List(ctx context.Context, userID int) ([]*records.Record, error) {
	query := `
		SELECT id, user_id, name, site, login, password_hash, info, created_at
		FROM records
		WHERE user_id=$1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %v", err)
	}
	defer rows.Close()

	var recordsList []*records.Record
	for rows.Next() {
		var record records.Record
		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Name,
			&record.Site,
			&record.Login,
			&record.Password,
			&record.Info,
			&record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan record: %v", err)
		}
		recordsList = append(recordsList, &record)
	}

	return recordsList, nil
}

func (r *recordRepo) GetById(ctx context.Context, id int) (*records.Record, error) {
	query := `
		SELECT id, user_id, name, site, login, password_hash, info, created_at
		FROM records
		WHERE id=$1
	`

	var record records.Record
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.UserID,
		&record.Name,
		&record.Site,
		&record.Login,
		&record.Password,
		&record.Info,
		&record.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found or not owned by user")
		}
		return nil, fmt.Errorf("failed to get record: %v", err)
	}

	return &record, nil
}
