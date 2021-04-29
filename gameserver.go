package d1pg

import (
	"context"

	"github.com/kralamoure/d1"
	"github.com/kralamoure/d1/d1typ"
)

func (r *Repo) CreateGameServer(ctx context.Context, gameServer d1.GameServer) error {
	query := "INSERT INTO d1.gameservers (id, host, port, state, completion)" +
		" VALUES ($1, $2, $3, $4, $5);"

	_, err := r.pool.Exec(ctx, query,
		gameServer.Id, gameServer.Host, gameServer.Port, gameServer.State, gameServer.Completion)
	return repoError(err)
}

func (r *Repo) GameServers(ctx context.Context) (gameServers map[int]d1.GameServer, err error) {
	query := "SELECT id, host, port, state, completion" +
		" FROM d1.gameservers;"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	gameServers = make(map[int]d1.GameServer)
	for rows.Next() {
		var gameServer d1.GameServer
		err = rows.Scan(&gameServer.Id, &gameServer.Host, &gameServer.Port, &gameServer.State, &gameServer.Completion)
		if err != nil {
			return
		}
		gameServers[gameServer.Id] = gameServer
	}
	return
}

func (r *Repo) GameServer(ctx context.Context, id int) (gameServer d1.GameServer, err error) {
	query := "SELECT id, host, port, state, completion" +
		" FROM d1.gameservers" +
		" WHERE id = $1;"

	err = repoError(
		r.pool.QueryRow(ctx, query, id).
			Scan(&gameServer.Id, &gameServer.Host, &gameServer.Port, &gameServer.State, &gameServer.Completion),
	)
	return
}

func (r *Repo) SetGameServerState(ctx context.Context, id int, state d1typ.GameServerState) error {
	query := "UPDATE d1.gameservers" +
		" SET state = $2" +
		" WHERE id = $1;"

	tag, err := r.pool.Exec(ctx, query, id, state)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return d1.ErrNotFound
	}
	return nil
}
