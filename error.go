package d1pg

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/kralamoure/d1"
)

const (
	errUniqueViolation errCode = "23505"
)

type errCode string

func repoError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", d1.ErrNotFound, err)
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}

	if errCode(pgErr.Code) != errUniqueViolation {
		return err
	}

	var repoErr error
	switch pgErr.ConstraintName {
	case "gameservers_host_port_key":
		repoErr = d1.ErrGameServerHostAndPortAlreadyExist
	case "characters_name_gameserver_id_key":
		repoErr = d1.ErrCharacterNameAndGameServerIdAlreadyExist
	case "tickets_account_id_key":
		repoErr = d1.ErrTicketAccountIdAlreadyExists
	default:
		repoErr = d1.ErrAlreadyExists
	}

	return fmt.Errorf("%w: %s", repoErr, err)
}
