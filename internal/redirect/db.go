package redirect

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/tombuente/redirect-url/xerrors"
)

type RepositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RepositoryImpl {
	repo := RepositoryImpl{
		db: db,
	}

	return repo
}

func (repo RepositoryImpl) GetURL(ctx context.Context, id int64) (URL, error) {
	query, _, err := squirrel.Select("*").
		From("urls").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return URL{}, errors.Join(xerrors.ErrSQLInternal, err)
	}

	row := repo.db.QueryRowxContext(ctx, query, id)
	var i URL
	err = row.Scan(&i.ID, &i.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return URL{}, errors.Join(xerrors.ErrSQLNotFound, err)
		}
		return URL{}, errors.Join(xerrors.ErrSQLInternal, err)
	}

	return i, err
}

func (repo RepositoryImpl) CreateURL(ctx context.Context, params URLParams) (URL, error) {
	query, args, err := squirrel.Insert("urls").
		Columns("url").
		Values(params.URL).
		Suffix("RETURNING id, url").
		ToSql()
	if err != nil {
		return URL{}, errors.Join(xerrors.ErrSQLInternal, err)
	}

	row := repo.db.QueryRowxContext(ctx, query, args...)
	var i URL
	err = row.Scan(&i.ID, &i.URL)
	if err != nil {
		return URL{}, errors.Join(xerrors.ErrSQLInternal, err)
	}

	return i, err
}
