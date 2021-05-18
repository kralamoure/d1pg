package retropg

import (
	"context"
	"time"

	"github.com/kralamoure/retro"
)

func (r *Db) CreateTicket(ctx context.Context, ticket retro.Ticket) (id string, err error) {
	query := "INSERT INTO retro.tickets (account_id, gameserver_id)" +
		" VALUES ($1, $2)" +
		" ON CONFLICT (account_id) DO UPDATE" +
		" SET id = DEFAULT, gameserver_id = EXCLUDED.gameserver_id, created = DEFAULT" +
		" RETURNING id;"

	err = dbError(
		r.pool.QueryRow(ctx, query,
			ticket.AccountId, ticket.GameServerId).
			Scan(&id),
	)
	return
}

func (r *Db) DeleteTickets(ctx context.Context, before time.Time) (count int, err error) {
	query := "DELETE FROM retro.tickets" +
		" WHERE created <= $1;"

	tag, err := r.pool.Exec(ctx, query, before)
	if err != nil {
		return
	}
	count = int(tag.RowsAffected())
	return
}

func (r *Db) UseTicket(ctx context.Context, id string) (ticket retro.Ticket, err error) {
	query := "DELETE FROM retro.tickets" +
		" WHERE id = $1" +
		" RETURNING id, account_id, gameserver_id, created;"

	err = dbError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&ticket.Id, &ticket.AccountId, &ticket.GameServerId, &ticket.Created),
	)
	return
}

func (r *Db) Tickets(ctx context.Context) (tickets map[string]retro.Ticket, err error) {
	query := "SELECT id, account_id, gameserver_id, created" +
		" FROM retro.tickets;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	tickets = make(map[string]retro.Ticket)
	for rows.Next() {
		var ticket retro.Ticket
		err = rows.Scan(&ticket.Id, &ticket.AccountId, &ticket.GameServerId, &ticket.Created)
		if err != nil {
			return
		}
		tickets[ticket.Id] = ticket
	}
	return
}

func (r *Db) Ticket(ctx context.Context, id string) (ticket retro.Ticket, err error) {
	query := "SELECT id, account_id, gameserver_id, created" +
		" FROM retro.tickets" +
		" WHERE id = $1;"

	err = dbError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&ticket.Id, &ticket.AccountId, &ticket.GameServerId, &ticket.Created),
	)
	return
}
