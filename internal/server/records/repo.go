package records

import "context"

type Repo interface {
	Create(ctx context.Context, record *Record) (*Record, error)
	Update(ctx context.Context, id, userID int, name, site, login, password, info string) (*Record, error)
	Delete(ctx context.Context, id, userID int) error
	List(ctx context.Context, userID int) ([]*Record, error)
	GetById(ctx context.Context, id int) (*Record, error)
}
