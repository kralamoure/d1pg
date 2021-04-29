// Package d1pg is a library that implements d1.Repo interface (https://github.com/kralamoure/d1) for a PostgreSQL database.
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
