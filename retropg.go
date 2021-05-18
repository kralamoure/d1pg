// Package retropg is a library that implements retro.Storer interface (https://github.com/kralamoure/retro) for a PostgreSQL database.
package retropg

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

type Storer struct {
	pool *pgxpool.Pool
}

func NewStorer(pool *pgxpool.Pool) (*Storer, error) {
	if pool == nil {
		return nil, errors.New("pool is nil")
	}

	login := &Storer{pool: pool}

	return login, nil
}
