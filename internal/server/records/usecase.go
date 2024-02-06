package records

import (
	"context"

	"github.com/pkg/errors"
)

type UseCase struct {
	repo Repo
}

func NewUseCase(repo Repo) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Create(ctx context.Context, record *Record) (*Record, error) {
	r, err := u.repo.Create(ctx, record)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating record")
	}
	return r, nil
}

func (u *UseCase) Update(ctx context.Context, id int, name, site, login, password, info string) (*Record, error) {
	record, err := u.repo.Update(ctx, id, name, site, login, password, info)
	if err != nil {
		return nil, errors.Wrapf(err, "error updatng record, id: %d", id)
	}
	return record, nil
}

func (u *UseCase) Delete(ctx context.Context, id int) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "error deleteing record, id: %d", id)
	}
	return nil
}

func (u *UseCase) List(ctx context.Context, userID int) ([]*Record, error) {
	records, err := u.repo.List(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "get records by userID err, id: %d", userID)
	}
	return records, nil
}

func (u *UseCase) GetById(ctx context.Context, id int) (*Record, error) {
	record, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "get record by id err, id: %d", id)
	}
	return record, nil
}
