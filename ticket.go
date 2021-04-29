package d1pg

import (
	"context"
	"time"

	"github.com/kralamoure/d1"
)

func (r *Repo) CreateTicket(ctx context.Context, ticket d1.Ticket) (id string, err error) {
	query := "INSERT INTO d1.tickets (account_id, gameserver_id)" +
		" VALUES ($1, $2)" +
		" ON CONFLICT (account_id) DO UPDATE" +
		" SET id = DEFAULT, gameserver_id = EXCLUDED.gameserver_id, created = DEFAULT" +
		" RETURNING id;"

	err = repoError(
		r.pool.QueryRow(ctx, query,
			ticket.AccountId, ticket.GameServerId).
			Scan(&id),
	)
	return
}

func (r *Repo) DeleteTickets(ctx context.Context, before time.Time) (count int, err error) {
	query := "DELETE FROM d1.tickets" +
		" WHERE created <= $1;"

	tag, err := r.pool.Exec(ctx, query, before)
	if err != nil {
		return
	}
	count = int(tag.RowsAffected())
	return
}

func (r *Repo) UseTicket(ctx context.Context, id string) (ticket d1.Ticket, err error) {
	query := "DELETE FROM d1.tickets" +
		" WHERE id = $1" +
		" RETURNING id, account_id, gameserver_id, created;"

	err = repoError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&ticket.Id, &ticket.AccountId, &ticket.GameServerId, &ticket.Created),
	)
	return
}

func (r *Repo) Tickets(ctx context.Context) (tickets map[string]d1.Ticket, err error) {
	query := "SELECT id, account_id, gameserver_id, created" +
		" FROM d1.tickets;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	tickets = make(map[string]d1.Ticket)
	for rows.Next() {
		var ticket d1.Ticket
		err = rows.Scan(&ticket.Id, &ticket.AccountId, &ticket.GameServerId, &ticket.Created)
		if err != nil {
			return
		}
		tickets[ticket.Id] = ticket
	}
	return
}

func (r *Repo) Ticket(ctx context.Context, id string) (ticket d1.Ticket, err error) {
	query := "SELECT id, account_id, gameserver_id, created" +
		" FROM d1.tickets" +
		" WHERE id = $1;"

	err = repoError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&ticket.Id, &ticket.AccountId, &ticket.GameServerId, &ticket.Created),
	)
	return
}
