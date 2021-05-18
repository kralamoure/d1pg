package retropg

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kralamoure/retro"
)

const (
	errUniqueViolation errCode = "23505"
)

type errCode string

func storerError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", retro.ErrNotFound, err)
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}

	if errCode(pgErr.Code) != errUniqueViolation {
		return err
	}

	var storerErr error
	switch pgErr.ConstraintName {
	case "gameservers_host_port_key":
		storerErr = retro.ErrGameServerHostAndPortAlreadyExist
	case "characters_name_gameserver_id_key":
		storerErr = retro.ErrCharacterNameAndGameServerIdAlreadyExist
	case "tickets_account_id_key":
		storerErr = retro.ErrTicketAccountIdAlreadyExists
	default:
		storerErr = retro.ErrAlreadyExists
	}

	return fmt.Errorf("%w: %s", storerErr, err)
}
