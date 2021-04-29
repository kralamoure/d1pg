// Package d1pg implements d1 repository interface for PostgreSQL.
package d1pg

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var defaultTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadWrite,
	DeferrableMode: pgx.NotDeferrable,
}

var errInvalidAssertion = errors.New("invalid assertion")

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) (*Repo, error) {
	if pool == nil {
		return nil, errors.New("pool is nil")
	}

	login := &Repo{pool: pool}

	return login, nil
}
