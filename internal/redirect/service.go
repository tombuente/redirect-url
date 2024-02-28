package redirect

import (
	"context"
	"errors"

	"github.com/tombuente/redirect-url/xerrors"
)

type URL struct {
	ID  int64
	URL string
}

type URLParams struct {
	URL string
}

type DB interface {
	GetURL(ctx context.Context, id int64) (URL, error)
	CreateURL(ctx context.Context, params URLParams) (URL, error)
}

type RedirectService struct {
	db DB
}

func NewRedirectService(db DB) RedirectService {
	return RedirectService{
		db: db,
	}
}

func (s RedirectService) GetURL(ctx context.Context, id int64) (URL, error) {
	url, err := s.db.GetURL(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, xerrors.ErrSQLNotFound):
			return URL{}, errors.Join(xerrors.ErrNotFound, err)
		default:
			return URL{}, errors.Join(xerrors.ErrInternal, err)
		}
	}

	return url, nil
}

func (s RedirectService) CreateURL(ctx context.Context, params URLParams) (URL, error) {
	url, err := s.db.CreateURL(ctx, params)
	if err != nil {
		return URL{}, errors.Join(xerrors.ErrInternal, err)
	}

	return url, nil
}
