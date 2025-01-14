// Package retropg is a PostgreSQL database implementation of the retro.Storer
// interface (https://github.com/kralamoure/retro).
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

type Db struct {
	pool *pgxpool.Pool
}

func NewDb(pool *pgxpool.Pool) (*Db, error) {
	if pool == nil {
		return nil, errors.New("pool is nil")
	}

	login := &Db{pool: pool}

	return login, nil
}
